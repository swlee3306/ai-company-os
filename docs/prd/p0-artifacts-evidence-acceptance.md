# P0-1 Acceptance Criteria — Artifacts + Evidence Bundle (MVP)

Goal: Make Evidence Bundle and Artifacts first-class so approvals/reviews can rely on concrete outputs.

## In Scope
- Persist artifacts in local store.
- Provide API endpoints to list and inspect artifacts.
- Implement Artifacts page in the dashboard.
- Connect Project Detail evidence bundle to artifacts.

## Out of Scope
- Real build pipelines.
- Binary artifact uploads; for MVP, artifacts are metadata (links/paths) only.

## Data Model (MVP)
Artifact fields:
- `id` (string)
- `type` (enum/string: ADR, schema-diff, benchmark, runbook, screenshot, export, log, other)
- `title` (string)
- `project_id` (string, optional)
- `task_id` (string, optional)
- `uri` (string) — can be file path, URL, or repo-relative reference
- `created_at` (RFC3339)
- `meta` (object, optional)

Evidence bundle:
- A project’s `evidence_bundle` can reference artifact ids OR freeform strings.
- UI should render artifact ids as clickable entries when resolvable.

## Acceptance Criteria

### AC-P0A-BE-001: Store
- Artifacts are stored under the data dir (default `~/.ai-company-os/`) as `artifacts.json`.
- Store is file-backed; restarting server preserves artifacts.

### AC-P0A-BE-002: API
- `GET /api/artifacts` returns 200 and JSON array.
- `POST /api/artifacts` creates an artifact and returns 201.
- `GET /api/artifacts/:id` returns 200 for existing artifact, 404 otherwise.

### AC-P0A-AUD-001: Audit
- Creating an artifact emits `artifact.create` audit event with `artifact_id` and context (project/task when present).

### AC-P0A-FE-001: Artifacts page
- Artifacts page lists artifacts with: type, title, project/task, created_at, uri.
- Empty/Loading/Error states exist.

### AC-P0A-FE-002: Project Detail evidence links
- Project Detail renders an **Evidence bundle** section.
- If an evidence entry matches an artifact id, it links to the artifact detail view.

### AC-P0A-QA-001: Smoke
- Seed or create at least 1 artifact via API.
- Confirm it appears on the Artifacts page.
- Confirm Project Detail evidence bundle links resolve for known artifact ids.

## Verification
- API smoke via curl.
- Manual UI check.
