# TESTING.md — TDD plan for RIG

Enable tests: `go test -tags tdd ./...`

Suites:
- T1 Hello (ack map) — Rama REST + Tutorial 1
- T2 Word Count — Tutorial 2, Paths, PStates
- T3 Partitioners — Tutorial 3, Depots
- T4 Dataflow — Paths + transforms
- T5 ETLs — Microbatch (exactly-once), Stream retries
- T6 Social — Follow/Unfollow, subindexing, mirrors

See TODO.md for phased work.
