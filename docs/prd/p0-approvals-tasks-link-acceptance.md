# P0-2 Acceptance Criteria — Link Approvals ↔ Tasks

Goal: Make approvals actionable by linking them to specific tasks and reflecting approval decisions in task state.

## In Scope
- Add `task_id` field to approval items (where applicable).
- When an approval decision is made, update the linked task state when relevant.
- Ensure audit trail captures both the approval decision and resulting task transition.
- Minimal UI linkage: show task_id in Approval evidence panel when present.

## Out of Scope
- Full reviewer workflow gating for all task transitions.
- Complex policy engine.

## Rules (MVP)
- If an approval has `task_id` and decision is `approve`:
  - if task state is `blocked`, transition it to `planned`.
- If decision is `reject`:
  - task remains `blocked`.

## Acceptance Criteria

### AC-P0T-BE-001: approval model stores task_id
- Approval items include an optional `task_id`.
- Seed data includes at least one approval item linked to a task.

### AC-P0T-BE-002: decision endpoint applies task transition
- On `POST /api/approvals/:id/decision`:
  - if approval has `task_id` and decision=approve, task transitions blocked→planned.
  - if decision=reject, no task transition is performed.

### AC-P0T-AUD-001: audit evidence
- Approval decision emits `approvals.decision`.
- If a task transition occurs due to approval, emit `task.transition` with `from/to` and `cause=approval`.

### AC-P0T-API-001: tasks reflect change
- After approve decision, `GET /api/tasks/:id` reflects new state.

### AC-P0T-FE-001: approvals UI shows linkage
- In Approvals evidence panel, show `task_id` when present.

### AC-P0T-QA-001: smoke
- Create or seed:
  - a blocked task
  - an approval pointing to that task
- Approve it and verify:
  - task becomes planned
  - audit includes both events
