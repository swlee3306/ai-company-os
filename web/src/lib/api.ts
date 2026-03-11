const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

export type Status = {
  driver?: string;
  health?: string;
  note?: string;
};

export type Agent = {
  id: string;
  name: string;
  persona_role: string;
  ops_specialty?: string;
  status: string;
  scope?: string[];
  version: string;
  heartbeat_seconds: number;
  approval_required: boolean;
  risk_scope?: string[];
};

export type Project = {
  id: string;
  name: string;
  status: string;
  phase: string;
  owner_ceo: string;
  team_lead: string;
  due: string;
  summary: string;
  evidence_bundle?: string[];
  agents?: string[];
};

export type ApprovalItem = {
  id: string;
  type: string;
  requester: string;
  target: string;
  risk: string;
  action: string;
  status: string;
};

async function getJson<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`);
  if (!res.ok) throw new Error(`${path} fetch failed: ${res.status}`);
  return res.json();
}

export async function getStatus(): Promise<Status> {
  return getJson('/api/status');
}

export async function getDoctor(): Promise<any> {
  return getJson('/api/doctor');
}

export async function getAudit(): Promise<string> {
  const res = await fetch(`${API_BASE}/api/audit`);
  if (!res.ok) throw new Error(`audit fetch failed: ${res.status}`);
  return res.text();
}

export async function listAgents(): Promise<Agent[]> {
  return getJson('/api/agents');
}

export async function listProjects(): Promise<Project[]> {
  return getJson('/api/projects');
}

export async function listApprovals(): Promise<ApprovalItem[]> {
  return getJson('/api/approvals');
}
