# docintel

Document intelligence platform — upload documents, ask questions, get answers.

## Getting started

```bash
cp .env.example .env
# fill in your API keys

docker compose -f deploy/docker-compose.yml up -d postgres redis
go run ./cmd/api
```

## Structure

| Path | Purpose |
|------|---------|
| `cmd/api` | API server entrypoint |
| `cmd/worker` | Background worker entrypoint |
| `internal/auth` | JWT, password hashing, middleware |
| `internal/config` | Environment loading |
| `internal/db` | Migrations and sqlc queries |
| `internal/documents` | Document service layer |
| `internal/chunks` | Chunking and embedding logic |
| `internal/chat` | Prompt construction and retrieval |
| `internal/conversations` | Conversation management |
| `internal/llm` | Groq client wrapper |
| `internal/embeddings` | Jina client wrapper |
| `internal/queue` | asynq task definitions |
| `internal/storage` | File storage abstraction |
| `internal/cache` | Redis cache helpers |
| `internal/ratelimit` | Rate limiting |
| `internal/observability` | Logging, metrics, tracing |
| `internal/http` | HTTP server, handlers, middleware |
| `deploy/` | Dockerfiles and docker-compose |
