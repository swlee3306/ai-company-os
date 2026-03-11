# B2 Acceptance Criteria — E2E 운영 시나리오: Production Deploy 승인 플로우

## Goal
운영자가 “프로덕션 배포”를 요청했을 때, **Task → Approval → Evidence/Artifact → Side-effect**가 한 번의 시나리오로 끝까지 이어지도록 한다.

이 단계는 UI/UX polish가 아니라 **운영 흐름의 추적성(traceability) 완성**이 목적이다.

## Scope (In)

### 1) Seed/데모 데이터
- `company seed` 실행 시 아래 데모 시나리오가 재현 가능해야 한다.
- 최소 1개의 project가 `production deploy` 성격의 approval을 필요로 하도록 구성한다.

### 2) Task ↔ Approval ↔ Artifact ↔ Audit 연결
- Task가 생성되면 다음이 가능해야 한다:
  - Task detail에서 연결된 Approval(있다면)로 이동
  - Approval evidence 패널에서 해당 Task로 이동
  - Approval decision 이후 Task/Project 상태 변화가 audit에 남음
- Artifact를 1개 이상 생성해서(예: 배포 플랜 문서/런북 링크),
  - Approval evidence에서 **Artifact 링크가 클릭 가능**해야 한다.

### 3) 승인(approve/reject) + 사유
- reject 시 reason required 정책 유지.
- approve 시 side-effect가 발생:
  - project status 전이(예: blocked → running) 또는 deploy 관련 상태 변화
  - audit에 `cause=approval` 등 연쇄 원인이 남아야 한다.

### 4) UI 플로우(최소)
- 사용자는 웹에서 다음 순서로 확인 가능해야 한다:
  1) Task 목록에서 배포 Task 선택
  2) Task detail에서 연결된 Approval 확인/이동
  3) Approval Center에서 evidence 확인
  4) approve/reject 수행
  5) Project detail에서 상태 변화 및 evidence 링크 확인
  6) Audit Logs에서 일련의 이벤트 추적 가능

## Out of Scope
- 실제 배포 실행(쿠버네티스 apply, helm 등)
- 멀티 승인자/다단계 승인
- 외부 시스템 연동(GitHub Actions 등)

## Acceptance Criteria

### AC-B2-E2E-001: Traceability
- 아래 엔티티들이 서로 링크/ID로 추적 가능:
  - Task ↔ Approval
  - Approval ↔ Project
  - Approval evidence ↔ Artifact
  - 모든 side-effect는 audit로 추적 가능

### AC-B2-E2E-002: Evidence completeness
- Approval evidence 패널에 최소 포함:
  - approval 자체 정보
  - linked project
  - linked task(있다면)
  - 최근 audit
  - evidence bundle(artifact id 포함) 클릭 링크

### AC-B2-E2E-003: Side-effect + audit
- approve 시:
  - 상태 전이가 발생(프로젝트 또는 태스크)
  - audit에 `cause=approval`와 함께 기록

### AC-B2-KR-DOC-001: 한국어 데모 가이드
- `README.ko.md` 또는 `docs/guide/ko.md`에 B2 데모 절차(스크린 기준)가 추가되어야 한다.

### AC-B2-QA-001: QA 배치 1회
- QA는 1회 배치로 다음을 확인:
  - seed로 시나리오 데이터 준비
  - UI에서 링크 이동 최소 3번 이상( Task→Approval→Artifact ) 성공
  - approve/reject 동작 확인(최소 approve)
  - audit에서 연쇄 이벤트 확인
