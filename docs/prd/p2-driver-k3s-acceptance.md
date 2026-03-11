# P2 Acceptance Criteria — Driver support: k3d default + k3s (Linux only)

Policy:
- macOS: k3d only
- Linux: k3d default + k3s optional (systemd lifecycle)

## Goals
- Users can choose driver in Settings.
- CLI/API reflect the selected driver.
- On Linux, k3s lifecycle is supported via systemd when selected.

## In Scope
### P2-1: k3s checks (Linux only)
- `status/doctor` includes real k3s checks when driver.selected=k3s:
  - `k3s --version` (or equivalent) runs
  - systemd service state (k3s) when available
  - basic kubectl connectivity when available

### P2-2: k3s lifecycle (Linux only)
- `company up` (driver=k3s): starts k3s service if installed.
- `company down` (driver=k3s): stops k3s service.
- Emits audit events: `driver.k3s.up`, `driver.k3s.down`.

### P2-3: Settings integration
- `settings.driver.selected` controls:
  - `company up/down`
  - `company status/doctor`
  - `/api/status` and `/api/doctor`

## Out of Scope
- k3s installation automation (can be P2.5).
- k3s on macOS.

## Acceptance Criteria

### AC-P2-SET-001: driver selection routes behavior
- When settings driver is k3d:
  - up/down use k3d paths
  - driver checks run k3d checks
- When settings driver is k3s:
  - on Linux: status/doctor/up/down use k3s paths
  - on macOS: status/doctor warn "k3s linux-only" and up/down fail with actionable message

### AC-P2-K3S-CHK-001: doctor checks
- On Linux with k3s installed:
  - doctor includes `k3s.version` check ok
  - doctor includes `k3s.service` check ok
- On Linux without k3s:
  - doctor includes `k3s.version` warn with install guidance

### AC-P2-K3S-UP-001: lifecycle
- On Linux with k3s installed:
  - `company up` starts k3s via systemd
  - `company down` stops k3s via systemd
- Both are idempotent.

### AC-P2-QA-001: smoke matrix
- macOS:
  - k3d path unaffected
  - selecting k3s yields clear linux-only warning
- Linux:
  - k3d path unaffected
  - k3s checks/lifecycle work when installed

