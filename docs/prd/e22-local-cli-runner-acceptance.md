# E2.2 Acceptance Criteria — local_cli runner + init + secrets (no manual env)

## Goal
OpenClaw 없이도 AI Company OS에서 Task → Run 버튼으로 로컬 CLI 에이전트를 실행할 수 있게 한다.
또한 사용자가 환경변수를 직접 만지지 않아도 되도록, 설치/초기화 단계에서 비밀값을 로컬에 저장하고 실행 시 주입한다.

## Principles
- default runner: codex CLI
- secrets는 audit에 기록하지 않는다.
- secrets 저장 파일은 로컬 권한(0600)으로 보호한다.

## In Scope

### 1) Secrets store (local)
- 데이터 디렉토리 하위에 secrets 파일 추가:
  - `${AI_COMPANY_OS_HOME}/secrets.json`
- API key 등 민감정보는 settings.json이 아니라 secrets.json에 저장.
- 파일 권한: 0600

### 2) Init wizard (CLI)
- `company init` 커맨드 추가:
  - runner backend 선택: `local_cli` / `openclaw_acp`
  - local_cli 선택 시:
    - runner.type=codex_cli (default)
    - runner.command(default: `codex`)
    - API key 입력(선택) → secrets.json에 저장
  - openclaw_acp 선택 시:
    - gateway url/token 입력(선택) → secrets.json에 저장
    - role agentIds 입력(최소 pm)

### 3) Runner backend: local_cli
- settings의 `runner.backend=local_cli` 지원
- `POST /api/tasks/:id/run`에서 local_cli일 때:
  - `runner.command`를 subprocess로 실행
  - 표준 출력/에러를 run 디렉토리에 저장
  - 최소 pm_only pipeline 동작

### 4) Preflight
- Settings 화면에 runner 상태(preflight) 표시:
  - local_cli: command 존재 여부
  - openclaw_acp: gateway reachability + sessions_spawn 가능 여부

### 5) Audit + Artifacts
- run.start/step/done 이벤트 유지
- run_log artifact 자동 생성 유지

## Out of Scope
- Keychain(macOS) 통합
- 자동 PR 생성/merge
- distributed workers

## Acceptance Criteria

### AC-E22-INIT-001
- `company init`으로 settings + secrets를 생성 가능.

### AC-E22-SEC-001
- secrets.json은 0600 권한으로 생성된다.
- audit log에는 secrets 값이 포함되지 않는다.

### AC-E22-RUN-001
- local_cli(pm_only) run이 실제로 subprocess를 실행하고 stdout/stderr가 저장된다.

### AC-E22-QA-001 (single batch)
- QA는 1회 배치로:
  - init 실행
  - task 생성
  - run 실행
  - runs/ artifacts/ audit 확인
