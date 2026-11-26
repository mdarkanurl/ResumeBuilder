# API

Base URL in development: `http://localhost:8080` (or `/api` through the Vite
proxy).

## Health check

```http
GET /healthz
```

```http
200 OK
ok
```

## Upload a resume

```http
POST /upload
Content-Type: multipart/form-data
```

| Field  | Type | Required | Description        |
| ------ | ---- | -------- | ------------------ |
| `file` | file | yes      | PDF resume, ≤ 5 MiB |

### Success

```http
201 Created
Content-Type: application/json
```

```json
{
  "filename": "resume.pdf",
  "size": 248193
}
```

### Errors

| Status | Reason                                                    |
| ------ | --------------------------------------------------------- |
| 400    | Missing file, non-PDF content, invalid filename           |
| 400    | Not multipart, missing boundary, malformed form           |
| 413    | File exceeds `MAX_UPLOAD_SIZE`                            |
| 500    | Storage failure                                           |

Error bodies are plain text.

## Kafka events

Not implemented yet. The intended event is published to `KAFKA_TOPIC`
(`resume-uploaded`) after a successful upload.
