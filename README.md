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

## Task T-20260312-091706

### Acceptance criteria

- Delivery is made on a branch and proposed through a PR; direct commits to the main branch are out of scope
- The MVP change set is limited to `README.md
- The plan is audit-first: review the current repository state and documentation before making the single approved edit
- QA runs once as a single batch after the README change is prepared
- The README includes both acceptance criteria and an execution plan for this task

### Execution plan

1. Audit the current repository state and existing `README.md` content to confirm scope and avoid unrelated edits
2. Create a branch for `T-20260312-091706
3. Update `README.md` with the task-specific acceptance criteria and this execution plan only
4. Run one QA batch covering the README diff and any required repository checks
5. Open a PR from the task branch with the audited scope, QA result, and task ID recorded
