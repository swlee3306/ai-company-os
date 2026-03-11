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

  useEffect(() => {
    (async () => {
      try {
        setError('');
        const res = await fetch(`${API_BASE}/api/artifacts`);
        if (!res.ok) throw new Error(`artifacts fetch failed: ${res.status}`);
        setData(await res.json());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Artifacts</h1>
      <p style={{ color: color.text.muted }}>Evidence outputs referenced by projects and approvals.</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No artifacts yet. Create one via API: <code>POST /api/artifacts</code>.</div>
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
