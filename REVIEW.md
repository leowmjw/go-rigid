# RIG Test & Design Review (Initial Pass)

Date: 2025-11-03

## Executive Summary
Current tests are skeletal smoke placeholders. They do not assert behavior, invariants, determinism, error paths, or Temporal-specific semantics. Immediate need: convert placeholders into spec-driving TDD examples aligned with Phases 0–2 (Depots/Streams, PStates/Paths) before implementing later microbatch/query/mirrors.

## Critical Gaps (Highest Priority First)
1. Missing assertions everywhere (all *_test.go). Tests neither fail nor guide implementation.
2. No coverage of deterministic vs non-deterministic separation (workflows vs activities) for Temporal.
3. Absent Temporal workflow replay invariants (e.g., fixed decisions, no time.Now/use of randomness in workflows, stub activities).
4. No negative tests (invalid ackLevel, bad JSON path, invalid partition count n<=0, etc.).
5. No concurrency / race / idempotency checks for append semantics (ack vs appendAck vs none).
6. No exactly-once semantics model for microbatch (attempt workflow interaction, rollback/commit markers, durable batch ledger).
7. PState tests do not validate subindexing mechanics, path composition, transform atomicity, or error when selecting non-existent paths.
8. Path codec tests do not round-trip actual navigator types (only JSON arrays placeholders) and lack mismatch/unknown token tests.
9. Partition tests ignore determinism boundaries (HashBy must yield consistent mapping, Random should be activity-level only). No tests for distribution quality or error on n==0.
10. No tests for server HTTP handlers (status codes, JSON bodies, integration route matching).
11. No aggregation tests beyond Count; Count doesn't assert initialization & increment semantics or concurrency safety.
12. Query tests do not assert input validation, result shape, or error wrapping.
13. Missing table-driven tests for edge cases (empty depot name/module, large batch sizes, deep paths, large partition counts).
14. No tests for temporal retry policies, backoff, cancellation, heartbeats in long-running workflows.
15. Absence of benchmarks for hot path operations (path encode/decode, partition pick).

## Phase Alignment Recommendations
Phase 0 (Foundations):
- Introduce assertion framework (std testing + testify/require if allowed). Keep minimal if external deps undesired.
- Define canonical test fixtures for modules, depots, pstates.
- Provide deterministic stubs for time, UUID, randomness (inject via anonymous function replacement per AGENTS.md).

Phase 1 (Depots & Streams):
- Add tests for Append ack levels: ack returns non-empty metadata; appendAck waits for durability; none returns immediately. Simulate durability via fake activity.
- StreamProcessorWorkflow: test partition assignment, message ordering guarantees, retry semantics on activity failure.
- Depot append ordering & idempotency (appendAck should not duplicate on replay).

Phase 2 (PStates & Paths):
- Path navigator round-trip (Key, Must, PKey, MultiPath, FilterFunc). Test JSON schema strictness.
- LocalTransform atomicity (Batch vs single operations). Add failure injection.
- Subindexing & multi-path fanout selection; test concurrency (parallel transforms) and isolation.

## Temporal-Specific Checklist
- Use workflow test environment to mock activities (env.RegisterActivity).
- Assert history length invariants after workflow run (ensures deterministic decisions).
- Freeze now: inject workflow.Now override vs using time.Now directly.
- Replace random/UUID with injected closures captured outside workflow logic.
- Test workflow cancellation & timeout (context cancellation) semantics where long polling will exist.

## Detailed Recommendations Per Package
Agg:
- CountAgg.Fold: test initial nil state -> 1; repeated increments; concurrency (run folds in goroutines with channel serialization or confirm non-thread-safe by design).
- Add error path test (e.g., if state type mismatch).

Depot:
- AppendResult contract: define fields (topology, offsets, batchID). Tests should assert presence based on ack level.
- Workflow tests: simulate AppendWorkflow (deterministic creation of append ledger entry). Add test for replay (run twice, histories equal).
- Error tests: invalid depot/module, nil data, oversize payload.

Microbatch:
- CoordinatorWorkflow: test scheduling of AttemptWorkflow child workflows or activities.
- AttemptWorkflow: test exactly-once by simulating failure mid-batch then retry; assert committed set matches expected, duplicates suppressed.
- Introduce durable attempt ledger mock (in-memory) with injected activity for commit.

Partition:
- HashBy.Pick: assert deterministic mapping for same key & n; test n change effects; test panic/error for n<=0.
- Random: ensure test marks it as activity-scoped; run multiple picks and assert non-deterministic distribution using seeded PRNG in activity.

Path:
- Build concrete types for navigators; test encode/decode across each variant.
- Negative: unknown tag, malformed multiPath, empty must sequence.
- Deep nesting performance (optional bench).

PState:
- LocalTransform followed by LocalSelect returns updated value.
- Batch atomicity: partial failure rolls back entire batch (simulate failure injection).
- Subindexing: nested path updates vs independent keys.

Query:
- Invoke: test registered query function map; unregistered name returns error; argument count mismatch.
- Side-effect isolation: ensure queries are pure (no global mutation) via test harness.

Server:
- Handler tests: use httptest to send JSON requests, assert status & body for each ackLevel; invalid ackLevel returns 400 with detail.
- Path-based routing: ensure variable extraction for module names (use realistic modules).
- Add test for proper Content-Type headers.

Stream:
- StreamProcessorWorkflow: test message fetch loop using mocked activity providing messages; assert ordering & backoff on transient errors.
- Retry modes: configure different retry policies and assert number of attempts.

## Test Structuring Improvements
- Adopt table-driven patterns consistently.
- Provide helper: withTestEnv(t, func(env *testsuite.TestWorkflowEnvironment)) to reduce boilerplate.
- Centralize fakes in /internal/testutil (time, rand, UUID, storage).
- Use t.Helper() in assertion functions.

## Suggested New Test Files
- internal/server/http_test.go (REST handlers)
- internal/depot/api_test.go (ClientAppend ack behaviors)
- internal/microbatch/coordinator_test.go & attempt_test.go
- internal/stream/processor_test.go
- internal/path/path_negative_test.go
- internal/pstate/pstate_batch_test.go
- internal/query/query_error_test.go

## Risk Areas If Untested
- Non-determinism causing Temporal workflow replay failures leading to stuck tasks.
- Data duplication in microbatch attempts (exactly-once violation).
- Incorrect path encoding causing silent state corruption.
- Partitioning skew leading to hotspots.
- Missing error handling surfacing as 500s instead of typed domain errors.

## Minimal Initial Assertion Examples (Illustrative)
(Do not implement yet; outline only.)
- Depot: TestDepot_Append_AckLevels asserts result map non-empty for ack; empty for none; appendAck blocks until mock durable flag set.
- Microbatch: TestMicrobatch_ExactlyOnce_AcrossAttempts runs AttemptWorkflow twice with injected failure; committed set size == original unique messages.
- Path: TestPath_JSONRoundTrip_Navigators ensures EncodeJSONPath output equals original JSON (after canonical ordering) for each navigator.
- Partition: TestPartition_HashBy_Deterministic checks two consecutive picks produce same partition index and < n.
- PState: TestPState_LocalSelectAndTransform_Subindexing verifies transformed value equals expected and unrelated path unaffected.

## Checklist (Action-Ordered)
[ ] Define deterministic stubs (time, rand, uuid) & inject into workflows.
[ ] Flesh out navigator concrete types + path round-trip tests with assertions.
[ ] Add assertions to existing tests (no silent passes).
[ ] Implement server handler tests for depot append ackLevel variants.
[ ] Add negative tests for invalid input across modules.
[ ] Establish microbatch attempt ledger mock and tests for exactly-once.
[ ] Workflow determinism tests (history equality on replay).
[ ] Partition tests for error on n<=0 and distribution sanity.
[ ] PState batch atomicity tests with failure injection.
[ ] Query registration and error path tests.
[ ] Aggregation error + success tests (CountAgg increments).
[ ] Stream processor retry/backoff tests.

## Clarifications Needed
Please confirm:
1. Are external test libs (testify, go-cmp) acceptable or restrict to stdlib only?
2. Expected persistence layer for depots/pstates (Pebble implied) — is a real Pebble dependency planned or abstractions only for Phase 0?
3. Target scope for Phase 0 foundations: should we introduce full navigator types now or only test scaffolds?
4. Acceptable to add small testing helpers in internal/testutil now?
5. Do we model microbatch exactly-once via Temporal child workflows, activities with idempotent side-effect markers, or both?
6. Preferred error taxonomy (package-local Err* vs sentinel types)?

## Closing
Transform placeholders into behavior-driven specs before coding implementations; otherwise risk drifting from Rama semantics. Await clarifications to refine and prioritize.
