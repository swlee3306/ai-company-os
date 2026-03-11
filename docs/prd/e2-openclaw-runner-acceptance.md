# E2 Acceptance Criteria — OpenClaw Runner Integration (ACP harness)

## Goal
E1에서 만든 Runs/Run 버튼의 “placeholder 실행”을 제거하고,
실제로 **OpenClaw ACP harness**를 통해 역할별 에이전트를 실행한다.

최종 사용자 경험:
- CEO가 Task를 만들고 Run을 누르면
- PM/FE/BE/QA/Reviewer 에이전트가 순서대로 실행되고
- 산출물/로그/결과가 Runs + Artifacts + Audit에 남는다.

## Principles
- 실행 엔진은 `openclaw sessions_spawn(runtime=acp)` 기반.
- 외부 비밀키를 OS가 저장하지 않는다(사용자는 OpenClaw/CLI 로그인 또는 env로 처리).
- 실패 시 재시도/재개가 가능하도록 step 단위로 기록.

## In Scope

### 1) Runner settings
- `settings.json`에 runner backend 설정 추가:
  - `runner.backend = openclaw_acp`
  - role→agent 매핑(최소): `pm`, `fe`, `be`, `qa`, `reviewer`
- Settings UI에서:
  - backend 선택(openclaw)
  - 각 role별 agentId 입력(또는 드롭다운; MVP는 입력)

### 2) Run execution
- `POST /api/tasks/:id/run`이 실제로:
  - run 디렉토리 생성
  - step 계획 생성(pm_only/full)
  - 각 step에서 OpenClaw ACP 세션 spawn
  - step stdout/stderr/요약을 run 디렉토리에 저장

OpenClaw 호출 규약(MVP):
- `sessions_spawn` with:
  - `runtime: "acp"`
  - `thread: true`, `mode: "session"`
  - `label: "aicos:<run_id>:<role>"`
  - `task:` role prompt + 작업 지시( Task title/desc + repo path + constraints )

### 3) Step state + resume
- run 폴더에 step state 저장:
  - `plan.json` (steps list)
  - `state.json` (current step, statuses)
- API:
  - `POST /api/runs/:id/resume` (failed step 이후 이어서 진행)

### 4) Artifacts + Audit
- step마다 최소 artifact 1개(로그) 생성:
  - `type=run_log`, `uri=file://.../runs/<id>/<role>.stdout.log`
- Audit:
  - `run.start`, `run.step.start`, `run.step.done`, `run.fail`, `run.done`

### 5) UX
- TaskDetail에서:
  - Run 시작
  - 현재 step/상태 표시
  - 실패 시 Resume 버튼
- Runs 페이지에서:
  - run 목록 + status
  - run detail 링크(최소)

## Out of Scope
- 자동 PR 생성/merge (E3로 분리)
- 분산 실행(k3s/k3d worker)
- 비밀키 저장/암호화 금고

## Acceptance Criteria

### AC-E2-EXEC-001: pm_only executes via OpenClaw
- pm_only pipeline으로 Run 실행 시:
  - OpenClaw ACP 세션이 실제로 생성됨
  - stdout/stderr/요약이 run 폴더에 남음
  - Audit에 step start/done 기록

### AC-E2-EXEC-002: full pipeline step chain
- full pipeline은 최소 3 step(PM→BE→QA)까지 순차 실행이 가능해야 함.
  - (FE/Reviewer는 enabled 옵션이더라도 MVP에서 warn로 스킵 가능)

### AC-E2-RESUME-001: resume
- 실패한 run은 Resume로 재시도 가능.

### AC-E2-QA-001: single batch QA
- QA는 1회 배치로:
  - runner backend openclaw 설정
  - task 생성
  - pm_only run 실행
  - artifacts/audit 확인

## UX Review requirement
- PM/Designer/Reviewer가 Settings(Runner) + TaskDetail(Run) UX를 리뷰하고
  최소 3개의 개선점(문구/상태/다음 행동)을 이슈 또는 후속 PR로 남긴다.
