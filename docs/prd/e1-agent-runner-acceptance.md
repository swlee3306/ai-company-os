# E1 Acceptance Criteria — Agent Runner MVP (OpenClaw-style registration)

## Goal
CEO가 자연어로 만든 Task를 기반으로, **OS에서 버튼을 눌러 역할별 에이전트를 실행**하고
결과(로그/산출물/승인)를 audit-first로 기록하는 "회사형 개발 파이프라인"의 MVP를 제공한다.

핵심 원칙:
- Runner(실행 엔진)는 Settings에서 **등록**한다.
- API key/로그인은 OS가 보관하지 않는다. (OpenClaw 스타일: 사용자가 CLI 로그인 또는 env 설정)
- 실행 결과는 Run 디렉토리 + Artifacts + Audit로 남긴다.

## In Scope

### 1) Runner registry (Settings)
- Settings UI에 Runner 등록 섹션 추가:
  - runner type: `codex_cli` (default template), `claude_code`, `gemini_cli`, `custom`
  - executable / command template
  - working dir (default repo root)
  - optional env var hints (문서 안내용)
- 저장: `settings.json`에 `runner` 블록 추가.

### 2) Run entity + storage
- 파일 스토어에 runs 저장:
  - `~/.ai-company-os/runs/<run_id>/`에:
    - `request.json` (task id, selected runner, pipeline)
    - `prompt.txt`
    - `stdout.log`, `stderr.log`
    - `RESULT.json` (status, summary, artifacts, evidence)
- API:
  - `POST /api/tasks/:id/run` : run 생성 + 실행 시작
  - `GET /api/runs` : run 목록
  - `GET /api/runs/:id` : run detail (메타 + tail log 링크/요약)

### 3) Pipelines (minimal)
- 최소 2가지 실행 모드 제공:
  - `pm_only`
  - `full` (pm→fe→be→qa→reviewer) 는 step만 구성하고, 실패 시 중단
- 각 step은 role tag를 가진 prompt로 실행된다.

### 4) Audit integration
- audit events:
  - `run.start`, `run.step.start`, `run.step.done`, `run.fail`, `run.done`
  - fields: run_id, task_id, runner_type, step(role)

### 5) Artifact integration
- run 완료 시 최소 1개 artifact 자동 생성:
  - type: `run_log`
  - uri: `file://.../runs/<run_id>/stdout.log` (또는 상대 경로)
  - meta에 run_id, task_id 포함

### 6) UI integration (minimal)
- Tasks 페이지/TaskDetail에:
  - Run 버튼
  - 최근 run 상태 표시(최근 1~3개)
- Runs 페이지(간단): run 목록 + status/summary 표시

## Out of Scope
- 분산 실행(k3s/k3d worker) / 멀티 노드
- 비밀키 저장/암호화 금고
- 실제 git PR 생성/merge 자동화(후속 E2로 분리 가능)

## Acceptance Criteria

### AC-E1-RUN-001: Task run works
- seed 없이:
  - Task 1개 생성
  - Run 버튼 클릭
  - run이 생성되고 상태가 UI/API에서 조회 가능

### AC-E1-LOG-001: Logs persisted
- run 디렉토리에 stdout/stderr/RESULT.json이 남는다.

### AC-E1-AUD-001: Audit trail
- run의 시작/단계/완료가 audit log에 남는다.

### AC-E1-ART-001: Artifact created
- run 완료 후 artifacts 리스트에 run_log artifact 1개 이상 존재.

### AC-E1-QA-001: Single batch QA
- QA는 1회 배치로:
  - runner 등록(템플릿)
  - task 생성
  - run 실행
  - audit/artifact 확인

## UX review requirement (PM/Designer/Reviewer)
- PM/Designer/Reviewer가 Runner 등록 화면과 Tasks/Run UX를 리뷰하고,
  "혼동되는 필드/누락된 상태/다음 행동"을 개선 제안으로 남긴다(후속 PR로 분리).
