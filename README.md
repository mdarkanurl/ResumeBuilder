# ResumeBuilder

Upload a PDF resume, and (eventually) get help turning it into a sharper one.

- **Web:** React 19 + TypeScript + Vite (`apps/web`)
- **API:** Go standard library (`apps/api`)
- **Infrastructure:** Kafka via Docker Compose

## Project layout

```
apps/
  web/                  React frontend
    src/
      assets/           Static assets imported by the app
      components/       Reusable UI components
      hooks/            React hooks
      pages/            Route-level pages
      services/         API calls
      types/            Shared TypeScript types
      utils/            Helpers
      App.tsx
      main.tsx
  api/                  Go backend
    cmd/server/         Application entrypoint
    internal/
      config/           Environment configuration
      handler/          HTTP handlers
      middleware/       Logging, recovery, CORS
      service/          Business logic
      repository/       Persistence
      database/         Reserved for database wiring
    migrations/         Reserved for SQL migrations
    uploads/            Local resume storage (gitignored)
docs/                   Architecture and API notes
scripts/                Development scripts
```

`pkg/` is intentionally omitted: nothing in the API is reusable outside the
module yet. Add it when a package genuinely needs to be imported externally.

## Prerequisites

- Go 1.27+
- Node.js 20+ and npm
- Docker (optional, for Kafka)

## Setup

```bash
cp .env.example .env      # optional; defaults work for local development
./scripts/setup.sh
```

## Run

```bash
./scripts/dev.sh
```

This starts:

- API on <http://localhost:8080>
- Web on <http://localhost:5173>

The Vite dev server proxies `/api/*` to the API, stripping the `/api` prefix.

### Run apps separately

```bash
cd apps/api && go run ./cmd/server
cd apps/web && npm run dev
```

### Kafka

```bash
docker compose up -d
```

The compose file creates the `resume-uploaded` topic. Publishing events is not
wired up yet; see `apps/api/internal/service/kafka.go`.

## Build

```bash
./scripts/build.sh
```

## API

| Method | Path       | Description                          |
| ------ | ---------- | ------------------------------------ |
| GET    | `/healthz` | Health check                         |
| POST   | `/upload`  | Multipart upload of a PDF (`file`)   |

```bash
curl -F "file=@resume.pdf" http://localhost:8080/upload
```

Limits and behaviour:

- PDF content is verified from the file header, not the extension.
- Maximum upload size defaults to 5 MiB (`MAX_UPLOAD_SIZE`).
- Files are stored in `UPLOAD_DIR` using the client-provided base name.

## Configuration

| Variable              | Default                 | Used by |
| --------------------- | ----------------------- | ------- |
| `PORT`                | `8080`                  | API     |
| `UPLOAD_DIR`          | `./uploads`             | API     |
| `MAX_UPLOAD_SIZE`     | `5242880`               | API     |
| `CORS_ALLOWED_ORIGIN` | `http://localhost:5173` | API     |
| `KAFKA_BROKERS`       | `localhost:9092`        | API     |
| `KAFKA_TOPIC`         | `resume-uploaded`       | API     |
| `VITE_API_BASE_URL`   | `/api`                  | Web     |

More detail: [docs/architecture.md](docs/architecture.md) and
[docs/api.md](docs/api.md).
