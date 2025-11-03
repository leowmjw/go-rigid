# AGENTS.md — Implementing “Rama-in-Go” (RIG)

- Go 1.25, stdlib HTTP + slog, Temporal workflows/activities.
- Mirror Rama REST (`/rest/{module}/depot/*/append`, `/rest/{module}/pstate/$$/select[One]`, `/rest/{module}/query/{name}/invoke`).
- TDD: tutorials 1–6 + testing guidance; use Temporal testsuite.
- Anonymous function replacement to decouple side effects.
- Deterministic workflows; non-determinism in activities.
