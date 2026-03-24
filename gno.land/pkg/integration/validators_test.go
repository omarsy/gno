package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	"github.com/gnolang/gno/gno.land/pkg/gnoland/ugnot"
	vmm "github.com/gnolang/gno/gno.land/pkg/sdk/vm"
	abci "github.com/gnolang/gno/tm2/pkg/bft/abci/types"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/log"
	"github.com/gnolang/gno/tm2/pkg/std"

	gnoenv "github.com/gnolang/gno/gnovm/pkg/gnoenv"
)

// TestReproduceTest11ChainHalt reproduces the test11 chain halt at block 702499.
//
// On test11, proposals 13 and 14 were executed in consecutive blocks:
//   - Block 702497: Proposal 13 executed → 9 validators added
//   - Block 702498: Proposal 14 executed → 3 validators removed
//   - Block 702499: EndBlocker re-fetched the removals → CRASH
//
// This test starts a real in-memory node with the r/sys/validators/v2 realm,
// executes two consecutive validator proposals, and checks if the node
// survives the next block.
//
// IMPORTANT: The default test config uses PubKeyTypeURLs=[] which silently
// filters all validator updates. We patch it so updates reach consensus.
func TestReproduceTest11ChainHalt(t *testing.T) {
	rootdir := gnoenv.RootDir()

	// Full node config with all packages loaded (validators, GovDAO, valopers)
	cfg, _ := TestingNodeConfig(t, rootdir)

	// CRITICAL FIX: enable validator key types so EndBlocker updates
	// actually reach the consensus layer instead of being silently dropped.
	cfg.Genesis.ConsensusParams.Validator = &abci.ValidatorParams{
		PubKeyTypeURLs: []string{"/tm.PubKeyEd25519", "/tm.PubKeySecp256k1"},
	}

	// Start node
	n, remoteAddr := TestingInMemoryNode(t, log.NewNoopLogger(), cfg)
	defer n.Stop()

	t.Log("Node started with PubKeyTypeURLs=[ed25519, secp256k1]")

	// Setup gnoclient
	kb := keys.NewInMemory()
	_, err := kb.CreateAccount(DefaultAccount_Name, DefaultAccount_Seed, "", "", 0, 0)
	require.NoError(t, err)

	rpcClient, err := rpcclient.NewHTTPClient(remoteAddr)
	require.NoError(t, err)

	client := gnoclient.Client{
		Signer: &gnoclient.SignerFromKeybase{
			Keybase:  kb,
			Account:  DefaultAccount_Name,
			Password: "",
			ChainID:  "tendermint_test",
		},
		RPCClient: rpcClient,
	}

	callerInfo, err := client.Signer.Info()
	require.NoError(t, err)
	caller := callerInfo.GetAddress()

	seq := uint64(0)

	broadcastRun := func(desc, code string) {
		t.Helper()
		msg := vmm.MsgRun{
			Caller: caller,
			Package: &std.MemPackage{
				Name:  "main",
				Files: []*std.MemFile{{Name: "run.gno", Body: code}},
			},
		}
		res, err := client.Run(gnoclient.BaseTxCfg{
			GasFee:         ugnot.ValueString(1000000),
			GasWanted:      95000000,
			SequenceNumber: seq,
		}, msg)
		require.NoError(t, err, "%s: broadcast error", desc)
		require.True(t, res.DeliverTx.IsOK(), "%s: DeliverTx failed: %v", desc, res.DeliverTx)
		seq++
		t.Logf("  %s → height %d", desc, res.Height)
	}

	broadcastCall := func(desc, pkgPath, fn string, args []string, send std.Coins) {
		t.Helper()
		msg := vmm.MsgCall{
			Caller:  caller,
			PkgPath: pkgPath,
			Func:    fn,
			Args:    args,
			Send:    send,
		}
		res, err := client.Call(gnoclient.BaseTxCfg{
			GasFee:         ugnot.ValueString(1000000),
			GasWanted:      95000000,
			SequenceNumber: seq,
		}, msg)
		require.NoError(t, err, "%s: broadcast error", desc)
		require.True(t, res.DeliverTx.IsOK(), "%s: DeliverTx failed: %v", desc, res.DeliverTx)
		seq++
		t.Logf("  %s → height %d", desc, res.Height)
	}

	// Step 0: Init GovDAO
	t.Log("\n--- Init GovDAO ---")
	broadcastRun("Init GovDAO", `package main
import i "gno.land/r/gov/dao/v3/init"
func main() {
	i.InitWithUsers(address("g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"))
}`)

	// Step 1: Register valoper
	t.Log("\n--- Register valoper ---")
	broadcastCall("Register valoper",
		"gno.land/r/gnops/valopers", "Register",
		[]string{
			"myval", "Test validator", "on-prem",
			"g1td0cgmt9uz7kq4hcv7fkkwvp3z4lq4dsewffwr",
			"gpub1pgfj7ard9eg82cjtv4u4xetrwqer2dntxyfzxz3pqg0lte6srklm3tuyja9489n3dsnx4wcadq43wrwnz6nln8s7lf9uyptc3nm",
		},
		std.MustParseCoins("20000000ugnot"),
	)

	const valAddr = "g1td0cgmt9uz7kq4hcv7fkkwvp3z4lq4dsewffwr"

	// Helper: query if address is in the realm's validator set
	queryIsValidator := func(addr string) bool {
		t.Helper()
		res, err := rpcClient.ABCIQuery(context.Background(),
			"vm/qeval",
			[]byte(fmt.Sprintf("gno.land/r/sys/validators/v2.IsValidator(address(\"%s\"))", addr)),
		)
		require.NoError(t, err)
		return strings.Contains(string(res.Response.Data), "true")
	}

	// Verify validator is NOT in the set initially
	require.False(t, queryIsValidator(valAddr), "validator should not exist before ADD")

	// Step 2: Add validator — block N (create + vote + execute in one TX)
	t.Log("\n--- Block N: Execute ADD proposal ---")
	broadcastRun("Add validator", `package main
import (
	"gno.land/r/gnops/valopers/proposal"
	"gno.land/r/gov/dao"
)
func main() {
	pr := proposal.NewValidatorProposalRequest(cross, address("`+valAddr+`"))
	pid := dao.MustCreateProposal(cross, pr)
	dao.MustVoteOnProposalSimple(cross, int64(pid), "YES")
	dao.ExecuteProposal(cross, pid)
}`)

	// Verify validator IS in the set after ADD
	require.True(t, queryIsValidator(valAddr), "validator should exist after ADD proposal")

	// Step 3: Remove validator — block N+1 (mark + create + vote + execute)
	t.Log("\n--- Block N+1: Execute REMOVE proposal ---")
	broadcastRun("Remove validator", `package main
import (
	valopers "gno.land/r/gnops/valopers"
	"gno.land/r/gnops/valopers/proposal"
	"gno.land/r/gov/dao"
)
func main() {
	addr := address("`+valAddr+`")
	valopers.UpdateKeepRunning(cross, addr, false)
	pr := proposal.NewValidatorProposalRequest(cross, addr)
	pid := dao.MustCreateProposal(cross, pr)
	dao.MustVoteOnProposalSimple(cross, int64(pid), "YES")
	dao.ExecuteProposal(cross, pid)
}`)

	// Verify validator is NOT in the set after REMOVE
	require.False(t, queryIsValidator(valAddr), "validator should not exist after REMOVE proposal")

	// Step 4: Trigger block N+2 — this is where the crash happens
	// The stale event from block N+1's fireEvents triggers the EndBlocker,
	// which re-fetches the same removal via GetChanges → crash.
	t.Log("\n--- Block N+2: trigger next block (crash point) ---")
	broadcastRun("Trigger next block", `package main
func main() { println("alive") }`)

	// Step 5: Verify node survived and validator set is correct
	t.Log("\n--- Verify node is alive and valset is correct ---")
	status, err := rpcClient.Status(context.Background(), nil)
	require.NoError(t, err)
	t.Logf("  Node alive at block %d", status.SyncInfo.LatestBlockHeight)

	require.False(t, queryIsValidator(valAddr),
		"validator should still be removed after crash point block")
	t.Log("\n✓ No crash — validator set is correct!")
}
