# C1.5 Acceptance Criteria — k3d bootstrap + up/down (Linux + macOS)

Goal: Provide a practical local runtime bootstrap path (install k3d if missing) and minimal lifecycle (`up/down`) using k3d.

## In Scope
- Detect missing k3d and provide an actionable install path.
- Provide `company up` and `company down` using k3d.
- Extend `status`/`doctor` to reflect real k3d cluster state.
- Linux + macOS support.

## Out of Scope
- Full k3s install automation.
- Deploying apps into the cluster.
- Sophisticated configuration matrix (tunable later).

## Acceptance Criteria

### AC-C15-DOCTOR-001: actionable k3d missing guidance
- If k3d is not installed:
  - `company doctor` reports k3d check as warn/error
  - includes actionable remediation instructions (Linux + macOS)
  - does not crash

### AC-C15-INSTALL-001: optional bootstrap helper
- Provide one of:
  - `company install k3d` (preferred)
  - or `company doctor --fix` (acceptable)
- It must:
  - on Linux: install k3d via a documented method (e.g., official install script)
  - on macOS: install via Homebrew when available, else provide instructions
- It must be safe:
  - prints what it will do
  - fails fast with a clear error if prerequisites (curl/brew) missing

### AC-C15-UP-001: minimal `up`
- `company up` creates a k3d cluster if missing.
- `company up` is idempotent (second run does not error; reports already running).
- On success, emits audit event `driver.k3d.up`.

### AC-C15-DOWN-001: minimal `down`
- `company down` deletes the k3d cluster if it exists.
- `company down` is idempotent (second run does not error).
- On success, emits audit event `driver.k3d.down`.

### AC-C15-STATUS-001: status reflects cluster state
- `company status` reports:
  - k3d present
  - cluster count
  - whether the default cluster exists

### AC-C15-QA-001: smoke (Linux + macOS)
- On Linux and macOS test machines:
  - if k3d missing: doctor shows actionable steps
  - install helper installs k3d (if allowed)
  - up → status shows cluster
  - down → status shows cluster removed

## Notes / Defaults
- Default cluster name: `company-os` (unless specified)
- Default k3d create args are minimal (tunable later)
