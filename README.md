# AI Company OS

Local, installable Company OS for running your AI engineering org like a real operating company.

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
