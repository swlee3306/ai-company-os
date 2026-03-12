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

- A dedicated implementation branch is created for Task `T-20260312-093559`; direct commits to the default branch are not used
- A single pull request is opened for the task and contains the full MVP scope
- The change set is limited to `README.md
- The delivery includes an audit-first review of the current repository state and task constraints before implementation
- QA is executed once as a single batch before merge approval
- The pull request description includes the task ID, scope limit, acceptance criteria, execution plan, and QA result

### Execution plan

1. Audit the current repository state, task scope, and delivery constraints before making changes
2. Create a task branch for `T-20260312-093559
3. Update only `README.md` with the required PM artifacts for this task
4. Run one QA batch covering the README-only change and record the result
5. Open a pull request from the task branch and include the audit summary, acceptance criteria, execution plan, and QA outcome
