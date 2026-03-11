import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listProjects, type Project } from '../lib/api';
import { color, font, radius, space } from '../ui/tokens';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

export default function Projects() {
  const [data, setData] = useState<Project[] | null>(null);
  const [error, setError] = useState<string>('');

  const [name, setName] = useState('');
  const [summary, setSummary] = useState('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  async function refresh() {
    try {
      setError('');
      setData(await listProjects());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function onCreate() {
    setSaveError('');
    if (!name.trim() || !summary.trim()) {
      setSaveError('name and summary are required');
      return;
    }
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/api/projects`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: name.trim(), summary: summary.trim() }),
      });
      if (!res.ok) throw new Error(await res.text());
      setName('');
      setSummary('');
      await refresh();
    } catch (e: any) {
      setSaveError(e?.message ?? String(e));
    } finally {
      setSaving(false);
    }
  }

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Projects</h1>
      <p style={{ color: color.text.muted }}>Project overview, participating agents, memory, and evidence bundles.</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface, marginBottom: space.card }}>
        <div style={{ fontWeight: 800, marginBottom: 10 }}>Create project</div>
        <input value={name} onChange={(e) => setName(e.target.value)} placeholder="name" style={{ width: '100%', padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}`, marginBottom: 8 }} />
        <textarea value={summary} onChange={(e) => setSummary(e.target.value)} placeholder="summary" rows={3} style={{ width: '100%', padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}`, resize: 'vertical', marginBottom: 8 }} />
        {saveError ? <div style={{ color: color.text.danger, fontSize: 12, marginBottom: 8 }}>{saveError}</div> : null}
        <button onClick={onCreate} disabled={saving} style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.text.primary, color: 'white' }}>
          {saving ? 'Creating…' : 'Create'}
        </button>
      </div>

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No projects yet.</div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: space.card }}>
          {data.map((p) => (
            <div key={p.id} style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8 }}>
                <div>
                  <div style={{ fontWeight: 700 }}>
                    <Link to={`/projects/${p.id}`} style={{ color: '#111827', textDecoration: 'none' }}>
                      {p.name}
                    </Link>
                  </div>
                  <div style={{ color: color.text.muted, fontSize: 12, fontFamily: font.mono }}>{p.id}</div>
                </div>
                <div style={{ fontSize: 12, color: '#6b7280' }}>{p.status}</div>
              </div>

              <div style={{ marginTop: 10, color: '#374151', fontSize: 13 }}>{p.summary}</div>
              <div style={{ marginTop: 10, color: '#6b7280', fontSize: 12 }}>Phase: {p.phase}</div>
              <div style={{ marginTop: 6, color: '#6b7280', fontSize: 12 }}>
                Owner: {p.owner_ceo} · Team Lead: {p.team_lead} {p.due ? `· Due: ${p.due}` : ''}
              </div>

              {p.evidence_bundle?.length ? (
                <div style={{ marginTop: 10, fontSize: 12 }}>
                  <div style={{ color: '#6b7280' }}>Evidence bundle</div>
                  <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{p.evidence_bundle.join(', ')}</div>
                </div>
              ) : null}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
