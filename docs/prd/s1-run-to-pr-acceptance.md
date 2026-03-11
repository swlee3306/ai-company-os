# S1 Acceptance Criteria — Run(full) → Branch/Commit/PR → QA artifact

## Goal
Task 하나로 실제 개발 파이프라인을 끝까지 실행한다:
- PM 계획 생성
- BE 구현(브랜치/커밋)
- GitHub PR 생성
- QA 1회 배치 실행
- 결과를 Runs + Artifacts + Audit로 남김

## Assumptions
- local_cli backend is used with codex CLI.
- GitHub CLI(`gh`)가 로컬에 설치되어 있고 인증되어 있다.
- Repo path는 Settings `workspace.repo_path`로 지정한다.

## In Scope

### 1) full pipeline steps (minimal)
Run(full) 실행 시 step 순서:
1) pm: AC + plan 생성 (stdout.log)
2) be: 작업 브랜치 생성 + 변경 적용 + 커밋
3) pr: `gh pr create`로 PR 생성, URL을 artifact로 저장
4) qa: `scripts/verify.sh` 1회 실행, 결과 로그 artifact 저장

### 2) Safe command allowlist
- git/gh/verify.sh 등 필요한 커맨드만 허용(초기 allowlist)
- 위험 커맨드는 approval gate로 분리(후속)

### 3) Evidence
- Artifacts 최소:
  - run_log (step별 stdout)
  - pr_link (PR URL)
  - qa_log (verify output)
- Audit:
  - run.step.start/done for each step

## Out of Scope
- auto-merge
- multi-reviewer approval
- advanced diff apply tooling

## Acceptance Criteria

### AC-S1-PR-001
- Run(full) 이후 PR이 실제로 생성되고 URL이 artifact로 남는다.

### AC-S1-QA-001
- QA step은 정확히 1회 실행되고 결과가 artifact로 남는다.

### AC-S1-TRACE-001
- Runs 페이지에서 run_id를 따라가면 step 로그와 PR/QA artifact를 확인할 수 있다.

### AC-S1-KR-DOC-001
- README.ko.md에 "Task → Run(full) → PR" 데모 절차 추가.
