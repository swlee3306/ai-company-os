# E2.3 Acceptance Criteria — local_cli 실제 실행 (non-interactive) + 로그/결과 수집

## Goal
E2.2의 local_cli runner가 단순 `--version` 검증 단계를 넘어,
Task를 입력으로 받아 **비인터랙티브 방식으로 실제 실행**하고(stdout/stderr/RESULT.json),
Artifacts/Audit에 증거를 남긴다.

## Principles
- MVP: non-interactive 실행(한 번 실행하고 종료)
- 안전: 실행 커맨드 allowlist/템플릿 기반(임의 shell injection 방지)
- 결과는 run 폴더에 남기고, artifact로 자동 연결

## In Scope

### 1) Runner command template
- settings에 `runner.command`는 "프로그램명" + 고정 옵션 형태로 제한.
  - 예: `codex` 또는 `/usr/local/bin/codex`
- OS는 shell로 실행하지 않고 `exec.Command`로 argv를 구성한다.

### 2) Prompt bridging
- run 생성 시 `prompt.txt`를 작성하고,
  local_cli 실행은 아래 중 하나를 사용(가능한 방식 우선):
  - `codex --prompt-file <path>` (지원 시)
  - 또는 stdin으로 prompt 전달

### 3) Output capture
- stdout/stderr를 `runs/<run_id>/stdout.log`, `stderr.log`에 저장
- 종료 코드/에러는 `RESULT.json`에 기록

### 4) Minimal “pm_only” behavior
- pm_only pipeline은 실제로 로컬 CLI 1회 실행되어야 함.
- 실행 결과(요약/AC/next steps)가 stdout에 포함되도록 prompt를 구성.

### 5) Artifacts/Audit
- run_log artifact 자동 생성 유지
- audit:
  - `run.start`
  - `run.step.start` (role=pm)
  - `run.step.done` or `run.fail`
  - `run.done`

### 6) UX
- TaskDetail에서 run 실행 후:
  - run status(done/failed)
  - 최근 run 링크(/runs)

## Out of Scope
- full pipeline multi-step chaining(별도)
- PR 자동 생성/merge

## Acceptance Criteria

### AC-E23-EXEC-001
- local_cli(pm_only) run이 실제로 CLI를 실행한다.
  - stdout.log에 모델 출력이 남는다.

### AC-E23-SEC-001
- runner.command는 shell injection 없이 exec.Command로만 실행된다.

### AC-E23-QA-001 (single batch)
- QA는 1회 배치로:
  - init으로 local_cli 설정
  - task 생성
  - run 실행
  - stdout/stderr/RESULT + artifacts + audit 확인
