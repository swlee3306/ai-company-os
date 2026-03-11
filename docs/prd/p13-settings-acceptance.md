# P1-3 Acceptance Criteria — Settings (driver + approval policy)

Goal: Provide minimal settings to control driver selection and display approval risk policy.

## In Scope
- Persist settings in local store.
- Expose settings via API.
- Implement Settings page UI to view/update settings.

## Out of Scope
- Auth/RBAC.
- Complex policy DSL.

## Settings (MVP)
- `driver.selected`: `k3d` | `k3s`
- `approval.policy_text`: string (display-only)

## Acceptance Criteria

### AC-P13-BE-001: Store
- Settings stored under data dir as `settings.json`.
- Defaults exist if file missing.

### AC-P13-BE-002: API
- `GET /api/settings` returns 200 with settings.
- `POST /api/settings` updates allowed fields and returns 200.

### AC-P13-FE-001: Settings UI
- Settings page shows:
  - Driver select dropdown (k3d/k3s)
  - Approval policy text block
- Saving updates settings via API and persists.

### AC-P13-INTEG-001: Driver selection affects status/doctor
- `company status` and `company doctor` reflect selected driver in output (`driver.selected`).
- (MVP) if k3s selected but not implemented, checks warn gracefully.

### AC-P13-QA-001: Smoke
- Change driver to k3s, refresh status/doctor, confirm selected shows.
- Change back to k3d.
