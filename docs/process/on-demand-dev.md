# On-demand Dev Company Process

This repo is operated like an on-demand dev company.

## Roles
- CEO (human): sets priorities and makes final decisions.
- Team Lead (orchestrator): breaks work into tasks, assigns roles, integrates outputs.
- PM, Designer, FE, BE, QA, Reviewer: on-demand execution roles.

## Default Rule: No direct main pushes
- Work is done on a branch.
- QA runs once per batch (smoke + regression scope appropriate to the change).
- Reviewer approves (or rejects with reason).
- Only after approval: merge into `main`.

## Minimal Flow
1) CEO intent → Task (PM) with acceptance criteria
2) Design (if needed) → `.pen` + exports
3) FE/BE implementation on a branch
4) QA validation (batch)
5) Reviewer gate (approve/reject)
6) Merge to `main`

## Evidence
- Every meaningful change should have:
  - acceptance criteria reference
  - audit/log evidence when applicable
  - screenshots/exports when UI/design changes
