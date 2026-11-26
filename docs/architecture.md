# Architecture

## Overview

```
Browser (React) ──/api──▶ Vite dev proxy ──▶ Go API ──▶ uploads/
                                                   └──▶ Kafka (planned)
```

The web app and the API are independent applications in one repository. During
development the Vite dev server proxies `/api` to the Go server so the browser
only ever talks to one origin.

## API layers

| Layer          | Package               | Responsibility                                  |
| -------------- | --------------------- | ----------------------------------------------- |
| Entrypoint     | `cmd/server`          | Wiring, HTTP server, graceful shutdown          |
| Transport      | `internal/handler`    | Multipart parsing, status codes, JSON responses |
| Cross-cutting  | `internal/middleware` | Logging, panic recovery, CORS                   |
| Business rules | `internal/service`    | PDF validation, filename sanitisation           |
| Persistence    | `internal/repository` | Writing files to `UPLOAD_DIR`                   |
| Configuration  | `internal/config`     | Environment variables and defaults              |

Dependencies point inward: handlers depend on services, services depend on a
small repository interface, and the repository knows nothing about HTTP.

### Why no database package yet

The API currently stores resumes on the local filesystem. `internal/database`
and `migrations/` are reserved for the first real database integration so the
layout is ready when that happens.

## Adding a feature

1. Add the business rule to `internal/service`.
2. Add persistence to `internal/repository` and an interface next to the service
   that consumes it.
3. Expose it through `internal/handler`.
4. Register the route in `cmd/server/main.go`.
5. Add the client call in `apps/web/src/services` and the UI in
   `apps/web/src/pages`.

## Frontend

- `services/` is the only place that calls `fetch`.
- `hooks/` owns async state; components stay presentational where practical.
- `types/` mirrors API payloads.
- `VITE_API_BASE_URL` defaults to `/api`.
