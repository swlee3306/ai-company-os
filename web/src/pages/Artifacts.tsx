import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { color, font, radius, space } from '../ui/tokens';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Artifact = {
  id: string;
  type: string;
  title: string;
  project_id?: string;
  task_id?: string;
  uri: string;
  created_at: string;
};

export default function Artifacts() {
  const [data, setData] = useState<Artifact[] | null>(null);
  const [error, setError] = useState('');

  const [title, setTitle] = useState('');
  const [uri, setURI] = useState('');
  const [type, setType] = useState('runbook');
  const [projectId, setProjectId] = useState('');
  const [taskId, setTaskId] = useState('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  async function refresh() {
    try {
      setError('');
      const res = await fetch(`${API_BASE}/api/artifacts`);
      if (!res.ok) throw new Error(`artifacts fetch failed: ${res.status}`);
      setData(await res.json());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function onCreate() {
    setSaveError('');
    if (!title.trim() || !uri.trim()) {
      setSaveError('title and uri are required');
      return;
    }
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/api/artifacts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          type: type.trim() || 'other',
          title: title.trim(),
          uri: uri.trim(),
          project_id: projectId.trim() || undefined,
          task_id: taskId.trim() || undefined,
        }),
      });
      if (!res.ok) throw new Error(await res.text());
      setTitle('');
      setURI('');
      setProjectId('');
      setTaskId('');
      await refresh();
    } catch (e: any) {
      setSaveError(e?.message ?? String(e));
    } finally {
      setSaving(false);
    }
  }

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Artifacts</h1>
      <p style={{ color: color.text.muted }}>Evidence outputs referenced by projects and approvals.</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface, marginBottom: space.card }}>
        <div style={{ fontWeight: 800, marginBottom: 10 }}>Create artifact</div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8, marginBottom: 8 }}>
          <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="title" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
          <input value={uri} onChange={(e) => setURI(e.target.value)} placeholder="uri" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '160px 1fr 1fr', gap: 8, marginBottom: 8 }}>
          <input value={type} onChange={(e) => setType(e.target.value)} placeholder="type" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
          <input value={projectId} onChange={(e) => setProjectId(e.target.value)} placeholder="project_id (optional)" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
          <input value={taskId} onChange={(e) => setTaskId(e.target.value)} placeholder="task_id (optional)" style={{ padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }} />
        </div>
        {saveError ? <div style={{ color: color.text.danger, fontSize: 12, marginBottom: 8 }}>{saveError}</div> : null}
        <button onClick={onCreate} disabled={saving} style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.text.primary, color: 'white' }}>
          {saving ? 'Creating…' : 'Create'}
        </button>
      </div>

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No artifacts yet.</div>
      ) : (
        <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, overflow: 'hidden', background: color.bg.surface }}>
          <div style={{ padding: 12, fontFamily: font.mono, fontSize: 12, color: color.text.muted, borderBottom: `1px solid ${color.border.default}` }}>
            type / title / project / task / created / uri
          </div>
          {data.map((a) => (
            <div key={a.id} style={{ padding: 12, borderBottom: `1px solid ${color.border.subtle}`, fontSize: 12 }}>
              <div style={{ fontFamily: font.mono }}>
                {a.type || 'other'} /{' '}
                <Link to={`/artifacts/${a.id}`} style={{ color: '#111827' }}>{a.title}</Link>
                {' '} / {a.project_id || '-'} / {a.task_id || '-'} / {a.created_at} / {a.uri}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
