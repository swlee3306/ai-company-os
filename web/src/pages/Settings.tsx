import { useEffect, useState } from 'react';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Settings = {
  driver: { selected: 'k3d' | 'k3s' };
  approval: { policy_text: string };
  runner?: { type: string; command: string; workdir?: string };
};

export default function Settings() {
  const [data, setData] = useState<Settings | null>(null);
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

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
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Approval policy</div>
            <div style={{ color: '#374151' }}>{data.approval.policy_text}</div>
          </div>

          <div style={{ marginBottom: 14 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Runner</div>
            <div style={{ display: 'grid', gridTemplateColumns: '160px 1fr', gap: 8 }}>
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
