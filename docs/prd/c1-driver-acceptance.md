# C1 Acceptance Criteria — k3d/k3s driver integration (Phase 1)

Goal: Wire the CLI and `doctor/status` to real local environment checks and (optionally) real runtime start/stop using k3d first, with a path to k3s.

## In Scope (C1)
- Detect Docker availability and basic health.
- Detect k3d availability and list clusters.
- Implement real `company status` details for driver layer.
- Implement real `company doctor` checks for driver layer.

## Optional (C1.5)
- Implement `company up/down` to create/delete a k3d cluster (minimal).

## Out of Scope (C1)
- Full Kubernetes app deployment.
- Task execution on cluster.
- k3s installation automation.

## Driver selection
- Default driver: `k3d`.
- `--driver k3s` is allowed but may return "not implemented" in C1 if k3s tooling is missing.

## Acceptance Criteria

### AC-C1-CLI-001: status reports driver signals
- When running `company status`
- Then it reports:
  - selected driver (k3d/k3s)
  - docker: available/unavailable
  - k3d: available/unavailable
  - k3d clusters count (if available)

### AC-C1-CLI-002: doctor performs real checks
- When running `company doctor`
- Then it performs and reports checks with severity:
  - docker socket reachable
  - `docker info` succeeds
  - `k3d version` succeeds
  - `k3d cluster list` succeeds
- Output is persisted to `doctor.json`.

### AC-C1-API-001: /api/status includes driver checks
- `GET /api/status` includes a `driver` block with the same signals as CLI status.

### AC-C1-API-002: /api/doctor reflects real driver checks
- `GET /api/doctor` returns the latest driver check results.

### AC-C1-QA-001: smoke
- On a machine with Docker + k3d installed:
  - status and doctor return OK and show at least one driver signal.
- On a machine missing k3d:
  - doctor returns a clear actionable error/warning (does not crash).

## Verification
- Manual: `company status`, `company doctor`.
- Automated: minimal bash smoke in QA step.
