# AGENTS.md — hillviewCoreAPI

## 1. What this repo is

The back-office API behind **team.hillview.tv**: videos, playlists, the homepage spotlight, equipment assets and checkouts, short links, staff users, mobile users, and file uploads to S3. Every `/admin` route requires a Hillview access token; student accounts are limited to `StudentAllowedRoutes` (videos, playlists, upload). It does **not** issue tokens (hillviewAuthAPI does) or serve the public site's catalogue (hillviewVideoAPI does).

Served at `https://api.hillview.tv/core/v1.1` (NPM path routing on the shared `api.hillview.tv` host). Part of the Hillview platform — legacy codebase lineage; see §9 before copying any of its patterns elsewhere.

## 2. Stack & dependencies

- Go 1.26.9, `net/http` + `github.com/gorilla/mux`
- `github.com/Masterminds/squirrel` + `github.com/go-sql-driver/mysql` (MariaDB, no ORM)
- `github.com/gorilla/handlers` for CORS — a legacy deviation from the estate standard (`rs/cors`); leave it unless doing a deliberate migration
- `github.com/aidenappl/go-monitor` **v0.0.9** — telemetry to Monitor zone `hillview` (§6)
- Also: AWS SDK v2 (S3 uploads), SendGrid (mail), golang-jwt v5, google/uuid

## 3. Project structure

```
hillviewCoreAPI/
  main.go                # router, middleware order, background jobs, http.Server + graceful shutdown
  env/env.go             # every environment variable (getEnv / getEnvOrPanic)
  background/            # healthchecks.io heartbeat loop
  db/                    # `DB` (hillviewtv) and `AssetDB` (hillview_assets) connections
  docs/                  # API reference (`docs/API.md`)
  jwt/                   # Hillview JWT parsing/validation (`HVJwtClaims`)
  mailer/                # SendGrid mail
  middleware/            # request id, logging, headers, `TokenHandlers` (request log), `AccessTokenMiddleware`
  query/                 # squirrel queries
  responder/             # success/error JSON; `SendError` records the failure cause for Monitor
  routers/               # handlers (`handle*.router.go`)
  structs/               # entities
  util/                  # helpers
  telemetry/             # Monitor client: Init, request middleware, failure recording, Go(), shutdown flush.
                         #   IDENTICAL in all five Hillview Go APIs — change it in all of them
  Dockerfile             # golang:alpine build → alpine run
  Devfile.yaml           # dev CLI commands
  .github/workflows/push-to-registry.yml   # main → registry.appleby.cloud → Lattice redeploy
```

## 4. Running, building & testing

```bash
dev          # go run .
dev build    # go build -o bin/app .
dev test     # go test ./...
dev check    # gofmt -w -s . && go vet ./... && go test ./...
```

`env` panics at startup on any missing required variable:

| Variable | Required | Purpose |
|---|---|---|
| `DATABASE_DSN` | yes | MariaDB DSN for `hillviewtv` |
| `ASSET_DATABASE_DSN` | yes | MariaDB DSN for `hillview_assets` |
| `JWT_SIGNING_KEY` | yes | Shared Hillview JWT key |
| `SENDGRID_API_KEY` | yes | SendGrid |
| `AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / AWS_REGION` | yes (SDK) | S3 uploads (`us-west-1`) |
| `PORT` | no | Listen port (default 8000; Lattice sets 8001–8005) |
| `HEALTH_CHECK_URL` | yes | healthchecks.io ping URL, hit once a minute |
| `MONITOR_INGEST_URL` | no | Monitor ingest, `https://hillview-monitor-api.appleby.cloud/v1/events`. Unset = no telemetry |
| `MONITOR_API_KEY` | no | Ingest-scoped key minted **on the hillview zone** |
| `MONITOR_ENV` | no | Event `env` (default `production`) |
| `MONITOR_SPOOL_DIR` | no | Durable spool directory; events survive a Monitor outage or restart |
| `MONITOR_STDOUT` | no | `true` also prints events to stdout (default off — the container log keeps its own lines) |

Local runs need a reachable MariaDB. Leave the `MONITOR_*` variables unset locally unless you are testing telemetry — and then set `MONITOR_ENV` to something other than `production`.

## 5. How code is written here

- Handlers call `query` directly; there is no service layer.
- Errors go through `responder.SendError` where a responder exists. It calls `telemetry.RecordFailureW(w, message, err)` so the request's Monitor event carries the cause — **always pass the underlying `err`**, or a 5xx reaches Monitor with a message and no cause.
- Code that writes errors another way (`http.Error`, hand-written JSON) still gets a cause: the middleware keeps the first 512 bytes of any 5xx body.
- Every goroutine starts through `telemetry.Go(ctx, name, fn)` so a panic is reported instead of killing the process.
- Background failures are reported with `telemetry.Warn` / `telemetry.Error` using `{resource}.{action}.failed` names, alongside the existing `log.Printf`.
- No tests exist outside `telemetry/` (table-driven, standard `testing`, `monitor.StartRecording`).

## 6. Domain & architecture

**Request flow:** Cloudflare → NPM (`api.hillview.tv`, path location `/core/v1.1`) → worker-3 `10.15.10.75:8001` → middleware → handler → MariaDB.

**Middleware order** on the `/core/v1.1` subrouter: (request id, where present) → `telemetry.Middleware` → `LoggingMiddleware` → `MuxHeaderMiddleware` → `TokenHandlers` (where present) → route. `/healthcheck` is on the root router, outside all of it, and emits nothing.

**Auth:** Hillview JWTs minted by hillviewAuthAPI and signed with the shared `JWT_SIGNING_KEY`, read from `Authorization: Bearer` or the `hv_access` cookie. `AccessTokenMiddleware` calls `telemetry.SetUserID` once it resolves the user.

### Monitor telemetry

Service name **`hillview-core-api`**, zone **`hillview`** (`https://hillview-monitor-api.appleby.cloud`, Lattice stack 37). go-monitor checks at startup that the ingest URL really is the hillview zone. Monitor is never a boot requirement: with `MONITOR_INGEST_URL` unset the service runs as before, and with Monitor down events wait in `MONITOR_SPOOL_DIR`. SIGTERM drains HTTP for up to 10s, then flushes telemetry before exit.

Event catalogue:

| Event | Level | Meaning |
|---|---|---|
| `service.startup / service.stopping / service.shutdown` | info | process lifecycle; `stopping` carries the signal |
| `service.listen.failed` | fatal | the HTTP listener could not start |
| `service.shutdown.drain_incomplete` | warn | in-flight requests outlived the 10s drain |
| `http.request.end` | info / warn / error | one per request under the API prefix: `method`, `path` (mux route template — what Monitor groups issues by), `request_path`, `status_code`, `duration_ms`, `client_ip`, `user_agent`, `user_id` once the access token resolves, and on failure `error_message` / `error` / `error_type`. 5xx = error (files an issue), 401/404/405 = info, other 4xx = warn |
| `http.panic` | error | a handler panicked; the client gets a 500 |
| `healthcheck.ping.failed / healthcheck.ping.rejected` | warn | the healthchecks.io heartbeat could not be delivered |

`X-Request-ID` is set on every response under the prefix and exposed via CORS, so the web apps' Monitor events can be joined to this service's.

## 7. Ecosystem & related repos

- **Called by:** `hillview-team-web` (team.hillview.tv), `hillview-tv-web` (hillview.tv), the other Hillview APIs' clients.
- **Data:** MariaDB `hillviewtv` / `hillview_assets` on a host no MariaDB MCP reaches — there is no SQL access from agents.
- **Siblings:** hillviewCoreAPI, hillviewVideoAPI, hillviewAuthAPI, hillviewLinksAPI, hillviewAssetsAPI — same layout, same `telemetry/` package.
- **Monitor:** zone `hillview`; agents query it with `mcp__monitor_hillview__*`.

## 8. Operations

- **Deployed on Lattice stack 19 "Hillview API"** (worker-3, `10.15.10.75`), container `hillview-core-api`, host port `8001`. Stack env holds the secrets plus `MONITOR_HILLVIEW_INGEST_URL` / `MONITOR_HILLVIEW_INGEST_KEY`; the container maps them to `MONITOR_INGEST_URL` / `MONITOR_API_KEY` and mounts a named volume at `/var/lib/monitor-spool`.
- **CI:** push to `main` → `push-to-registry.yml` builds `registry.appleby.cloud/hillview-core-api:latest` → `LATTICE_DEPLOY_URL?container=…` redeploys this container. Env or compose changes need a **stack Deploy** in Lattice — a container redeploy reuses the old config.
- **Logs:** `lattice_get_container_logs`; **errors:** Monitor zone hillview (`monitor_list_issues` with `service: hillview-core-api`).
- **Health:** `GET /healthcheck` (Docker healthcheck) and the healthchecks.io heartbeat. `https://api.hillview.tv/` itself has no upstream — probe a prefixed route.

## 9. Rules & guardrails

- Never point `MONITOR_API_KEY` at a key from another zone or the control plane — events would file under someone else's tenant silently.
- Never add `ReadTimeout`/`WriteTimeout` to the server without exempting uploads; only `ReadHeaderTimeout` is set on purpose.
- Keep `telemetry/` identical across the five Hillview APIs.
- This is a legacy codebase: don't use its patterns (gorilla/handlers CORS, `interface{}` responders, `log.Printf` logging) as precedent for new repos.

## 10. Verification

```bash
gofmt -l .        # must print nothing
go build ./...
go vet ./...
go test ./...
```

## 11. Keeping this file updated

Any change to structure, stack, commands, conventions, endpoints, environment variables, events, or deployment updates this file (and `README.md` where setup or role changed) in the same change.
