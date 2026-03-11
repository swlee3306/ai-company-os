# AI Company OS (한국어 가이드)

로컬(설치형)로 실행되는 **AI 개발회사 운영 시스템**입니다.

- 백엔드: Go(Gin)
- 프론트엔드: React + Vite
- 로컬 저장소: `~/.ai-company-os/` (환경변수 `AI_COMPANY_OS_HOME`로 변경 가능)

## 1분 데모

### 1) 백엔드 실행

```bash
cd ai-company-os

go build -o company ./cmd/company
./company seed
./company serve --listen 127.0.0.1:8787
```

### 2) 프론트 실행

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

브라우저에서 Vite 주소로 접속 후 메뉴를 이동해보세요:
- Dashboard / Projects / Agents / Tasks / Workflows / Approvals / Artifacts / Audit Logs / Settings

## 흔한 오류 해결

### Docker Desktop이 꺼져있을 때(k3d)
- 증상: `company up` 실행 시 docker daemon 연결 오류
- 해결: Docker Desktop을 실행하고 완전히 기동된 후 다시 실행

### k3d 미설치
- `company install k3d --dry-run`으로 설치 플랜 확인
- macOS: `brew install k3d`

## 드라이버(driver) 선택
- 기본: k3d
- Linux에서는 k3s 옵션도 선택 가능(Settings에서 변경)
  - 운영 노드에서는 up/down 테스트 주의(별도 테스트 노드 권장)

## B2 데모(운영 시나리오): Production Deploy 승인 플로우

목표: **Task → Approval → Evidence/Artifact → Side-effect → Audit**가 한 흐름으로 이어지는지 확인합니다.

1) seed 실행
```bash
./company seed
```

2) 웹에서 확인
- Tasks로 이동 → `T-184` 클릭
- Linked approval 안내를 확인하고 Approval Center로 이동
- Approval Center에서 `appr-1 (production deploy)` 선택
  - Evidence bundle에서 `A-20260311-DEPLOYPLAN` 링크 클릭 → Artifact detail로 이동
- approve 실행
- Projects에서 `billing-core` 상태/증거(evidence bundle) 확인
- Audit Logs에서 `approvals.decision` + `cause=approval` 이벤트 추적

## UI 규칙(개발 참고)
- page padding: 32
- section gap: 24
- card padding: 16

## 릴리즈 빌드

```bash
./scripts/release.sh
ls -la release/
./release/company version
```
