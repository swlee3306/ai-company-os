package model

type Agent struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	PersonaRole      string   `json:"persona_role"`
	OpsSpecialty     string   `json:"ops_specialty,omitempty"`
	Status           string   `json:"status"`
	Scope            []string `json:"scope,omitempty"`
	Version          string   `json:"version"`
	HeartbeatSeconds int      `json:"heartbeat_seconds"`
	ApprovalRequired bool     `json:"approval_required"`
	RiskScope        []string `json:"risk_scope,omitempty"`
}

type Project struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Phase    string   `json:"phase"`
	OwnerCEO string   `json:"owner_ceo"`
	TeamLead string   `json:"team_lead"`
	Due      string   `json:"due"`
	Summary  string   `json:"summary"`
	Evidence []string `json:"evidence_bundle,omitempty"`
	Agents   []string `json:"agents,omitempty"`
}

type ApprovalItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Requester string `json:"requester"`
	Target    string `json:"target"`
	Risk      string `json:"risk"`
	Action    string `json:"action"`
	Status    string `json:"status"`
}
