# D1 Acceptance Criteria — UI-first registration (create agents/projects/approvals/artifacts)

## Goal
데모(seed) 없이도, 사용자가 **웹 UI만으로** 운영 데이터를 등록하고 시스템을 실제로 사용 가능하게 만든다.

핵심 원칙:
- ID는 **서버 자동 생성**.
- 각 리스트 페이지에 최소 Create 섹션을 추가(큰 라우팅/구조 변경 최소화).
- 모든 생성/변경은 audit log로 남는다.

## In Scope

### 1) Backend API: Create endpoints
다음 create endpoints 제공(모두 JSON body):

- Agents
  - `POST /api/agents`
  - required: `name`, `persona_role`
  - optional: `ops_specialty`, `scope[]`, `version`, `heartbeat_seconds`, `approval_required`, `risk_scope[]`

- Projects
  - `POST /api/projects`
  - required: `name`, `summary`
  - optional: `phase`, `owner_ceo`, `team_lead`, `due`, `status`, `agents[]`, `evidence_bundle[]`

- Approvals
  - `POST /api/approvals`
  - required: `type`, `requester`, `target`, `risk`
  - optional: `action`(default: "approve or reject"), `task_id`

- Artifacts
  - `POST /api/artifacts` (already exists, extend UI usage)
  - required: `title`, `uri`
  - optional: `type`, `project_id`, `task_id`, `meta`

Validation:
- 잘못된 요청은 400 + `{error: string}`.

Audit:
- create마다 audit action 추가:
  - `agent.create`, `project.create`, `approval.create`, `artifact.create`
  - fields: created id + 주요 연관(project_id/task_id/target 등)

### 2) Frontend UI: Create forms on list pages
다음 페이지에 Create 섹션 추가:
- Agents 페이지: Agent 등록 폼
- Projects 페이지: Project 등록 폼
- Approvals 페이지: Approval request 생성 폼
- Artifacts 페이지: Artifact 등록 폼

UX 규칙:
- 성공 시 리스트 즉시 refresh
- Error/Loading 상태는 기존 UI tokens 스타일 사용
- 입력 최소화(advanced fields는 optional)

### 3) Traceability
- Approval 생성 시 `task_id`가 있으면:
  - Approval evidence 패널에서 task link 노출
- Artifact 생성 시 project/task 연계가 있으면:
  - Project evidence bundle(또는 Artifact list/detail)에서 추적 가능

## Out of Scope
- Edit/Delete UI, bulk import
- Auth/permissions
- 복잡한 폼(모든 필드를 다 노출하는 설정 UI)

## Acceptance Criteria

### AC-D1-API-001: Create endpoints work
- 각 `POST` endpoint가 파일 스토어에 append 저장 + audit 기록.

### AC-D1-UI-001: UI-only onboarding
- seed 없이도 웹 UI에서 다음을 수행 가능:
  1) Agent 생성 1개
  2) Project 생성 1개
  3) Task 생성 1개(기존 기능)
  4) Approval 생성 1개(task_id optional)
  5) Artifact 생성 1개(project/task 연결)

### AC-D1-QA-001: Single batch QA
- QA는 1회 배치로:
  - 위 5단계를 실제로 실행
  - Audit Logs에서 create 이벤트 확인
  - Approvals evidence에서 linked task 확인(해당 시)

### AC-D1-KR-DOC-001: 한국어 가이드 업데이트
- `README.ko.md`에 "seed 없이 시작" 절차와 UI 등록 방법 추가.
