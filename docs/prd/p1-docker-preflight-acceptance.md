# P1 Acceptance Criteria — Docker preflight UX for k3d up/down

Goal: Make `company up` failures actionable when Docker daemon is not running (common on macOS).

## In Scope
- Preflight check for Docker daemon before invoking k3d operations.
- Clear, short error messaging with remediation hints.
- Preserve current behavior when Docker is healthy.

## Out of Scope
- Full Docker Desktop install automation.
- Deep diagnostics beyond `docker info` reachability.

## Acceptance Criteria

### AC-P1-UP-001: `company up` docker preflight
- When Docker daemon is not reachable:
  - `company up` exits non-zero
  - prints a concise message (one screen) including:
    - that Docker daemon is not running/reachable
    - on macOS: hint to start Docker Desktop
    - on Linux: hint to start docker service
- It must NOT run `k3d cluster create` when docker preflight fails.

### AC-P1-DOWN-001: `company down` docker preflight (soft)
- If Docker daemon is not reachable:
  - `company down` returns a warning
  - still exits 0 (idempotent cleanup semantics), unless a strict flag is added later.

### AC-P1-DOCTOR-001: doctor surfaces docker reachability
- `company doctor` includes a check:
  - `docker.daemon` (ok/warn)
  - detail includes the socket path error when failing.

### AC-P1-QA-001: regression
- On a machine with Docker running + k3d installed:
  - `company up` and `company down` still work as before.
- On a machine with Docker stopped:
  - `company up` fails fast with actionable hint.
