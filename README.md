# docintel

Document intelligence platform — upload documents, ask questions, get answers.

## Getting started

```bash
cp .env.example .env
# fill in API keys (GROQ_API_KEY, JINA_API_KEY, JWT_SECRET at minimum)

docker compose -f deploy/docker-compose.yml up -d postgres redis
go run ./cmd/migrate migrate
go run ./cmd/api
```

## Database

### Connection

The app reads connection details from env vars. `DATABASE_URL` takes precedence; individual vars are used as fallback.

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | — | Full Postgres DSN (used when set, e.g. on Fly.io) |
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `postgres` | Postgres user |
| `DB_PASSWORD` | `postgres` | Postgres password |
| `DB_NAME` | `docintel` | Database name |

The local docker-compose postgres uses `docintel/docintel/docintel`, matching the default `DATABASE_URL` in `.env.example`.

### Migrations

Migration files live in `scripts/db/migrations/` as paired `.up.sql` / `.down.sql` files, using [golang-migrate](https://github.com/golang-migrate/migrate) with the `file://` source.

```bash
# apply all pending migrations
go run ./cmd/migrate migrate

# roll everything back
go run ./cmd/migrate migrate down   # (not yet wired — run manually if needed)

# scaffold a new migration pair
go run ./cmd/migrate create -name <migration_name>
# e.g. go run ./cmd/migrate create -name add_documents_table
# → scripts/db/migrations/20260527123456_add_documents_table.up.sql
# → scripts/db/migrations/20260527123456_add_documents_table.down.sql
```

> **Note:** always run `go run ./cmd/migrate` from the project root. The tool resolves the migrations directory relative to `os.Getwd()`.

### Seeds

Seed files live in `scripts/db/seeds/`. There is no automated seed runner yet — execute them manually with `psql` when needed:

```bash
psql "$DATABASE_URL" -f scripts/db/seeds/<file>.sql
```

## Structure

| Path | Purpose |
|------|---------|
| `cmd/api` | API server entrypoint |
| `cmd/worker` | Background worker entrypoint |
| `cmd/migrate` | Migration CLI (`migrate` / `create`) |
| `internal/domain` | Core entities and port interfaces |
| `internal/app` | Use cases / application services |
| `internal/adapters/postgres` | Postgres repo implementations |
| `internal/adapters/redis` | Cache and rate-limit adapters |
| `internal/adapters/storage` | File storage (local + S3) |
| `internal/adapters/llm` | Groq client |
| `internal/adapters/embeddings` | Jina client |
| `internal/adapters/queue` | asynq task queue |
| `internal/adapters/pdf` | PDF text extraction |
| `internal/transport/http` | HTTP server, handlers, middleware |
| `internal/transport/worker` | asynq task handlers |
| `pkg/config` | DB config and connection |
| `pkg/migrate` | Migration runner and file scaffolder |
| `scripts/db/migrations` | SQL migration files |
| `scripts/db/seeds` | SQL seed files |
| `deploy/` | Dockerfiles and docker-compose |
