# AI Company OS

## 한 줄 소개
로컬 환경에서 AI 작업 흐름을 운영하기 위한 실험 저장소입니다.

## 저장소 성격
- 분류: 백엔드 / 플랫폼 실험
- 목적: 서버, CLI, 웹 UI를 함께 다루는 로컬 운영 도구
- 핵심 기술: Go, React, Vite, CLI

Experimental local workspace for coordinating AI engineering tasks through a Go CLI/API and React UI. This is not a production multi-tenant service or evidence of autonomous company operation.

Verification (2026-09-11): `go build -mod=readonly ./...` passed. Real provider calls, long-running workflows and production reliability were not validated by that build.

## Quick start (MVP)

### 1) Backend (Go)

#### Build

```bash
go build -o company ./cmd/company
```

### 2) Run API server (for dashboard)

```bash
./company serve --listen 127.0.0.1:8787
```

### 3) CLI basics

```bash
./company status
./company doctor
```

## Local storage

By default, runtime state is stored under:

- `~/.ai-company-os/state.json`
- `~/.ai-company-os/doctor.json`
- `~/.ai-company-os/audit.log`

Override with:

- `AI_COMPANY_OS_HOME=/path/to/dir`

## Design

Design source:
- `design/company-ui.pen`

Exports:
- `design/exports/`

Check exports are up to date:

```bash
./scripts/check-design-exports.sh
```

### 2) Frontend (React + Vite)

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

Then open the Vite URL and ensure the backend API is running at `VITE_API_BASE`.

## Korean guide

- See [README.ko.md](./README.ko.md)
