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

## Task T-20260312-093559

### Acceptance criteria

- Document work is delivered on a dedicated branch and proposed via PR only; no direct commits to the default branch
- The implementation path is audit-first: inspect current behavior and constraints before making changes in the pipeline
- QA is executed as a single batched pass after implementation work is complete
- The `full` pipeline execution plan is explicitly defined before work begins

### Execution plan

1. Audit the current repository state, pipeline expectations, and release constraints before changing implementation scope
2. Create a dedicated task branch for `T-20260312-093559` and keep all work isolated to that branch
3. Execute the `full` pipeline changes required by the task after the audit confirms scope and dependencies
4. Run QA once as a single batch covering the completed branch state
5. Open a PR with the audit summary, implementation summary, and QA results for review and merge
