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

## Korean guide

- See [README.ko.md](./README.ko.md)

## Task T-20260312-092203

### Acceptance criteria

- PM plan is documented in `README.md` for Task `T-20260312-092203
- Execution is explicitly limited to a branch and PR workflow; no direct mainline delivery
- The plan explicitly states audit-first sequencing
- The plan explicitly states a single QA batch occurs after implementation is complete

### Execution plan

1. Create a task branch for `T-20260312-092203` and keep all work isolated to that branch
2. Perform an audit-first review of the current repo state, requirements, and impacted surfaces before making changes
3. Implement the minimal scoped changes required by the task in the branch
4. Run one consolidated QA batch after implementation, capturing the results once
5. Open a PR with the audit summary, change summary, QA evidence, and task linkage for review
