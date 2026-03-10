# Company Platform UI

Design source artifacts:
- [Pencil board (`company-ui.pen`)](../../design/company-ui.pen)
- `design/app.pen` is currently an unused placeholder (may be removed or populated later).
- [Export folder](../../design/exports/)

## Goals
- Provide a five-screen UI spec aligned to PRD Section 10 for install, operations, tasks, logs/audit, and diagnostics/settings.
- Keep dashboard actions semantically aligned with CLI/system operations (`install`, `up`, `down`, `status`, `task`, `logs`, `doctor`).
- Ensure exported design artifacts are reviewable from repository paths and filename-compliant with PRD `REQ-UX-005`.

## Non-Goals
- Define backend implementation details for command handlers, persistence schema, or notification transport internals.
- Replace CLI flows; this spec mirrors CLI/system semantics in UI form.
- Specify enterprise IAM/RBAC or non-Discord notification providers for this phase.

## User Flows
1. Team Lead opens **Installation / Onboarding**, validates prerequisites, confirms driver (`k3d` default or `k3s`), runs install, then transitions to operations.
2. Team Lead or CEO uses **Operations Dashboard** to start/stop runtime (`up`/`down`) and inspect current platform state (`status`).
3. Orchestrator/Worker uses **Task Board / Detail** to submit tasks, track state transitions, and inspect task details/log previews.
4. Reviewer/Operator uses **Logs And Audit** to filter and inspect system logs, task logs, and audit events.
5. Team Lead uses **Diagnostics / Settings** to run diagnostics (`doctor`), review remediation guidance, and verify integration/settings health.

## Screens List (PRD Required)
1. **Installation / Onboarding**  
   Artifact: [install-onboarding.png](../../design/exports/install-onboarding.png)
2. **Operations Dashboard**  
   Artifact: [operations-dashboard.png](../../design/exports/operations-dashboard.png)
3. **Task Board / Detail**  
   Artifact: [task-board-detail.png](../../design/exports/task-board-detail.png)
4. **Logs And Audit**  
   Artifact: [logs-audit.png](../../design/exports/logs-audit.png)
5. **Diagnostics / Settings**  
   Artifact: [diagnostics-settings.png](../../design/exports/diagnostics-settings.png)

## Screen Specs
### Installation / Onboarding
- Purpose: first-run setup and prerequisite validation across macOS, Linux, WSL.
- Components: prerequisite checklist, driver selector (`k3d` default, `k3s` optional), install progress/state panel, next-step CTA.
- Key states:
  - Empty: no install started; prerequisites pending.
  - Loading: install/probe in progress with streamed step status.
  - Error: prerequisite or install failure with remediation guidance.

### Operations Dashboard
- Purpose: operational control plane for runtime health and lifecycle.
- Components: system summary cards, lifecycle action bar, incident/status banner, worker/runtime status panel, audit/notification status summary.
- Key states:
  - Empty: no active runtime; prompt to run `up`.
  - Loading: status refresh underway.
  - Error: runtime degraded/down with surfaced incident context.

### Task Board / Detail
- Purpose: task submission, queue visibility, and deep task inspection.
- Components: task board/list, task detail drawer/panel, status timeline, log preview, audit linkage block.
- Key states:
  - Empty: no queued/running tasks.
  - Loading: task list/detail refresh or submit in progress.
  - Error: failed task or unavailable task/log payload.

### Logs And Audit
- Purpose: investigation workspace with clear evidence chain.
- Components: filter controls (actor/task/severity/time), log viewer (task/system), audit event table, event detail panel.
- Key states:
  - Empty: no records in selected filter scope.
  - Loading: log/audit query in progress.
  - Error: log backend unavailable or filter query failed.

### Diagnostics / Settings
- Purpose: diagnose and recover environment/runtime issues; maintain operational settings.
- Components: diagnostic summary panel, findings list (severity + remediation), integration status section (Discord/audit retention), environment/settings form.
- Key states:
  - Empty: no diagnostic run yet.
  - Loading: `doctor` execution in progress.
  - Error: diagnostics command failure or stale/unavailable health data.

## REQ-UX-002: Primary UI Action -> CLI/System Mapping
| Screen | Primary UI action | CLI/system command | Expected system effect |
| --- | --- | --- | --- |
| Installation / Onboarding | Run install | `company install` | Validates prerequisites, applies driver config, bootstraps platform, emits install audit record. |
| Installation / Onboarding | Select driver | System config action (`driver=k3d` default or `driver=k3s`) used by `install`/runtime | Persists runtime driver selection surfaced later by `status`. |
| Operations Dashboard | Start platform | `company up` | Starts runtime services and updates health/incident status. |
| Operations Dashboard | Stop platform | `company down` | Stops runtime services cleanly and records lifecycle event. |
| Operations Dashboard | Refresh status | `company status` | Returns driver, health, queue, and issue summary for dashboard cards/banner. |
| Task Board / Detail | Submit task | `company task "..."` | Creates queued task with ID and initial audit/log trail. |
| Task Board / Detail | Inspect task | `company task inspect <task-id>` | Returns task metadata, owner/state timeline, and linked evidence. |
| Task Board / Detail | List tasks | `company task list` | Returns filtered task queue/board data. |
| Logs And Audit | View logs | `company logs --task <task-id>` or `company logs --system` | Retrieves task-scoped or system-scoped logs. |
| Logs And Audit | View audit events | System audit query action (dashboard audit feed) | Retrieves actor/action/target/timestamp/result evidence records. |
| Diagnostics / Settings | Run diagnostics | `company doctor` | Checks prerequisites/runtime/integrations and returns severity-classified remediation guidance. |
| Diagnostics / Settings | Validate notifications/settings | System integration health check surfaced with `doctor` output | Shows Discord configuration and delivery health indicators. |

## Component Inventory
- Prerequisite checklist
- Driver selector
- Progress/state panel
- Next-step CTA/action strip
- Summary metric cards
- Lifecycle action bar
- Incident/status banner
- Worker/runtime status panel
- Task board/list and task cards
- Task detail drawer/panel
- Status timeline
- Log preview/log viewer
- Filter controls
- Audit table and event detail panel
- Diagnostic summary and findings list
- Integration status panel
- Environment/settings form

## Interaction States
Across all screens, the minimum interaction states are defined and represented:
- Empty: first-run or no-data condition with recovery CTA.
- Loading: command/query in progress with explicit activity indicator/progress messaging.
- Error: actionable failure state with remediation guidance and retry path.

## Responsive Rules
- Desktop (>= 1200 px): two-rail layouts are allowed (primary content + side detail panel).
- Tablet (768 px to 1199 px): collapse secondary panels into tabs/stacked sections.
- Mobile (< 768 px): single-column flow; action bar becomes sticky bottom/top action group; tables convert to stacked cards with key-value rows.

## Layout and Spacing Rules
- All five PRD screens use an 8pt grid baseline; all frame sizes/offsets/gaps align to 8px increments.
- Root screen frame padding: 32px.
- Primary section spacing (vertical rhythm between major blocks): 24px.
- Card-to-card spacing (within columns/rows): 16px.
- Card internal spacing: 16px (padding and internal content gaps standardized to 16px).
- Column gutters are equalized at 16px across each two-column screen composition.
- Header rows are normalized to a shared height and spacing pattern; action button containers and button heights are aligned to keep button baselines consistent.
- Off-by-1 spacing and uneven vertical rhythm are removed by enforcing shared constants (`32 / 24 / 16`) and 8pt multiples across screen layouts.
- Applied consistently to: Installation / Onboarding, Operations Dashboard, Task Board / Detail, Logs And Audit, and Diagnostics / Settings.

## Design Tokens
No new token system is introduced in this revision. Existing Pencil design tokens/styles from `design/company-ui.pen` are reused.

## Export Artifacts (REQ-UX-005)
Canonical screenshot filenames present under `design/exports/`:
- [install-onboarding.png](../../design/exports/install-onboarding.png)
- [operations-dashboard.png](../../design/exports/operations-dashboard.png)
- [task-board-detail.png](../../design/exports/task-board-detail.png)
- [logs-audit.png](../../design/exports/logs-audit.png)
- [diagnostics-settings.png](../../design/exports/diagnostics-settings.png)

## Acceptance Criteria Checklist
- [x] Design spec covers all five PRD-required screens.
- [x] Primary UI actions map to CLI/system actions (`REQ-UX-002`).
- [x] Canonical export filenames required by `REQ-UX-005` are present under `design/exports/`.
- [x] Components and interaction states are specified for each screen.
- [x] Responsive behavior expectations are documented.

## Traceability
### Task Reference
- Task summary: Fix design doc alignment with PRD required screens, add export filename compliance, and add CLI action mapping.
- Source review decision: Prior reviewer decision was `reject`; this revision applies the requested screen alignment, CLI mapping, and canonical export naming fixes.

### Requirement Mapping (Reviewer Checkpoints)
| Requirement | Where satisfied in this document | Artifact/check |
| --- | --- | --- |
| REQ-UX-001 | `Screens List (PRD Required)` + `Screen Specs` | `rg -n "Installation / Onboarding|Operations Dashboard|Task Board / Detail|Logs And Audit|Diagnostics / Settings" docs/design/company-platform-ui.md` |
| REQ-UX-002 | `REQ-UX-002: Primary UI Action -> CLI/System Mapping` table | `rg -n "REQ-UX-002|CLI/system command|company" docs/design/company-platform-ui.md` |
| REQ-UX-003 | `Screen Specs` and `Component Inventory` | Manual reviewer read of component lists |
| REQ-UX-004 | `Interaction States` + `Responsive Rules` | `rg -n "## Interaction States|## Responsive Rules" docs/design/company-platform-ui.md` |
| REQ-UX-005 | `Export Artifacts (REQ-UX-005)` canonical filenames | `ls -lh design/exports/install-onboarding.png design/exports/operations-dashboard.png design/exports/task-board-detail.png design/exports/logs-audit.png design/exports/diagnostics-settings.png` |
| REQ-UX-006 | `Export Artifacts (REQ-UX-005)` links resolve to repository exports | Manual link check in markdown preview |
