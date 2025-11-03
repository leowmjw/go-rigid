# TODO.md — Phases & Future

## Phases
0) Foundations (tests only)
1) Depots & Streams
2) PStates & Paths
3) Microbatch
4) Query
5) Mirrors
6) Hardening

## Future Options
- Reactive PState diffs (SSE/gRPC/long-poll)
- Idempotent appends (keys/hash)
- Microbatch→external exactly-once (outbox/WAL)
- Custom navigators (allowlist/WASM)
- Advanced partitioners (ring/locality)
- Client sharding/redirects (discovery)
