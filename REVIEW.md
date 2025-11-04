# RIG TDD Test Review

Date: 2025-11-03
Scope: Review of current test scaffolding vs stated goals in AGENTS.md / TESTING.md (Tutorials 1–6) for Go + Temporal implementation of Rama concepts.

## Executive Summary
Current tests are skeletal placeholders that execute no real assertions and exercise no implemented logic (all core functions return ErrNotImplemented). Immediate priority is to define observable contracts per subsystem (Depot, PState, Path, Partitioner, Microbatch, Query, Agg) and encode them as failing, focused tests before adding implementation. Emphasize Temporal workflow determinism, side‑effect isolation, and exactly-once semantics early to avoid architectural rework.

## Critical Checklist (Highest Priority First)
1. Depot Append / Ack Levels
   - Define expected behaviors for Ack, AppendAck, None (timing of client completion, guarantees delivered, idempotency requirements).
   - Tests: success path, unknown depot/module, partition routing, duplicate append (same payload / key), ordering guarantees within a partition, cross-partition ordering (explicitly NOT guaranteed?), append during workflow replay.
   - Temporal specifics: signal vs activity usage, retry policies, failure escalation, heartbeat handling.
2. PState Select / Transform
   - Round‑trip transforms (set, merge, conditional update) and visibility through select.
   - Subindexing semantics (multi-level path, missing intermediate creation, nil vs absent distinction).
   - Concurrency: two concurrent transform operations (simulate via workflow local state vs activity) – conflict resolution / last-write-wins vs CAS.
3. Path Encoding/Decoding
   - Canonical JSON encoding order and stability; rejecting malformed sequences.
   - Coverage for each navigator type: implicit key, explicit must, pkey, multiPath, filterFunc; invalid token detection.
   - Property: Encode(Decode(x)) == canonical(x); fuzz decoding robustness.
4. Microbatch Exactly-Once
   - Coordinator & Attempt workflows: dedupe across retries, partial failure recovery, commit vs visible boundary (when results become observable downstream).
   - Idempotent side-effects (activity simulation) and WAL/outbox placeholder tests.
   - Backoff, retry exhaustion, metrics hooks existence (even if stubbed).
5. Partitioner
   - Deterministic hashing given same key & n.
   - Distribution quality (statistical test) for HashBy; behavior when n <= 0 or n=1.
   - Random partitioner still deterministic inside workflows (inject seed / function) – ensure tests enforce determinism requirement.
6. Query Invocation
   - Registry mechanism: unknown query name error, argument count/type validation, single vs multi return forms, panic isolation.
   - Determinism: queries must not mutate workflow state (tests that mutation attempts are rejected or no-oped).
7. Aggregations (Count)
   - Nil initial state becomes 0; increment semantics; large counts / overflow boundary (e.g., int64 wrap detection or use big.Int?).
8. Temporal Determinism Infrastructure
   - Tests that replay of recorded history (using testsuite) yields identical decisions (e.g., run once, capture commands, run again).
   - Guardrails: non-deterministic calls (time.Now, rand) must be wrapped; tests fail if direct usage detected (can implement lint/test stub later).
9. Error Taxonomy
   - Standardized error types (not implemented, invalid input, transient vs permanent) – tests asserting classification mapping to Temporal retries.
10. Performance / Load (Later Phase)
   - Benchmark skeletons for high-volume appends, path transformations, partition selection – maintain determinism.

## Observed Gaps Per Package
### depot
- No tests for append result structure (Ack map), ack level differences, partition workflow execution, error paths, or Temporal retry behavior.
- Missing tests for idempotency and ordering semantics.

### pstate
- Single test only checks that API calls compile; no assertions, no error case coverage (missing path, invalid op), no concurrency scenario.

### path
- Round-trip test lacks assertions; does not inspect encoded bytes or ensure expected navigators.
- Missing negative tests (malformed JSON arrays, unknown navigator tags, empty path handling) and property/fuzz tests.

### partition
- Determinism test lacks assertion; no tests for distribution, edge cases (n<=0), or behavior differences between Random and HashBy.

### microbatch
- Placeholder only; lacks attempt sequencing, dedupe across retries, failure injection, idempotent activity simulation, temporal timers & backoff.

### query
- Single call without checking result or error; no tests for registry, unknown name, argument validation, concurrency, determinism.

### agg
- Count test does not assert final value (expected 3). No tests for nil initial state semantics or large iteration counts.

### testutil
- Provided helper (NewEnv) untested for deterministic configuration (e.g., time skipping, interceptors). Consider tests ensuring environment customization (e.g., disabling deadlock detection configuration if needed) can be plugged in.

## Recommended Test Design Enhancements
- Use table-driven tests with explicit expected outputs & errors.
- Introduce golden tests for path JSON encoding to lock format.
- Add fuzz tests (Go 1.18+ built-in) for path DecodeJSONPath and partition hashing with random keys.
- Temporal workflow tests: use env.RegisterWorkflow / env.RegisterActivity and env.ExecuteWorkflow with arguments; assert env.GetWorkflowResult & env.GetWorkflowError; add tests for replay by re-running after env.Close(); test signal handling when scaling partitions.
- Use dependency injection / function variables to replace non-deterministic calls (timeNow, randInt) – test that replacement mechanism works.
- Add race condition tests with -race (CI) once implementation present (especially pstate store and aggregator updates).
- Consider property-based tests: CountFold(n times) == n; Partition distribution (Chi-squared) for large sample.

## Phasing Alignment (TODO.md)
- Phase 0 (Foundations): Flesh out failing tests for API contracts BEFORE code; current state incomplete – expand tests per checklist items 1–4 minimally.
- Phase 1 (Depots & Streams): Need concrete append/ack tests & partition routing.
- Phase 2 (PStates & Paths): Strengthen path encode/decode & pstate transform semantics tests.
- Phase 3 (Microbatch): Add exactly-once + dedupe tests early; they influence storage/ack design.
- Later phases (Query, Mirrors, Hardening): Prepare test placeholders with TODO comments enumerating contract assertions to prevent scope drift.

## Temporal-Specific Considerations
- Ensure tests cover: workflow versioning (Version API usage), activity retry classification (non-retryable errors), heartbeat timeouts, cancellation propagation (client cancel vs internal abort), workflow time skipping & timers, signal/child workflow interactions.
- Provide tests that detect accidental non-determinism (e.g., injecting random seed difference) by capturing history and re-executing with a deterministic replayer.

## Metrics & Observability (Prepare Early)
- Plan tests that assert logging (slog) structure or metric counters (use in-memory sink) for append counts, retry counts, dedupe occurrences.

## Open Clarifications Needed
1. Precise semantics for Ack vs AppendAck vs None (client latency vs durability vs visibility guarantees)?
2. Expected ordering guarantees across partitions (is global ordering ever required or explicitly avoided)?
3. Partitioning determinism: Should Random be deterministic inside workflows via a fixed seed or replaced by a hash? (Temporal determinism constraint.)
4. PState transform conflict policy: last-write-wins, compare-and-swap, or transactional batching semantics?
5. Path navigator specification: definitive grammar / reserved tokens list? (e.g., multiPath, must, pkey, filter function prefix.)
6. Microbatch commit protocol: two-phase (prepare/commit) or single-phase with idempotent externalization? Dedupe key basis (batchID, attemptID, combination)?
7. Aggregation state type expectations: Count uses int, int64, or arbitrary precision? Overflow handling strategy?
8. Query side-effect rules: are queries strictly read-only (Temporal definition) or allowed controlled side-effects via activities?
9. Error taxonomy: desired classification mapping to Temporal (application vs internal vs non-retryable) – formal list?
10. Performance targets (append throughput, latency) to guide benchmark thresholds & regression guards?

## Next Actions
- Author failing tests per checklist (focus on Depot ack semantics + Path encoding + Count aggregator correctness first).
- Draft minimal domain contracts in code comments to solidify test expectations.
- Introduce deterministic dependency injection shims for time/rand before implementing logic.

---
Please confirm or clarify the Open Clarifications list so test authoring can proceed accurately.
