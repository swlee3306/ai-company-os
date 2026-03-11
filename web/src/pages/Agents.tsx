import { useEffect, useState } from 'react';
import { listAgents, type Agent } from '../lib/api';

export default function Agents() {
  const [data, setData] = useState<Agent[] | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        setData(await listAgents());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Agent Registry</h1>
      <p style={{ color: '#6b7280' }}>Registered execution agents with health, scope, concurrency, and heartbeat visibility.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {data === null ? (
        <div>Loading…</div>
      ) : data.length === 0 ? (
        <div>No agents yet. Run <code>company seed</code> to add demo agents.</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            name / persona role / status / scope / version / heartbeat
          </div>
          {data.map((a) => (
            <div key={a.id} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>
                <div>
                  <strong>{a.name}</strong> {a.ops_specialty ? <span style={{ color: '#6b7280' }}>({a.ops_specialty})</span> : null}
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
