# E2.4 Acceptance Criteria — Workspace repo path + codex exec wiring

## Goal
사용자가 운영하는 프로젝트 레포(working repo)를 Settings에 등록하고,
local_cli runner가 해당 레포에서 **`codex exec`**로 non-interactive 실행을 수행하도록 한다.

## In Scope

### 1) Settings: workspace repo path
- `settings.json`에 workspace 추가:
  - `workspace.repo_path` (string)
- Settings UI에 입력 필드 추가:
  - Workspace repo path

### 2) Preflight validation
- API:
  - `POST /api/workspace/validate` with `{ repo_path }`
  - checks:
    - directory exists
    - `.git/` exists (git repo)
- UI에서 Validate 버튼 제공 + 결과 표시(ok/warn/fail)

### 3) Runner execution uses repo_path
- local_cli backend + codex command일 때:
  - 실행은 `codex exec` 서브커맨드를 사용
  - `-C <workspace.repo_path>` 적용
  - sandbox/approval 기본값:
    - `-s workspace-write`
    - `-a untrusted` (또는 on-request; MVP는 untrusted)
- stdout/stderr/RESULT.json 저장은 유지

### 4) Prompt
- pm_only pipeline prompt에 다음 포함:
  - Task title/desc
  - repo_path
  - constraints: branch/PR only, QA once, audit-first

## Out of Scope
- multi-repo 지원
- PR 자동 생성/merge

## Acceptance Criteria

### AC-E24-SET-001
- Settings에서 repo_path 저장/로드가 가능.

### AC-E24-VAL-001
- validate API가 정상/오류를 구분해 반환.

### AC-E24-EXEC-001
- codex 설치된 환경에서 pm_only run이 `codex exec -C <repo_path> ...`로 실행되고,
  stdout.log에 출력이 남는다.

### AC-E24-QA-001 (single batch)
- QA는 1회 배치로:
  - repo_path 설정 + validate
  - task 생성
  - run 실행
  - runs stdout/stderr 확인
