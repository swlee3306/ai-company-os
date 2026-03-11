# P1-2 Acceptance Criteria — Evidence bundle typing + UI improvements

Goal: Make evidence bundles more readable by typing entries and presenting them in a structured UI.

## In Scope
- Treat artifacts as first-class evidence where possible.
- Display evidence bundle entries grouped by type.
- Improve Project Detail evidence section UI (chips/badges/sections).

## Out of Scope
- Real CI integration.
- Automatic artifact extraction from git.

## Rules (MVP)
- If evidence entry is an artifact id (`A-*`): resolve to artifact detail and use its `type`.
- Otherwise infer type heuristically:
  - `ADR-` → ADR
  - `schema-diff` → schema-diff
  - `benchmark` → benchmark
  - `runbook` or `rollback` → runbook
  - else → other

## Acceptance Criteria

### AC-P12-FE-001: Evidence bundle grouping
- On Project Detail:
  - Evidence entries are grouped into sections by type.
  - Each entry is rendered as a compact badge/chip with label.

### AC-P12-FE-002: Links
- Artifact-backed entries link to `/artifacts/:id`.
- Non-artifact entries remain plain text but visually typed.

### AC-P12-FE-003: Empty state
- If no evidence bundle exists, show a clear placeholder.

### AC-P12-QA-001: Smoke
- With seeded project evidence bundle entries:
  - grouping renders
  - at least one artifact-backed link works
