# B1 Acceptance Criteria — Tasks & Workflows (MVP)

Goal: Implement minimal Tasks + Workflow Board features to support the dev-company operating loop.

## In Scope
- Task CRUD (minimal): create, list, inspect.
- Task state machine and transitions.
- Workflow board view (kanban) by state.
- Audit events for task creation and transitions.

## Out of Scope
- Real worker execution.
- Artifact generation.
- Reviewer approval workflows beyond task state.

## Data Model (MVP)
A task must include:
- `id` (string)
- `title` (string)
- `desc` (string, optional)
- `state` (enum)
- `assignee` (string, optional)
- `reviewer_required` (bool)
- timestamps: `created_at`, `updated_at`

## State Machine
Allowed states (MVP):
- `draft`
- `planned`
- `assigned`
- `running`
- `reviewing`
- `done`
- `blocked`

Allowed transitions (MVP):
- draft → planned
- planned → assigned
- assigned → running
- running → reviewing
- reviewing → done
- any → blocked
- blocked → planned (recovery)

## Acceptance Criteria

### AC-B1-BE-001: API endpoints
- `GET /api/tasks` returns 200 and JSON array.
- `POST /api/tasks` creates a task and returns 201 + created task.
- `GET /api/tasks/:id` returns 200 for existing task, 404 otherwise.
- `POST /api/tasks/:id/transition` applies a valid transition and returns 200.

### AC-B1-BE-002: Local store
- Tasks are stored under the data dir (default `~/.ai-company-os/`).
- Rebooting the server preserves tasks (file-backed store).

### AC-B1-AUD-001: Audit events
- Creating a task emits `task.create` audit event.
- Transitioning a task emits `task.transition` audit event including:
  - task id
  - from/to state
  - actor (api/cli)

### AC-B1-FE-001: Tasks page
- Tasks page shows:
  - list of tasks (id/title/state/assignee)
  - create form (title + optional desc)
- Empty/Loading/Error states exist.

### AC-B1-FE-002: Workflow board
- Workflows page shows a kanban board grouped by state.
- Tasks render under their current state.

### AC-B1-QA-001: Smoke flow
- With backend running and seeded or created tasks:
  - create a task in UI
  - transition it through at least 2 states
  - verify audit log reflects the actions

## Verification
- Manual UI smoke for B1.
- API smoke via curl.
