# Company Platform UI

Design source: `design/company-ui.pen`

The five required frames were created in the Pencil document:
- `Login`
- `Dashboard`
- `Agent Registry`
- `Task Board`
- `Approval Center`

PNG export paths:
- [Login](../../design/exports/login.png)
- [Dashboard](../../design/exports/dashboard.png)
- [Agent Registry](../../design/exports/agent-registry.png)
- [Task Board](../../design/exports/task-board.png)
- [Approval Center](../../design/exports/approval-center.png)

Notes:

## Login
- States: default sign-in state with SSO primary action and backup-code secondary action.
- Components: card shell, input groups, info alert, status labels, primary and outline buttons.

## Dashboard
- States: overview state with metrics, active navigation, action toolbar, warning context, and a detail rail.
- Components: sidebar, metric cards, action card, image card, data table, warning alert, confirmation modal.

## Agent Registry
- States: searchable index state with registration CTA and a policy-review notice on the side rail.
- Components: sidebar, search box, primary button, data table, info alert, image card, action card.

## Task Board
- States: kanban-style flow across queued, running, and blocked work, with blocked items elevated visually.
- Components: sidebar, tabs, primary button, card variants, warning/error alert, modal card.

## Approval Center
- States: reviewer workspace with queue filtering, warning banner, stacked approval items, and side summary.
- Components: sidebar, search box, outline button, warning alert, dialog card, center modal, summary card.

Export note:
- Screens were exported from Pencil as a ZIP and extracted into `design/exports/`.
- The linked PNG files now point to the exported images on disk.
