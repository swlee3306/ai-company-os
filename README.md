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

## Task T-20260312-090900

### Acceptance criteria

- A dedicated branch is created for the task before implementation starts
- All implementation work is proposed and reviewed through a pull request; no direct mainline delivery
- An audit-first pass is completed before code changes, with findings captured in the task or PR notes
- QA is executed once as a single batch after implementation is complete
- The task is considered done only after branch, PR, audit, and single-batch QA requirements are all satisfied

### Execution plan

1. Create a task branch for `T-20260312-090900
2. Perform the audit-first review and record findings before making changes
3. Implement the approved scope on the task branch
4. Open or update the PR with the audit summary and implementation details
5. Run QA one time as a single batch for the completed change set
6. Address blocking QA issues on the same branch, update the PR, and merge through the PR flow when accepted
