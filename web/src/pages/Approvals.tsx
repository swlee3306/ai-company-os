import { useEffect, useState } from 'react';
import { listApprovals, type ApprovalItem } from '../lib/api';

export default function Approvals() {
  const [data, setData] = useState<ApprovalItem[] | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        setData(await listApprovals());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Approval Center</h1>
      <p style={{ color: '#6b7280' }}>Queue for production deploys, sensitive permission elevation, and new agent activation requests.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {data === null ? (
        <div>Loading…</div>
      ) : data.length === 0 ? (
        <div>No approvals yet. Run <code>company seed</code> to add demo approvals.</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            type / requester / target / risk / action
          </div>
          {data.map((it) => (
            <div key={it.id} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12 }}>
              <span style={{ color: it.risk === 'HIGH' ? '#b91c1c' : it.risk === 'MEDIUM' ? '#b45309' : '#374151' }}>
                {it.type} / {it.requester} / {it.target} / {it.risk} / {it.action}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
