import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

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
      <p style={{ color: '#6b7280' }}>Evidence outputs referenced by projects and approvals.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {data === null ? (
        <div>Loading…</div>
      ) : data.length === 0 ? (
        <div>No artifacts yet. Create one via API: <code>POST /api/artifacts</code>.</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden', background: 'white' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            type / title / project / task / created / uri
          </div>
          {data.map((a) => (
            <div key={a.id} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', fontSize: 12 }}>
              <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>
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
