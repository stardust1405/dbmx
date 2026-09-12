# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

Run from the repo root:

| Task | Command |
| --- | --- |
| Dev (hot-reload, opens window) | `wails dev` |
| Production build → `build/bin/` | `wails build` |
| Rebuild embedded SQLite template | `bash scripts/build_db.sh` |
| Verify Go compiles | `go build ./...` |

From `frontend/` (pnpm 10.24.0, pinned via `packageManager`):

| Task | Command |
| --- | --- |
| Type/template check | `pnpm check` |
| Lint (prettier + eslint) | `pnpm lint` |
| Auto-format | `pnpm format` |

`wails dev` also serves at `http://localhost:34115`, where Go methods can be called from browser devtools. Vite is on `:5173` with `strictPort`.

**There is almost no test suite.** No vitest/playwright, and the only Go tests are `app/table_schema_test.go`. Verification normally means `go build ./...`, `pnpm check`, and running the app.

`app/table_schema_test.go` covers the schema editor (add/edit/drop for columns, indexes and constraints) against a **real postgres**, because that code's entire output is SQL text a server has to accept — a type modifier that survives the round trip, a transaction that rolls a half-applied edit back, a constraint body rendered the way `pg_get_constraintdef` prints it. None of that can be checked by compiling.

It skips unless `DBMX_LIVE_PG` names a server, so `go test ./...` stays green without one:

```bash
docker run -d --name dbmx-pg -e POSTGRES_PASSWORD=dbmx -e POSTGRES_DB=dbmxtest \
  -p 55432:5432 postgres:16
DBMX_LIVE_PG='postgres://postgres:dbmx@localhost:55432/dbmxtest?sslmode=disable' \
  go test ./app -run TestLive -v
```

Point it at a scratch database only — the fixtures drop and recreate their tables. Reach for it when changing anything that generates SQL; `liveConnections` gives you a `*Connections` backed by a real pool and an in-memory sqlite `tabs` row, which is the same path the Wails frontend takes.

## Required local setup

Two git-ignored files must exist **at compile time**. Both are `go:embed`ed, so their absence is a build failure, not a runtime one:

1. **`config/env/env.toml`** — embedded by `config/env/env.go` (`//go:embed *.toml`). The README documents only two of the three sections:

   ```toml
   [sqlite3]
   name = "app.db"        # parsed into Env.Sqlite3 but never read; keep for parity

   [supabase]
   project_id = "..."
   anon_key = "..."

   [dbmx]
   base_url = "..."       # hosted DBMX API for Stardust AI — MISSING from the README
   ```

2. **`config/database/app.dev.db`** (dev) / **`config/database/app.db`** (`-tags production`). `scripts/build_db.sh` only produces `app.db`; for `wails dev` you must also copy it to `app.dev.db`.

## Architecture

### The Wails binding list is the entire API surface

`main.go` constructs every dependency and passes the bound structs in `options.App.Bind`: `App`, `Connections`, `Tabs`, `QueryHistory`, `SavedQueries`, `Auth`, `Stardust`. `PoolManager` is deliberately *not* bound — it is an internal dependency threaded into `Connections`, `Tabs`, and `Stardust`.

Any exported method on a bound struct becomes callable from TypeScript. Bindings are generated into `frontend/src/lib/wailsjs` (`wailsjsdir` in `wails.json`) and **are committed to git** — after changing a Go method signature, re-run `wails dev`/`wails build` and commit the regenerated files, or the frontend will silently call a stale signature.

Methods use pointer receivers, with one stray exception (`func (s Stardust) Chat`).

### Local SQLite is a build artifact, not a runtime migration

This is the most surprising part of the codebase:

1. goose runs at **build time** — `scripts/build_db.sh`, wired as a `preBuildHooks` entry in `wails.json` — applying `migrations/*.sql` into `config/database/app.db`.
2. That file is `go:embed`ed through a build-tag split: `database_embed_dev.go` (`//go:build !production`) embeds `app.dev.db`; `database_embed_prod.go` (`//go:build production`) embeds `app.db`. `wails.json` sets `"tags": "production"` for builds.
3. On first launch, `ensureDB` in `config/database/sqlite3.go` copies the embedded bytes to `os.UserConfigDir()/dbmx/app.db` — **only if that file does not already exist**.

Consequence: adding a migration does **not** migrate any existing database, including your own. After adding `migrations/000N_*.sql`, rebuild the template *and* delete `~/Library/Application Support/dbmx/app.db`.

### Postgres connection pooling

`app/postgres_manager.go` keys `*pgxpool.Pool` by a generated `uuid.UUID` **pool ID**, not by saved-connection ID. One saved connection can back many live pools (one per opened database), so `ActiveConns[connID]` refcounts them and `DeletePool` only forgets the connection when the last pool closes. Guarded by a `sync.RWMutex`.

A tab remembers its pool in `tabs.active_db_id`; `Stardust.Chat` and query execution both resolve the pool from that column. `App.shutdown` calls `Connections.TerminateAllDatabaseConnections`.

### Stardust AI

`app/stardust.go` is a thin proxy to a hosted service via `github.com/stardust1420/dbmx-go`, pointed at `Env.DBMXConfig.BaseURL` and authorized with the Supabase access token from `Auth.getToken()` (`app/auth.go`, which transparently refreshes from the `active_session` table when expired).

`Chat` fetches live schema DDL from the tab's active pool and prepends it as a system prompt — the remote model never touches the database directly.

### `model/` vs `schema/`

`model/` holds the structs mirroring local SQLite tables that cross the Wails boundary (they surface in `wailsjs/go/models.ts`). `schema/` is nearly empty. `dbmx.*` types come from the external client package.

### Frontend

- SvelteKit with `adapter-static`, and `src/routes/+layout.ts` sets `prerender = true` / `ssr = false`. This is required for Wails' embedded asset server — do not add server routes or server `load` functions.
- Svelte 5 runes. **`src/lib/state.svelte.ts` is the global store and deliberately mixes two systems**: `$state(new SvelteMap(...))` for `postgresConnectionsMap`, `databasesMap`, `tabsMap`, loading maps; and classic `writable()` stores for grid `columns`/`rows`/pagination, user/session flags, and editor prefs. Match whichever system a value already uses.
- UI is shadcn-svelte (`frontend/components.json`, slate base, aliased to `$lib/components/ui`) — add components with the shadcn-svelte CLI rather than hand-writing into `ui/`. App-specific components live under `$lib/components/app/{sidebar,main_screen,connections,connection-forms,user}`. Flowbite Svelte is also a dependency.
- Monaco is loaded at runtime via `@monaco-editor/loader` in `main_screen/sql_editor.svelte`, which registers the custom `aurora-sql` / `aurora-sql-light` themes and a SQL completion provider fed by the `suggestions` store. `vite.config.ts` carries a `monaco-esm-fix` plugin setting COOP/COEP headers in dev.
- Result grids use `@tanstack/table-core`.
- Dark mode: Tailwind 3 `darkMode: ["class"]` plus `mode-watcher` (`<ModeWatcher />` in `+layout.svelte`).

## Known gaps

- Connection passwords and SSH keys are stored **in plaintext** in local SQLite (`migrations/0001_init.sql`). The README's "encrypted SQLite database" is aspirational.
- `app/connections.go` uses `ssh.InsecureIgnoreHostKey()`.
- PostgreSQL is the only implemented engine.
