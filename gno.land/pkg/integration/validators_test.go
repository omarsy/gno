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

// TestEndBlocker_AddRemoveSameValidator tests the case where the same validator
// is added in block N and removed in block N+1.
//
// Root cause: before the fix, EndBlocker called GetChanges(lastBlockHeight)
// with no upper bound. lastBlockHeight is the last *committed* block (N), so
// during block N+1's EndBlocker it fetched changes from block N onward —
// including block N's ADD that was already applied. The validator set received
// both ADD (power=1) and REMOVE (power=0) for the same address, which was
// rejected as a duplicate entry by processChanges() and crashed the node.
//
// The fix bounds GetChanges(from, to) so EndBlocker only fetches the current
// block's changes.
//
// IMPORTANT: The default test config uses PubKeyTypeURLs=[] which silently
// filters all validator updates. We patch it so updates reach consensus.
func TestEndBlocker_AddRemoveSameValidator(t *testing.T) {
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

	// Step 4: Liveness check.
	// Before the fix, the node crashes during Step 3: block N+1's EndBlocker
	// calls GetChanges(N) which re-fetches the ADD from block N alongside the
	// REMOVE from block N+1, sending contradictory updates to consensus.
	// With the fix (bounded GetChanges), the node survives and this confirms it.
	t.Log("\n--- Liveness check: broadcast after crash point ---")
	broadcastRun("Liveness check", `package main
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

// TestEndBlocker_Test11ChainHalt reproduces the test11 chain halt with
// distinct validators: add val X in block N, advance, add val Y in block N+3,
// remove val X in block N+4.
//
// Before the fix, the crash occurred at the liveness check (block after N+4):
//   - Block N+4's EndBlocker: GetChanges(N+3) re-fetched the stale ADD Y from
//     block N+3 alongside the REMOVE X from N+4. Different addresses, so no
//     duplicate error — but the stale ADD Y was re-applied (harmless here).
//   - Block N+5's EndBlocker: the collector still had a stale event from N+4,
//     so GetChanges(N+4) re-fetched the REMOVE X. But val X was already removed
//     → verifyRemovals failed with "failed to find validator to remove" → crash.
//
// The fix bounds GetChanges(from, to) so EndBlocker only fetches the current
// block's changes, preventing stale re-fetches entirely.
func TestEndBlocker_Test11ChainHalt(t *testing.T) {
	rootdir := gnoenv.RootDir()
	cfg, _ := TestingNodeConfig(t, rootdir)
	cfg.Genesis.ConsensusParams.Validator = &abci.ValidatorParams{
		PubKeyTypeURLs: []string{"/tm.PubKeyEd25519", "/tm.PubKeySecp256k1"},
	}

	n, remoteAddr := TestingInMemoryNode(t, log.NewNoopLogger(), cfg)
	defer n.Stop()

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

	queryIsValidator := func(addr string) bool {
		t.Helper()
		res, err := rpcClient.ABCIQuery(context.Background(),
			"vm/qeval",
			[]byte(fmt.Sprintf("gno.land/r/sys/validators/v2.IsValidator(address(\"%s\"))", addr)),
		)
		require.NoError(t, err)
		return strings.Contains(string(res.Response.Data), "true")
	}

	const (
		valAddrX = "g1td0cgmt9uz7kq4hcv7fkkwvp3z4lq4dsewffwr"
		valPubX  = "gpub1pgfj7ard9eg82cjtv4u4xetrwqer2dntxyfzxz3pqg0lte6srklm3tuyja9489n3dsnx4wcadq43wrwnz6nln8s7lf9uyptc3nm"
		valAddrY = "g1frduc0f2zw985cjqtmtml92mtxvpg2mgkzvxts"
		valPubY  = "gpub1pgfj7ard9eg82cjtv4u4xetrwqer2dntxyfzxz3pqgeaw3fky4c37fx53fm5kfptn6xlcdzxyd4ttu74stjch5fxmc9h5d9jmc6"
	)

	// Step 0: Init GovDAO
	t.Log("\n--- Init GovDAO ---")
	broadcastRun("Init GovDAO", `package main
import i "gno.land/r/gov/dao/v3/init"
func main() {
	i.InitWithUsers(address("g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"))
}`)

	// Step 1: Register both valopers
	t.Log("\n--- Register valopers ---")
	broadcastCall("Register valoper X",
		"gno.land/r/gnops/valopers", "Register",
		[]string{"valX", "Validator X", "on-prem", valAddrX, valPubX},
		std.MustParseCoins("20000000ugnot"),
	)
	broadcastCall("Register valoper Y",
		"gno.land/r/gnops/valopers", "Register",
		[]string{"valY", "Validator Y", "on-prem", valAddrY, valPubY},
		std.MustParseCoins("20000000ugnot"),
	)

	// Step 2: Add val X — block N
	t.Log("\n--- Block N: Add validator X ---")
	broadcastRun("Add validator X", `package main
import (
	"gno.land/r/gnops/valopers/proposal"
	"gno.land/r/gov/dao"
)
func main() {
	pr := proposal.NewValidatorProposalRequest(cross, address("`+valAddrX+`"))
	pid := dao.MustCreateProposal(cross, pr)
	dao.MustVoteOnProposalSimple(cross, int64(pid), "YES")
	dao.ExecuteProposal(cross, pid)
}`)
	require.True(t, queryIsValidator(valAddrX), "val X should be in set after ADD")

	// Step 3: Advance blocks (N+1, N+2)
	t.Log("\n--- Advance blocks ---")
	broadcastRun("no-op N+1", `package main
func main() {}`)
	broadcastRun("no-op N+2", `package main
func main() {}`)

	// Step 4: Add val Y — block N+3
	t.Log("\n--- Block N+3: Add validator Y ---")
	broadcastRun("Add validator Y", `package main
import (
	"gno.land/r/gnops/valopers/proposal"
	"gno.land/r/gov/dao"
)
func main() {
	pr := proposal.NewValidatorProposalRequest(cross, address("`+valAddrY+`"))
	pid := dao.MustCreateProposal(cross, pr)
	dao.MustVoteOnProposalSimple(cross, int64(pid), "YES")
	dao.ExecuteProposal(cross, pid)
}`)
	require.True(t, queryIsValidator(valAddrY), "val Y should be in set after ADD")

	// Step 5: Remove val X — block N+4
	t.Log("\n--- Block N+4: Remove validator X ---")
	broadcastRun("Remove validator X", `package main
import (
	valopers "gno.land/r/gnops/valopers"
	"gno.land/r/gnops/valopers/proposal"
	"gno.land/r/gov/dao"
)
func main() {
	addr := address("`+valAddrX+`")
	valopers.UpdateKeepRunning(cross, addr, false)
	pr := proposal.NewValidatorProposalRequest(cross, addr)
	pid := dao.MustCreateProposal(cross, pr)
	dao.MustVoteOnProposalSimple(cross, int64(pid), "YES")
	dao.ExecuteProposal(cross, pid)
}`)

	// Step 6: Liveness check
	t.Log("\n--- Liveness check ---")
	broadcastRun("Liveness check", `package main
func main() { println("alive") }`)

	// Step 7: Verify final state
	require.False(t, queryIsValidator(valAddrX), "val X should be removed")
	require.True(t, queryIsValidator(valAddrY), "val Y should still be in set")

	status, err := rpcClient.Status(context.Background(), nil)
	require.NoError(t, err)
	t.Logf("  Node alive at block %d", status.SyncInfo.LatestBlockHeight)
	t.Log("\n✓ No crash — test11 scenario passed!")
}
