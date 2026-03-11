# C1 Acceptance Criteria — UI polish (spacing/typography consistency)

## Goal
웹 UI를 MVP 수준에서 한 단계 정리: 화면별 padding/gap/card 스타일을 통일하고, 정보 밀도/가독성을 개선한다.

## Principles
- 불필요한 구조 리팩터링 금지(큰 diff 금지).
- 디자인 규칙(기존 합의)을 코드 스타일에 반영:
  - page padding=32, section gap=24
  - card gap=16
  - (modal height 규칙 등 기존 디자인 문서 내용은 유지)

## In Scope

### 1) Layout tokens (code-level)
- `web/src/ui/tokens.ts` (또는 유사 파일)로 spacing/typography 상수 정의:
  - `space.page=32`, `space.section=24`, `space.card=16`, `radius.card=12`, `border=1px #e5e7eb`
- 주요 페이지에서 하드코딩 숫자를 점진적으로 tokens로 치환(최소 3~5곳).

### 2) Page consistency
- 최소 아래 페이지들의 레이아웃 규칙을 맞춘다:
  - Projects / Project Detail
  - Approvals
  - Tasks / Task Detail
  - Artifacts / Artifact Detail
  - Audit Logs

### 3) Empty/Loading/Error states
- Empty/Loading/Error 컴포넌트 스타일을 통일(폰트/색상/여백).

## Out of Scope
- 디자인(Pencil .pen) 대대적 수정 및 PNG 재-export 강제.
- 새로운 컴포넌트 라이브러리 도입.

## Acceptance Criteria

### AC-C1-UI-001: Tokens introduced
- spacing/typography tokens 파일이 존재하고, 최소 3페이지에서 사용된다.

### AC-C1-UI-002: Visual consistency
- 페이지 상단 타이틀/설명 텍스트 간격이 일관.
- 카드 border/radius/padding이 일관.

### AC-C1-QA-001: Build
- `web npm run build` PASS

### AC-C1-KR-DOC-001: 한국어 노트
- `README.ko.md`에 UI 규칙(여백/카드) 한 단락 추가(개발자 참고).
