# Bilingual Tests

Real end-to-end tests for the message processing flow.

## Test Cases

| ID | Request | Description |
|----|---------|-------------|
| 1 | remind me to pay electricity tomorrow 8pm | English reminder with time |
| 2 | ingatkan saya untuk bayar listrik besok jam 8 malam | Indonesian reminder with time |
| 3 | list | English list reminders |
| 4 | daftar | Indonesian list reminders |
| 5 | ping | Ping |

## Run Tests

From project root (with Postgres running):

```bash
cd backend && go run ./cmd/bilingual_test
```

Results are written to `test/results.json`.
