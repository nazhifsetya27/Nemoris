# Manual Test Suite (Sprint 6)

Runnable test coverage for NEMORIS backend. Run from `backend/` directory.

## Quick Start

```bash
cd backend
go run ./cmd/tests/run_all.go
```

Results are saved to `backend/test_results/`.

## Individual Tests

| Test | Command | Notes |
|------|---------|-------|
| AI | `go run ./cmd/tests/test_ai.go` | No DB required |
| Reply builder | `go run ./cmd/tests/test_reply_builder.go` | No DB required |
| Reminder flow | `go run ./cmd/tests/test_reminder_flow.go` | Needs DB, .env |
| Repository | `go run ./cmd/tests/test_repository.go` | Needs DB |
| Scheduler | `go run ./cmd/tests/test_scheduler.go` | Needs DB; may take 2–6 min if WAHA slow |
| Health | `go run ./cmd/tests/test_health.go` | Needs DB, WAHA for full check |

## Prerequisites

- `.env` configured (DB_HOST, DB_USER, etc.)
- PostgreSQL running
- BOT_NUMBER set (for self-message test)

## Test Sender

All DB-modifying tests use `test-user@s.whatsapp.net` for isolation.
