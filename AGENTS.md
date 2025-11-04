# AGENTS.md — Implementing “Rama-in-Go” (RIG)

## Overview
- Go 1.25, stdlib HTTP + slog, Temporal workflows/activities.
- Mirror Rama REST (`/rest/{module}/depot/*/append`, `/rest/{module}/pstate/$$/select[One]`, `/rest/{module}/query/{name}/invoke`).
- TDD: tutorials 1–10 + testing guidance; use Temporal testsuite.
- Anonymous function replacement to decouple side effects.
- Deterministic workflows; non-determinism in activities.

## Key Learnings from TDD Implementation
- **URL Encoding Handling**: REST paths use URL-encoded characters (e.g., `%2A` for `*`, `%24` for `$`). Parse `r.URL.Path` after decoding; `strings.Split` works post-decoding.
- **In-Memory Stores for Tests**: Use simple maps in `Server` struct for tutorial modules to simulate Rama logic without full persistence.
- **Assertion-Driven Tests**: Stub tests pass vacuously; add explicit assertions for responses, status codes, and expected values.
- **Temporal Test Setup**: Use `testutil.NewEnv()` for `TestWorkflowEnvironment`; register workflows/activities; `ExecuteWorkflow` then `GetWorkflowResult`.
- **Path Parsing**: Extract module, depot/pstate/query from URL paths using `strings.Split`; handle operations like `selectOne`.
- **Error Handling**: Return JSON errors with `writeErr` for unhandled cases; status 501 for not implemented.
- **Inside-Out Focus**: Start with core server handlers, add internal components (e.g., partition), then expand tests.

## Phase-by-Phase Implementation Guide

### Phase 0: Foundations (Tests Only)
- Implement TDD stubs for T1–T10 in `tutorials/`.
- Add assertions to T1–T6 for expected behaviors.
- Set up `testutil` with Temporal `TestWorkflowEnvironment`.
- Update `TESTING.md` with suite descriptions.
- Goal: All tutorial tests pass with basic in-memory logic.

### Phase 1: Depots & Streams
- Implement `internal/depot/` for append operations.
- Add stream processing workflows in `internal/stream/`.
- Replace in-memory depot processing with actual stream logic.
- Test depot append triggers workflows.

### Phase 2: PStates & Paths
- Implement `internal/pstate/` with `Store` interface and KV backend (e.g., pebble).
- Add path navigation in `internal/path/`.
- Update server `handlePStateSelect` to use real pstate queries.
- Ensure selectOne returns correct data from persistent store.

### Phase 3: Microbatch
- Implement `internal/microbatch/` for exactly-once processing.
- Use Temporal workflows for batching and deduplication.
- Add idempotent appends with keys/hashes.
- Test microbatch scenarios in T5/T10.

### Phase 4: Query
- Implement `internal/query/` for named queries.
- Add `handleQueryInvoke` in server.
- Support parameterized queries with results.
- Test query invocation in T7.

### Phase 5: Mirrors
- Implement mirroring for replication.
- Add subindexing and adjacency lists.
- Update social tutorial logic to real mirrors.

### Phase 6: Hardening
- Add error handling, retries, and logging.
- Implement anonymous function decoupling for side effects.
- Ensure workflow determinism and activity non-determinism.
- Performance tuning and integration tests.
