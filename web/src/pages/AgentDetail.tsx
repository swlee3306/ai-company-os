import { Link, useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { type Agent } from '../lib/api';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

export default function AgentDetail() {
  const { id } = useParams();
  const [data, setData] = useState<Agent | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        const res = await fetch(`${API_BASE}/api/agents/${id}`);
        if (!res.ok) throw new Error(`agent fetch failed: ${res.status}`);
        setData(await res.json());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, [id]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Agent Detail</h1>
      <p style={{ color: '#6b7280' }}>ID: {id}</p>
      <p>
        <Link to="/agents">← Back to agents</Link>
      </p>
      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}
      <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{data ? JSON.stringify(data, null, 2) : 'loading…'}</pre>
    </div>
  );
}
