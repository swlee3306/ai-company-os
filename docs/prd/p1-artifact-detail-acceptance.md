# P1-1 Acceptance Criteria — Artifact detail + evidence link resolution

Goal: Improve evidence usability by providing an artifact detail view and making evidence bundle entries clickable.

## In Scope
- Add artifact detail route in FE: `/artifacts/:id`.
- Add BE endpoint already exists: `GET /api/artifacts/:id` (use it).
- Update Project Detail evidence bundle rendering:
  - if an evidence entry matches an artifact id (e.g., `A-...`), render it as a link to artifact detail.
  - otherwise render as plain text.

## Out of Scope
- Binary uploads.
- Artifact editing/deletion.

## Acceptance Criteria

### AC-P11-FE-001: Artifact detail page
- Navigating to `/artifacts/:id`:
  - fetches `GET /api/artifacts/:id`
  - shows Loading/Error/Empty states
  - displays fields: type/title/project_id/task_id/created_at/uri/meta.

### AC-P11-FE-002: Artifacts list links
- Artifacts list page links each artifact id/title to `/artifacts/:id`.

### AC-P11-FE-003: Project evidence link resolution
- In Project Detail, Evidence bundle section:
  - entries matching an artifact id become clickable links to artifact detail.

### AC-P11-QA-001: Smoke
- Create an artifact via API.
- Ensure it appears in Artifacts list and opens detail.
- Ensure a project with evidence entry matching that artifact id links correctly.
