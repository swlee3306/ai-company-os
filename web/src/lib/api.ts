const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

export type Status = {
  driver?: string;
  health?: string;
  note?: string;
};

export async function getStatus(): Promise<Status> {
  const res = await fetch(`${API_BASE}/api/status`);
  if (!res.ok) throw new Error(`status fetch failed: ${res.status}`);
  return res.json();
}

export async function getDoctor(): Promise<any> {
  const res = await fetch(`${API_BASE}/api/doctor`);
  if (!res.ok) throw new Error(`doctor fetch failed: ${res.status}`);
  return res.json();
}

export async function getAudit(): Promise<string> {
  const res = await fetch(`${API_BASE}/api/audit`);
  if (!res.ok) throw new Error(`audit fetch failed: ${res.status}`);
  return res.text();
}
