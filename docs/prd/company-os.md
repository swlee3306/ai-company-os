# PRD: AI Dev Company OS

## Document Control
- Product: AI Dev Company OS
- Document owner: Manager / Team Lead helper
- Status: Draft for design and implementation alignment
- Last updated: 2026-03-10
- Related design spec: `docs/design/company-platform-ui.md` (required companion document)
- Related design artifacts: `design/company-ui.pen`, `design/exports/` (`design/app.pen` is currently an unused placeholder)

## 1. Product Summary
AI Dev Company OS is an installable dev-company platform that runs with one operating model across macOS, Linux, and WSL. It combines a local/company CLI, an operational dashboard, and runtime drivers so a CEO, orchestrator, workers, and reviewers can install, launch, monitor, and audit AI-assisted software delivery with the same concepts, commands, and evidence trail on every supported environment.

The product must support a default local Kubernetes driver (`k3d`) and an optional lightweight driver (`k3s`) while keeping the user-facing workflow stable. The platform must also expose audit and notification integrations so operational state changes are visible beyond the local machine.

## 2. Problem Statement
Current AI-agent workflows are fragmented:
- Setup differs across macOS, Linux, and WSL, increasing onboarding time and failure rate.
- Operators lack a single control plane for installation, lifecycle management, task execution, and diagnostics.
- Leadership cannot reliably answer basic questions such as "what is running", "what failed", and "what evidence exists".
- Worker and reviewer roles often operate without a shared audit trail, increasing ambiguity and rework.
- UI and CLI are frequently designed separately, causing gaps between command semantics and dashboard representations.

## 3. Personas
### CEO (Flant)
- Needs confidence that the company OS is installable, observable, and auditable.
- Needs high-level status, throughput, risk indicators, and notification hooks.
- Cares about adoption, operational reliability, and delivery velocity.

### Team Lead / Orchestrator (openclaw)
- Owns environment setup, cluster driver selection, worker coordination, and exception handling.
- Converts CEO intent into structured tasks and assigns work across roles.
- Needs deterministic commands, actionable health checks, and a dashboard that maps directly to CLI actions.
- Cares about task flow, agent utilization, and recovery from broken states.

### Planner / PM (Product)
- Translates CEO intent into requirements, scope (in/out), and acceptance criteria.
- Identifies decision points and edge cases with minimal, high-signal questions.
- Keeps PRD and delivery plan aligned with what is actually being built.

### Designer (Product UI)
- Turns PRD outcomes into screen/component design and interaction states.
- Enforces layout rules (grid/spacing/tokens) and prevents clipping/readability issues.
- Ensures primary actions remain semantically aligned with CLI/system operations.

### FE Developer (Dashboard / Client)
- Implements UI behavior and maps primary actions to CLI/system operations.
- Owns state handling (empty/loading/error), usability, and regression prevention.
- Keeps changes minimal and consistent with design tokens and component rules.

### BE Developer (CLI / Runtime)
- Implements command handlers, task state model, logs/audit generation, and notifications.
- Ensures outputs are evidence-grade (traceable IDs, timestamps, actor/action records).
- Owns cross-platform execution constraints (dependencies, permissions, driver differences).

### QA (Validation)
- Defines test scenarios and pass/fail gates across macOS, Linux, and WSL.
- Validates install/runtime/task flows and checks for regressions.
- Produces reproducible evidence for failures and verifies fixes.

### Workers
- Execute assigned work units and need clear task state, logs, dependencies, and environment readiness.
- Care about low-friction task pickup, minimal setup ambiguity, and quick debugging.

### Reviewer
- Verifies outputs against task intent and design intent.
- Needs traceable evidence, audit logs, task lineage, and status visibility without inspecting internals manually.
- Acts as a release gate: approve or reject based on evidence and acceptance criteria.

## 4. Goals
### Product Goals
1. Provide a single installable platform experience across macOS, Linux, and WSL.
2. Offer one operational vocabulary shared by CLI and dashboard.
3. Make runtime state, task state, logs, and diagnostics visible and auditable.
4. Support default `k3d` operation with optional `k3s` selection without changing higher-level workflows.
5. Enable leadership visibility through audit records and Discord notifications.

### Non-Goals
1. Building a cloud-hosted multi-tenant SaaS in this phase.
2. Supporting every Kubernetes distribution or every operating system variant.
3. Replacing existing source control, CI, or chat tools.
4. Designing bespoke worker IDEs beyond the required dashboard surfaces.
5. Defining low-level implementation details for every backend service in this PRD.

## 5. Success Metrics
### Primary Metrics
- `Install success rate`: At least 90% successful first-run installs on supported environments in controlled validation.
- `Time to operational readiness`: Median time from `install` start to healthy `status` under 15 minutes on reference machines.
- `Task execution success`: At least 95% of valid task submissions reach terminal state with retained logs and audit entries.
- `Diagnostic usefulness`: At least 80% of test users report `doctor` output is sufficient to self-resolve common setup/runtime issues.

### Secondary Metrics
- `Command/UI consistency`: 100% of core dashboard actions map to defined CLI commands.
- `Audit coverage`: 100% of lifecycle-changing actions create an audit event.
- `Notification coverage`: 100% of failed installs, failed task runs, and cluster-down events can trigger Discord notifications.

## 6. Product Principles
1. CLI-first truth: the CLI defines the canonical operational model.
2. Dashboard mirrors, not invents: UI state and actions must map to underlying CLI concepts.
3. Diagnose before blame: every failure path should expose next steps.
4. Evidence over assumptions: audit logs, task logs, and health checks must be preserved and attributable.
5. Cross-platform parity: user intent should not change across supported host environments.

## 6.1 Default Role Tooling And Outputs
Unless explicitly overridden by the installed user or Team Lead, the system follows these defaults:

- Planner / PM: produces Markdown artifacts (requirements, scope, acceptance criteria, open questions, risks/trade-offs).
- Designer: uses Pencil (`pencil.dev` MCP) as the default design toolchain and maintains `.pen` source + exported PNG artifacts.
- FE Developer: defaults to TypeScript; framework choice should follow modern conventions (React + Vite by default) unless specified otherwise.
- BE Developer: defaults to TypeScript (Node.js) for stack consistency; alternative languages/frameworks are allowed when specified.
- QA: executes unit → integration → (where feasible) E2E validation and produces a reproducible evidence trail.
- Reviewer: confirms alignment across PRD intent, design intent, implementation, and QA evidence; approves or rejects with a recorded reason.

## 7. Scope
### In Scope
- CLI lifecycle commands: `install`, `up`, `down`, `status`, `task`, `logs`, `doctor`
- Driver support: `k3d` by default, `k3s` optional
- Dashboard requirements for five core screens
- Audit event generation and retention expectations
- Discord notification requirements
- Traceability from PRD requirements to UX and future implementation/testing work

### Out of Scope
- Billing, authentication federation, and enterprise policy administration
- Full RBAC matrix beyond role-based persona expectations
- Marketplace/plugin ecosystems
- Detailed data model/schema contracts for downstream implementation teams

## 8. User Stories And Acceptance Criteria
### 8.1 Installation And Environment
**User story:** As a Team Lead, I can install the company OS on macOS, Linux, or WSL so that the same workflow is available on each platform.

Acceptance criteria:
- A documented install flow exists for macOS, Linux, and WSL.
- `install` validates prerequisites and reports actionable failures.
- `install` defaults to `k3d` unless the user selects `k3s`.
- Successful install produces a clear next step to run `up` or view status in the dashboard.
- Install completion creates an audit record.

### 8.2 Runtime Lifecycle
**User story:** As a Team Lead, I can start, stop, and inspect the platform using stable commands and matching UI actions.

Acceptance criteria:
- `up` starts all required local services for the selected driver.
- `down` stops services cleanly and leaves the environment recoverable.
- `status` reports driver, service health, task queue state, and last known issues.
- Dashboard lifecycle controls match `up`, `down`, and `status` semantics.
- Lifecycle transitions create audit records and can trigger Discord notifications on failure.

### 8.3 Task Operations
**User story:** As a worker or orchestrator, I can submit and inspect tasks so work moves through the system with visible state and logs.

Acceptance criteria:
- `task` supports submitting, listing, and inspecting tasks.
- Each task has an identifiable state model from queued to terminal.
- `logs` can retrieve task-scoped and system-scoped logs.
- Dashboard task views show the same task identifiers, status, and log availability as the CLI.
- Task state changes generate audit events.

### 8.4 Diagnostics
**User story:** As a Team Lead, I can run diagnostics to identify broken setup, runtime, or dependency conditions.

Acceptance criteria:
- `doctor` checks host prerequisites, driver availability, cluster connectivity, storage/network assumptions, and notification configuration health.
- `doctor` groups findings by severity and suggests remediation steps.
- Dashboard surfaces the latest diagnostic summary and links users to more detailed evidence.
- Diagnostic runs create audit records.

### 8.5 Executive Visibility
**User story:** As the CEO, I can see system health, throughput, and exceptions without using the CLI directly.

Acceptance criteria:
- Dashboard provides an executive summary view with cluster health, task throughput, worker utilization, and active incidents.
- Critical failures and task exceptions can be sent to Discord.
- Executive metrics reference auditable source data rather than inferred UI-only state.

### 8.6 Review And Evidence
**User story:** As a Reviewer, I can verify delivered work against task intent and design intent from one evidence trail.

Acceptance criteria:
- Audit data records actor, action, timestamp, target, and outcome.
- Task detail views expose linked logs and relevant audit entries.
- The PRD and future implementation tasks reference `docs/design/company-platform-ui.md` for UI intent.
- Verification guidance includes the repository-standard command `PYTEST_DISABLE_PLUGIN_AUTOLOAD=1 VERIFY_TIMEOUT_SECONDS=900 ./scripts/verify.sh`.

### 8.7 Service Delivery Loop (Role Negotiation + Iteration)
**User story:** As an installed user, I can request that the Team Lead build a service (web app, API, automation workflow, or mixed system) so the system coordinates PM, Design, FE, BE, QA, and Reviewer roles through iterative negotiation until the service is delivered.

Acceptance criteria:
- A service request can be submitted as a structured task with: goal, scope constraints, and a definition of done.
- The Team Lead can assign roles (PM/Design/FE/BE/QA/Reviewer) to a request, with accountability recorded.
- The system supports iterative loops: a reviewer rejection can route work back to the appropriate role stage with a reason.
- The system preserves a decision trail for trade-offs (what was changed, why, and by whom) alongside logs and audit events.
- Completion produces an evidence bundle referencing requirements, commits/artifacts, test results, and review approval.

## 9. Functional Requirements
### 9.1 CLI Requirements
- `REQ-CLI-001`: The product shall provide a unified CLI entry point for platform operations.
- `REQ-CLI-002`: The CLI shall support `install`, `up`, `down`, `status`, `task`, `logs`, and `doctor`.
- `REQ-CLI-003`: Command help output shall use consistent naming and examples across supported platforms.
- `REQ-CLI-004`: Non-zero exit codes shall be used for failed operations.
- `REQ-CLI-005`: Long-running commands shall stream progress or status updates rather than remaining silent.

### 9.2 Driver Requirements
- `REQ-DRV-001`: `k3d` shall be the default runtime driver.
- `REQ-DRV-002`: `k3s` shall be selectable as an optional runtime driver.
- `REQ-DRV-003`: Driver selection shall be visible in both CLI status output and dashboard system views.
- `REQ-DRV-004`: Higher-level task and lifecycle workflows shall remain stable regardless of selected driver.

### 9.3 Task And Logging Requirements
- `REQ-TASK-001`: The system shall support task submission, task listing, and task inspection.
- `REQ-TASK-002`: The system shall persist task state transitions and timestamps.
- `REQ-TASK-003`: Each task shall support role assignment metadata (PM, Design, FE, BE, QA, Reviewer) and a recorded owner per role.
- `REQ-TASK-004`: The system shall support a task stage model that enables iteration (e.g., planned → designing → implementing → testing → review → done/rejected).
- `REQ-TASK-005`: When a task is rejected, the system shall record the rejection reason and route the task back to a prior stage.
- `REQ-LOG-001`: The system shall expose task logs and system logs.
- `REQ-LOG-002`: Logs shall be accessible from both CLI and dashboard contexts.
- `REQ-EVID-001`: Completed tasks shall produce an evidence bundle referencing acceptance criteria, implementation artifacts (commits/build outputs), test results, and review approval.
- `REQ-DEC-001`: The system shall support a decision/trade-off log linked to a task (who decided, what changed, and why).

### 9.4 Audit And Notification Requirements
- `REQ-AUD-001`: Every lifecycle-changing operation shall emit an audit record.
- `REQ-AUD-002`: Audit records shall include actor, action, target, timestamp, and result.
- `REQ-NOTIFY-001`: The system shall support Discord notifications for failures and notable operational events.
- `REQ-NOTIFY-002`: Notification delivery state or configuration health shall be visible in diagnostics.

### 9.5 Diagnostics Requirements
- `REQ-DOC-001`: The system shall provide a `doctor` command for environment and runtime diagnostics.
- `REQ-DOC-002`: Diagnostic output shall classify findings by severity and remediation guidance.
- `REQ-DOC-003`: Dashboard system health views shall surface the latest diagnostic summary.

## 10. UX/UI Requirements
The design team is producing five screens. The final design specification must be documented in `docs/design/company-platform-ui.md`, and exported screenshots must be placed under `design/exports/`. The PRD defines required user outcomes and minimum component expectations for each screen.

### Screen 1: Installation / Onboarding
- Purpose: guide first-time setup and prerequisite validation.
- Required UX outcomes:
  - Show supported host environments: macOS, Linux, WSL.
  - Show driver selection with `k3d` default and `k3s` optional.
  - Show install progress, failure states, and remediation guidance.
- Required components:
  - Prerequisite checklist
  - Driver selector
  - Progress/state panel
  - Next-step CTA

### Screen 2: Operations Dashboard
- Purpose: give CEO and Team Lead a real-time control plane.
- Required UX outcomes:
  - Surface current cluster/driver health, worker status, queue status, and active incidents.
  - Expose lifecycle actions matching `up`, `down`, and `status`.
  - Highlight audit/notification status.
- Required components:
  - System summary cards
  - Lifecycle action bar
  - Incident/status banner
  - Worker/runtime status panel

### Screen 3: Task Board / Task Detail
- Purpose: let orchestrators and workers submit, monitor, and inspect tasks.
- Required UX outcomes:
  - Visualize task states and ownership clearly.
  - Provide task detail with logs, timestamps, and audit linkage.
  - Keep identifiers consistent with CLI output.
- Required components:
  - Task list or board
  - Task detail drawer/page
  - Status timeline
  - Log preview or linked log panel

### Screen 4: Logs And Audit
- Purpose: support reviewers and operators investigating what happened.
- Required UX outcomes:
  - Differentiate task logs from system logs and audit events.
  - Allow filtering by actor, task, severity, and time.
  - Preserve a clear chain from action to evidence.
- Required components:
  - Filter controls
  - Log viewer
  - Audit event table
  - Event detail panel

### Screen 5: Diagnostics / Settings
- Purpose: help Team Leads recover from broken states and maintain integrations.
- Required UX outcomes:
  - Show latest `doctor` summary and remediation guidance.
  - Expose Discord configuration health and audit retention status.
  - Make driver/environment settings visible and editable where appropriate.
- Required components:
  - Diagnostic summary panel
  - Findings list with severity
  - Integration status section
  - Environment/settings form

## 11. Design Coordination Requirements
- `REQ-UX-001`: `docs/design/company-platform-ui.md` shall exist and describe the five required screens.
- `REQ-UX-002`: The design spec shall state how each primary UI action maps to a CLI command or system action.
- `REQ-UX-003`: The design spec shall name core components used on each screen.
- `REQ-UX-004`: Exported screenshots shall be stored under `design/exports/`.
- `REQ-UX-005`: The following screenshot artifacts are required at minimum:
  - `design/exports/install-onboarding.png`
  - `design/exports/operations-dashboard.png`
  - `design/exports/task-board-detail.png`
  - `design/exports/logs-audit.png`
  - `design/exports/diagnostics-settings.png`
- `REQ-UX-006`: If naming differs, `docs/design/company-platform-ui.md` shall include an explicit filename mapping table.

## 12. System Overview
### 12.1 Conceptual Flow
1. User installs the platform via CLI.
2. User selects or accepts default driver (`k3d`, optional `k3s`).
3. User starts the platform and verifies health.
4. Orchestrator submits and monitors tasks.
5. Workers execute tasks with logs and state changes recorded.
6. Reviewers inspect outputs, logs, and audit history.
7. Failures and important events can notify Discord recipients.

### 12.2 Core Interfaces
- CLI:
  - `install`: install prerequisites, validate host, configure driver
  - `up`: start runtime services
  - `down`: stop runtime services
  - `status`: report health and runtime state
  - `task`: submit/list/inspect work
  - `logs`: inspect task/system logs
  - `doctor`: diagnose environment and runtime issues
- Drivers:
  - `k3d` default local Kubernetes runtime
  - `k3s` optional lightweight runtime
- Observability:
  - Audit event stream
  - Task/system logs
  - Discord notification channel
- Dashboard:
  - Five-screen operational UI aligned to CLI semantics

## 13. Risks And Assumptions
### Assumptions
- The initial implementation targets local or single-node operator-managed environments.
- `k3d` and `k3s` are sufficient for the first release.
- Discord is the only required outbound notification target in this phase.
- Design work will land in `docs/design/company-platform-ui.md` and screenshots will be exported into `design/exports/`.

### Risks
- Cross-platform installation differences may create inconsistent prerequisite handling.
- Driver-specific edge cases may leak through the intended unified workflow.
- Missing audit or notification instrumentation can undermine reviewer and executive trust.
- UI/CLI drift may occur if dashboard flows are designed without direct command mapping.

## 14. Traceability Matrix
This matrix links PRD requirements to the required UX surfaces and future implementation/testing work.

| Requirement | UX screen / component target | Future task focus | Test focus |
| --- | --- | --- | --- |
| REQ-CLI-001, REQ-CLI-002 | Screen 2 lifecycle action bar; Screen 3 task controls; Screen 5 diagnostics actions | Implement canonical CLI entry point and command handlers | CLI smoke tests for command availability and help output |
| REQ-CLI-005 | Screen 1 progress/state panel; Screen 2 status banner | Progress streaming and live status updates | Long-running command output tests |
| REQ-DRV-001, REQ-DRV-002, REQ-DRV-003 | Screen 1 driver selector; Screen 2 system summary; Screen 5 settings form | Driver abstraction and configuration persistence | Driver selection and status reporting tests |
| REQ-TASK-001, REQ-TASK-002 | Screen 3 task board/detail; Screen 2 queue indicators | Task model, state transitions, CLI/UI task synchronization | Task lifecycle integration tests |
| REQ-LOG-001, REQ-LOG-002 | Screen 3 log preview; Screen 4 log viewer | Unified log retrieval and retention | Log access tests for task and system scopes |
| REQ-AUD-001, REQ-AUD-002 | Screen 4 audit event table/detail; Screen 2 audit status summary | Audit event schema and event emission | Audit coverage tests for lifecycle and task events |
| REQ-NOTIFY-001, REQ-NOTIFY-002 | Screen 2 incident banner; Screen 5 integration status | Discord integration and delivery status reporting | Notification trigger and configuration-health tests |
| REQ-DOC-001, REQ-DOC-002, REQ-DOC-003 | Screen 5 diagnostic summary/findings; Screen 2 health summary | Diagnostic engine and surfaced remediation guidance | Doctor command output and UI summary tests |
| REQ-UX-001 through REQ-UX-006 | All five screens and exported screenshots | Design documentation and design-export delivery | Design artifact presence checks and reviewer checklist |

## 15. Future Task And Test Guidance
Implementation and review tasks derived from this PRD should:
- Reference the specific requirement IDs above.
- Reference the relevant section in `docs/design/company-platform-ui.md`.
- Name the target files/components being changed.
- Include validation with:
  - `PYTEST_DISABLE_PLUGIN_AUTOLOAD=1 VERIFY_TIMEOUT_SECONDS=900 ./scripts/verify.sh`
- Add focused command-level or UI-level checks where applicable.

Recommended future workstreams:
1. CLI and driver implementation
2. Audit and notification pipeline
3. Dashboard shell and five-screen UI
4. Task/log data plumbing
5. Cross-platform installation and diagnostics hardening

## 16. Open Dependencies
- `docs/design/company-platform-ui.md` must be created or updated to match this PRD.
- Required screenshots must be exported under `design/exports/`.
- Downstream implementation tasks must attach tests mapped to the requirement IDs in Section 14.
