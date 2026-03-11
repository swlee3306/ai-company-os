import { Link, useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Artifact = {
  id: string;
  type: string;
  title: string;
  project_id?: string;
  task_id?: string;
  uri: string;
  created_at: string;
  meta?: any;
};

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white', marginBottom: 12 }}>
      <h2 style={{ marginTop: 0, marginBottom: 10, fontSize: 16 }}>{title}</h2>
      {children}
    </section>
  );
}

export default function ArtifactDetail() {
  const { id } = useParams();
  const [data, setData] = useState<Artifact | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        const res = await fetch(`${API_BASE}/api/artifacts/${id}`);
        if (res.status === 404) {
          setData(null);
          return;
        }
        if (!res.ok) throw new Error(`artifact fetch failed: ${res.status}`);
        setData(await res.json());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, [id]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Artifact Detail {id ? ` / ${id}` : ''}</h1>
      <p style={{ color: '#6b7280' }}>Evidence artifact metadata (MVP).</p>
      <p>
        <Link to="/artifacts">← Back to artifacts</Link>
      </p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {!data ? (
        <div>Loading / not found…</div>
      ) : (
        <>
          <Card title="Summary">
            <div><strong>type</strong>: {data.type || 'other'}</div>
            <div><strong>title</strong>: {data.title}</div>
            <div><strong>project</strong>: {data.project_id || '-'}</div>
            <div><strong>task</strong>: {data.task_id || '-'}</div>
            <div><strong>created</strong>: {data.created_at}</div>
          </Card>
          <Card title="URI">
            <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{data.uri}</div>
          </Card>
          <Card title="Meta">
            <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{JSON.stringify(data.meta || {}, null, 2)}</pre>
          </Card>
        </>
      )}
    </div>
  );
}
