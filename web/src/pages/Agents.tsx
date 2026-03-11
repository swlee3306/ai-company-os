import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listAgents, type Agent } from '../lib/api';
import { color, radius, space } from '../ui/tokens';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

export default function Agents() {
  const [data, setData] = useState<Agent[] | null>(null);
  const [error, setError] = useState<string>('');

  const [name, setName] = useState('');
  const [personaRole, setPersonaRole] = useState('BE');
  const [ops, setOps] = useState('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  async function refresh() {
    try {
      setError('');
      setData(await listAgents());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function onCreate() {
    setSaveError('');
    if (!name.trim()) {
      setSaveError('name is required');
      return;
    }
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/api/agents`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name.trim(), persona_role: personaRole, ops_specialty: ops.trim() || undefined }),
      });
      if (!res.ok) throw new Error(await res.text());
      setName('');
      setOps('');
      await refresh();
    } catch (e: any) {
      setSaveError(e?.message ?? String(e));
    } finally {
      setSaving(false);
    }
  }

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Agent Registry</h1>
      <p style={{ color: color.text.muted }}>Registered execution agents with health, scope, concurrency, and heartbeat visibility.</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface, marginBottom: space.card }}>
        <div style={{ fontWeight: 800, marginBottom: 10 }}>Create agent</div>
        <input value={name} onChange={(e) => setName(e.target.value)} placeholder="name" style={{ width: '100%', padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}`, marginBottom: 8 }} />
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
          <select value={personaRole} onChange={(e) => setPersonaRole(e.target.value)} style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }}>
            <option value="PM">PM</option>
            <option value="BE">BE</option>
            <option value="FE">FE</option>
            <option value="QA">QA</option>
            <option value="Reviewer">Reviewer</option>
            <option value="Designer">Designer</option>
          </select>
          <input value={ops} onChange={(e) => setOps(e.target.value)} placeholder="ops_specialty (optional)" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
        </div>
        {saveError ? <div style={{ color: color.text.danger, fontSize: 12, marginBottom: 8 }}>{saveError}</div> : null}
        <button onClick={onCreate} disabled={saving} style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.text.primary, color: 'white' }}>
          {saving ? 'Creating…' : 'Create'}
        </button>
      </div>

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No agents yet.</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            name / persona role / status / scope / version / heartbeat
          </div>
          {data.map((a) => (
            <div key={a.id} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>
                <div>
                  <strong>
                    <Link to={`/agents/${a.id}`} style={{ color: '#111827', textDecoration: 'none' }}>
                      {a.name}
                    </Link>
                  </strong>{' '}
                  {a.ops_specialty ? <span style={{ color: '#6b7280' }}>({a.ops_specialty})</span> : null}
                </div>
                <div style={{ color: '#6b7280', fontSize: 12 }}>
                  {a.persona_role} / {a.status} / {(a.scope || []).join(', ')} / {a.version}
                </div>
              </div>
              <div style={{ textAlign: 'right', fontSize: 12, color: a.heartbeat_seconds > 120 ? '#b91c1c' : a.heartbeat_seconds > 30 ? '#b45309' : '#065f46' }}>
                {a.heartbeat_seconds}s ago
                {a.approval_required ? <div style={{ color: '#b91c1c' }}>approval required</div> : null}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
