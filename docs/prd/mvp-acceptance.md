# MVP Acceptance Criteria (PM)

Scope: local installable CLI + local dashboard (web) with evidence/audit-first operating model.

## In Scope (MVP)
- Local Go backend (Gin) API server + CLI (`company`).
- Local file store under `~/.ai-company-os/` (override via `AI_COMPANY_OS_HOME`).
- Dashboard skeleton (React+Vite) with left navigation and page scaffolds.
- Demo data seeding for end-to-end smoke via `company seed`.
- Approval loop with decision gate + evidence lookup.

## Out of Scope (MVP)
- Real k3d/k3s driver control and real cluster lifecycle.
- Real task execution, worker orchestration, or artifact generation.
- Authentication/SSO.
- Multi-user/RBAC.

## Definitions
- Evidence bundle (MVP): minimal linked context shown in UI for approvals (linked agent/project + recent audit).
- Decision log (MVP): reject requires a reason and is persisted.

## Acceptance Criteria

### AC-CLI-001: CLI runs and writes audit
- Given a fresh machine state
- When running `company status` and `company doctor`
- Then both commands exit 0
- And `~/.ai-company-os/audit.log` contains appended JSONL entries with `ts`, `actor`, `action`.

### AC-CLI-002: Demo seeding
- When running `company seed`
- Then `agents.json`, `projects.json`, `approvals.json` exist under the data dir
- And subsequent API calls return non-empty lists.

### AC-API-001: Core read endpoints
- When the API server is running via `company serve --listen 127.0.0.1:8787`
- Then the following endpoints return HTTP 200:
  - `GET /api/status`
  - `GET /api/doctor`
  - `GET /api/audit`
  - `GET /api/agents`
  - `GET /api/projects`
  - `GET /api/approvals`

### AC-API-002: Detail endpoints
- `GET /api/agents/:id` returns 200 for a seeded agent id, else 404.
- `GET /api/projects/:id` returns 200 for a seeded project id, else 404.
- `GET /api/approvals/:id` returns 200 for a seeded approval id, else 404.

### AC-APPROVAL-001: Reject requires reason
- Given an approval item id
- When calling `POST /api/approvals/:id/decision` with `{decision:"reject", reason:""}`
- Then it returns HTTP 400.

### AC-APPROVAL-002: Approve updates approval status
- When calling `POST /api/approvals/:id/decision` with `{decision:"approve"}`
- Then it returns HTTP 200
- And the approval item `status` becomes `approve` in subsequent list/detail reads.

### AC-APPROVAL-003: Evidence endpoint
- `GET /api/approvals/:id/evidence` returns 200
- And includes:
  - `approval`
  - `audit_recent` (up to 10)
  - plus `agent` and/or `project` if linkable by target.

### AC-APPROVAL-004: Side effects (MVP rule engine)
- For type `tool permission`:
  - approve sets linked agent status to `active` and `approval_required=false`
  - reject sets linked agent status to `blocked` and `approval_required=true`
- For type `production deploy`:
  - approve sets linked project status to `running`
  - reject sets linked project status to `blocked`

### AC-FE-001: Dashboard boots and navigates
- When running `web` via `npm run dev`
- Then the app loads without runtime errors
- And the left navigation routes to:
  - Dashboard, Projects, Agents, Approvals, Audit Logs, Settings
- And pages show Loading/Error/Empty states reasonably.

### AC-FE-002: Approvals UI decision loop
- In the Approvals page:
  - selecting an item loads evidence panel
  - submitting `reject` without reason shows an error
  - submitting `approve` changes the item status and refreshes the list
  - agent/project side effects are observable via Agents/Projects pages after refresh.

## Verification
- Backend/API smoke can be validated with curl scripts.
- Frontend smoke is manual for MVP (automation in later phase).
