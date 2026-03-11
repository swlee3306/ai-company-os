import { useEffect, useState } from 'react';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Settings = {
  driver: { selected: 'k3d' | 'k3s' };
  approval: { policy_text: string };
  workspace?: { repo_path?: string };
  runner?: { backend?: string; type: string; command: string; workdir?: string; agents?: { pm?: string } };
};

export default function Settings() {
  const [data, setData] = useState<Settings | null>(null);
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const [wsStatus, setWsStatus] = useState<{ status: string; detail?: string } | null>(null);
  const [wsValidating, setWsValidating] = useState(false);

  async function refresh() {
    try {
      setError('');
      const res = await fetch(`${API_BASE}/api/settings`);
      if (!res.ok) throw new Error(`settings fetch failed: ${res.status}`);
      setData(await res.json());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function save() {
    if (!data) return;
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/api/settings`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });
      if (!res.ok) throw new Error(`settings save failed: ${res.status}`);
      setData(await res.json());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    } finally {
      setSaving(false);
    }
  }

  async function validateWorkspace() {
    if (!data?.workspace?.repo_path?.trim()) {
      setWsStatus({ status: 'fail', detail: 'repo_path is empty' });
      return;
    }
    setWsValidating(true);
    setWsStatus(null);
    try {
      const res = await fetch(`${API_BASE}/api/workspace/validate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ repo_path: data.workspace.repo_path }),
      });
      if (!res.ok) throw new Error(await res.text());
      const out = await res.json();
      setWsStatus({ status: out.status || (out.ok ? 'ok' : 'fail'), detail: out.detail });
    } catch (e: any) {
      setWsStatus({ status: 'fail', detail: e?.message ?? String(e) });
    } finally {
      setWsValidating(false);
    }
  }

  return (
    <div style={{ maxWidth: 900 }}>
      <h1 style={{ marginTop: 0 }}>Settings</h1>
      <p style={{ color: '#6b7280' }}>Driver selection and approval policy (MVP).</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {!data ? (
        <div>Loading…</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white' }}>
          <div style={{ marginBottom: 14 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Driver</div>
            <select
              value={data.driver.selected}
              onChange={(e) => setData({ ...data, driver: { selected: e.target.value as 'k3d' | 'k3s' } })}
              style={{ padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}
            >
              <option value="k3d">k3d</option>
              <option value="k3s">k3s</option>
            </select>
          </div>

          <div style={{ marginBottom: 14 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Workspace repo path</div>
            <input
              value={data.workspace?.repo_path || ''}
              onChange={(e) => setData({ ...data, workspace: { repo_path: e.target.value } })}
              placeholder="/Users/you/workspace/my-project"
              style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', marginBottom: 8 }}
            />
            <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
              <button onClick={validateWorkspace} disabled={wsValidating} style={{ padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: 'white' }}>
                {wsValidating ? 'Validating…' : 'Validate'}
              </button>
              {wsStatus ? (
                <div style={{ fontSize: 12, color: wsStatus.status === 'ok' ? '#065f46' : wsStatus.status === 'warn' ? '#b45309' : '#b91c1c' }}>
                  {wsStatus.status}{wsStatus.detail ? ` — ${wsStatus.detail}` : ''}
                </div>
              ) : null}
            </div>
            <div style={{ marginTop: 6, color: '#6b7280', fontSize: 12 }}>
              Used by local_cli/codex runner (codex exec -C).
            </div>
          </div>

          <div style={{ marginBottom: 14 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Approval policy</div>
            <div style={{ color: '#374151' }}>{data.approval.policy_text}</div>
          </div>

          <div style={{ marginBottom: 14 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Runner</div>

            <div style={{ marginBottom: 8 }}>
              <div style={{ color: '#6b7280', fontSize: 12, marginBottom: 6 }}>Backend</div>
              <select
                value={data.runner?.backend || 'local_placeholder'}
                onChange={(e) => setData({ ...data, runner: { ...(data.runner || { type: 'codex_cli', command: 'codex' }), backend: e.target.value } })}
                style={{ padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}
              >
                <option value="local_placeholder">local_placeholder</option>
                <option value="openclaw_acp">openclaw_acp</option>
              </select>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '160px 1fr', gap: 8, marginBottom: 8 }}>
              <select
                value={data.runner?.type || 'codex_cli'}
                onChange={(e) => setData({ ...data, runner: { ...(data.runner || { command: 'codex' }), type: e.target.value } })}
                style={{ padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}
              >
                <option value="codex_cli">codex_cli</option>
                <option value="claude_code">claude_code</option>
                <option value="gemini_cli">gemini_cli</option>
                <option value="custom">custom</option>
              </select>
              <input
                value={data.runner?.command || 'codex'}
                onChange={(e) => setData({ ...data, runner: { ...(data.runner || { type: 'codex_cli' }), command: e.target.value } })}
                placeholder="command (e.g., codex)"
                style={{ padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}
              />
            </div>

            {data.runner?.backend === 'openclaw_acp' ? (
              <div style={{ marginBottom: 8 }}>
                <div style={{ color: '#6b7280', fontSize: 12, marginBottom: 6 }}>PM agentId (ACP)</div>
                <input
                  value={data.runner?.agents?.pm || ''}
                  onChange={(e) => setData({ ...data, runner: { ...(data.runner || { type: 'codex_cli', command: 'codex' }), agents: { ...(data.runner?.agents || {}), pm: e.target.value } } })}
                  placeholder="e.g., acp.codex or your configured ACP agent id"
                  style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}
                />
                <div style={{ marginTop: 6, color: '#6b7280', fontSize: 12 }}>
                  Requires OPENCLAW_GATEWAY_URL and OPENCLAW_GATEWAY_TOKEN on the server.
                  Also, Gateway /tools/invoke must allow sessions_spawn.
                </div>
              </div>
            ) : null}

            <div style={{ marginTop: 6, color: '#6b7280', fontSize: 12 }}>
              Note: API keys are not stored here. Use CLI login or environment variables.
            </div>
          </div>

          <button onClick={save} disabled={saving} style={{ padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: '#111827', color: 'white' }}>
            {saving ? 'Saving…' : 'Save'}
          </button>
        </div>
      )}
    </div>
  );
}
