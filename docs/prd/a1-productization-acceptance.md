# A1 Acceptance Criteria — Productization + Korean guide

Goal: Make the repo easy to install/demo/release for real users, with Korean documentation.

## In Scope
- Provide a repeatable build/release script.
- Provide a 1-minute demo path.
- Add Korean user guide (KR).

## Out of Scope
- Signed binaries/notarization.
- Installer packaging (Homebrew formula, deb/rpm) — later.

## Acceptance Criteria

### AC-A1-REL-001: Build/release script
- Add a script under `scripts/` (e.g., `scripts/release.sh`) that:
  - builds Go `company` binary
  - builds web (`web/`) production bundle
  - outputs artifacts to a single folder (e.g., `release/`)
  - prints the output paths

### AC-A1-REL-002: Version info
- `company version` command exists and prints:
  - app version (git commit or semver placeholder)
  - build time

### AC-A1-DOC-001: English quickstart remains
- `README.md` includes:
  - build + run backend
  - run frontend
  - seed + approvals/tasks demo

### AC-A1-DOC-002: Korean guide
- Add Korean documentation:
  - `README.ko.md` (or `docs/guide/ko.md`) with:
    - 설치/실행/데모 절차
    - 흔한 오류 해결(예: Docker Desktop 미기동, k3d 미설치)
    - driver 설정(k3d/k3s) 설명

### AC-A1-QA-001: Smoke
- On a clean machine with Go/Node installed:
  - running the release script completes
  - demo steps in KR guide can be followed end-to-end
