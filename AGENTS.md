# AGENTS.md — Implementing “Rama-in-Go” (RIG)

<updated: 2025-11-04T16:46:49Z>

Vision: Incrementally reproduce Rama primitives (Depots, Streams, PStates, Paths, Queries, Microbatch, Mirrors) in Go with Temporal ensuring deterministic workflow logic and activity-isolated side effects.

Core Stack:
- Go 1.25, stdlib HTTP + slog.
- Temporal workflows (pure/deterministic) + activities (non-deterministic I/O, idempotent side effects).

REST Surface (to mirror Rama):
- Depot append: `POST /rest/{module}/depot/*/append` (ack | appendAck | none).
- PState select/selectOne: `POST /rest/{module}/pstate/$$/{select|selectOne}`.
- Query invoke: `POST /rest/{module}/query/{name}/invoke`.

Progress Summary:
- Initial HTTP depot append handler with ack level variants + test.
- Count aggregation implemented (spec-driven start).
- Partition hashing deterministic + error handling (n<=0) tested.
- Temporal test harness placeholders (workflows not yet implemented).
- Deterministic test stubs (time/rand/uuid) added for injection.
- Review authored outlining gaps and priority checklist.

Outstanding Core Work (next phases):
1. Implement navigator types + JSON codec (Key/Must/PKey/MultiPath/FilterFunc) minimally; add negative tests.
2. In-memory PState store (subindexing, LocalTransform/LocalSelect atomic semantics) + path application tests.
3. Depot append workflow & activity design (ledger + idempotent appendAck path); determinism/replay tests.
4. Microbatch exactly-once via idempotent commit activity (batch ledger) + fallback child workflow if complexity grows.
5. Stream processor workflow with retry/backoff and ordering guarantees; mock activity for events.
6. Query registry + Invoke tests (argument validation, pure function enforcement) using package-local sentinel errors.
7. Determinism validation tests (history length equality between replays) for each workflow.
8. Negative & edge tests (bad JSON path tokens, invalid ackLevel, deep path nesting, large partition counts).
9. Batch atomicity tests for PState (failure injection rollback) + multi-path fanout.
10. Distribution test for Random partitioner (activity-scoped) + optional benchmark for HashBy and path codec.

Guiding Principles:
- TDD first: add failing, asserting tests before implementing each primitive.
- Keep workflows free of non-determinism: inject side-effect closures only inside activities.
- Incremental complexity: finish Depots + PStates/Paths before Microbatch/Streams/Queries.
- Package-local sentinel errors now; propose typed errors only when classification required.

Immediate Next Step:
- Begin Phase 2 foundations: add minimal navigator structs & round-trip tests driving path codec implementation.

Clarifications Incorporated:
- Minimal external test libs permitted.
- In-memory stores for unit tests; Pebble integration later.
- Minimal navigators now; expand later with TODO markers.
- Idempotent activities for microbatch; child workflows only if necessary.

TODO Markers (embed in code as // TODO(rig-phase-X)) will track phase progression.
