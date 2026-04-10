# ADR: Recover Validator Changes After Node Restart

## Context

The `EndBlocker` in `gno.land/pkg/gnoland/app.go` relies on an in-memory event
collector to decide whether to query the VM for validator set changes. The
collector listens on the `EventSwitch` for `validatorUpdate` events fired during
transaction execution. When events are present, the `EndBlocker` calls
`GetChanges(from, to)` on `r/sys/validators/v2` and forwards the resulting
updates to Tendermint2's consensus layer.

**The bug:** on node restart, the in-memory event collector is empty. Events
from the last block before shutdown are lost. The `EndBlocker` sees an empty
collector, returns early, and never queries `GetChanges` — validator changes
that were committed to the realm but not yet applied to consensus are
permanently lost.

This was confirmed by an integration test: a validator added via GovDAO proposal
and verified in the realm (`IsValidator` returns `true`) disappears from the
consensus set after a restart.

## Decision

Use a `firstBlock` flag inside the `EndBlocker` closure. On the very first
invocation after startup, the EndBlocker bypasses the collector check and always
queries the VM for pending validator changes. On all subsequent blocks, the
collector gate applies as before.

```go
firstBlock := true

return func(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
    // ... auth/gas price logic ...

    if firstBlock {
        firstBlock = false
        collector.getEvents() // drain any accumulated events
        if app.LastBlockHeight() == 0 {
            return abci.ResponseEndBlock{} // genesis — nothing to recover
        }
    } else if len(collector.getEvents()) == 0 {
        return abci.ResponseEndBlock{}
    }

    // ... VM query + apply validator changes ...
}
```

This keeps all validator logic in one place (the EndBlocker) rather than
splitting it between init and EndBlocker. The `firstBlock` bool is safe because
the EndBlocker runs single-threaded in the ABCI consensus flow.

The collector itself is kept as a performance optimization — it avoids a VM
query on every block when no validator changes occurred.

## Alternatives Considered

### 1. Query VM at init time and seed the collector

After `vmk.Initialize()`, query the VM for pending changes and pre-populate
the collector with a synthetic event so the EndBlocker picks it up naturally.

**Rejected:** splits validator logic across init and EndBlocker. Also requires
handling the case where the validators realm isn't deployed (test nodes), adding
error-handling complexity to init.

### 2. Fire a synthetic event via `evsw.FireEvent`

Seed the collector by firing a fake `validatorUpdate` event on the `EventSwitch`
after restart.

**Rejected:** `FireEvent` dispatches to all listeners, not just the collector.
Unknown listeners could react to a synthetic event in unexpected ways, creating
subtle bugs.

### 3. Always query the VM (remove the collector)

Remove the event collector entirely and call `GetChanges` in every `EndBlocker`.

**Rejected after benchmarking:** the collector early-return path costs ~13ns/0
allocs per block, while a VM query costs ~830ns/8 allocs (mock) and
significantly more with a real VM. Over thousands of blocks with no validator
changes, the collector avoids substantial overhead.

### 4. Persist the collector to disk

Save collector state before shutdown and restore on restart.

**Rejected:** adds persistence complexity for a problem that only manifests at
the boundary between the last pre-shutdown block and the first post-restart
block. A one-time unconditional query on the first block is simpler.

## Consequences

- **Positive:** validator changes committed to the realm are always applied to
  consensus, even across restarts.
- **Positive:** all validator change logic stays in the EndBlocker — no init-time
  special cases.
- **Positive:** the event collector optimization is preserved for normal
  operation (all blocks after the first).
- **Positive:** on genesis (height 0), the first EndBlocker still early-returns
  — no wasted VM query.
- **Trade-off:** if `GetChanges` or the validators realm is unavailable on the
  first block after restart, the EndBlocker logs an error and continues (same as
  any other block). No panic, no silent data loss.
- **Testing:** a txtar integration test (`restart_validators.txtar`) verifies
  the full flow: add validator via GovDAO, restart, confirm validator appears in
  the consensus set. The test also required adding `gnorpc validators` and
  `sleep` testscript commands.
