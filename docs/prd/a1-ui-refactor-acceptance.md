# A1 Acceptance Criteria — UI refactor to match design intent

Goal: Bring the implemented web UI closer to `design/company-ui.pen` for Agent Detail, Project Detail, and Audit Logs, while preserving current MVP behavior.

## In Scope
1) Agent Detail page UI refactor (structure + labels)
2) Project Detail page UI refactor (structure + labels)
3) Audit Logs: simple filter/search + newest-first ordering

## Out of Scope
- New design creation or layout redesign beyond aligning to existing Pencil spec.
- Full evidence bundle integration (commits/tests/artifacts) — deferred to A2.
- Auth/SSO.

## References
- Design source: `design/company-ui.pen`
  - Frames: `AI Company OS / Agent Detail`, `AI Company OS / Project Detail`, `AI Company OS / Audit Logs`
- PRD: `docs/prd/company-os.md`
- MVP AC baseline: `docs/prd/mvp-acceptance.md`

## Acceptance Criteria

### AC-A1-AGENT-001: Agent Detail page structure
- Agent Detail page renders a clear 3-section layout (cards/sections):
  1) **Configuration and policy**
  2) **Health and recent executions**
  3) **Danger permissions**
- Labels follow the persona-aligned wording used in the design spec:
  - `Persona role` and `Ops specialty` are displayed separately when present.
  - Approval-gated messaging is visible when `approval_required=true`.

### AC-A1-AGENT-002: Agent Detail content mapping
- Agent Detail displays:
  - name/id, persona role, ops specialty (optional)
  - endpoint (if available), scope, version, heartbeat
  - approval_required + risk_scope (if available)
- If the agent is blocked/pending approval, a link/button to Approvals is visible.

### AC-A1-PROJECT-001: Project Detail page structure
- Project Detail page renders sections that match design intent:
  - Project overview
  - Goal / current work / evidence bundle
  - Participating agents and recent work
  - Project memory
  - Key decisions
- Owner and Team Lead are clearly labeled (avoid ambiguous "Manager AI" wording in UI).

### AC-A1-PROJECT-002: Evidence bundle rendering
- If project includes `evidence_bundle`, it renders as a dedicated block labeled **Evidence bundle**.

### AC-A1-AUDIT-001: Audit Logs ordering and parsing
- Audit logs render newest-first.
- JSONL entries are parsed and rendered as rows (ts/actor/action).

### AC-A1-AUDIT-002: Simple filtering
- Provide a search input that filters audit rows by substring match across:
  - actor
  - action
  - (optional) project/target fields when present

### AC-A1-QA-001: Smoke regression
- Existing MVP flows remain functional:
  - navigation routes still work
  - approvals decision loop still works
  - evidence panel still loads

## QA / Verification
- Manual UI check against the Pencil frames.
- Smoke test:
  - start backend + seed
  - start web
  - verify pages render and approvals still function.

## Reviewer Gate
- Reviewer approves only if:
  - UI structure matches the section layout above
  - wording is consistent with design intent
  - no regression in approvals loop
