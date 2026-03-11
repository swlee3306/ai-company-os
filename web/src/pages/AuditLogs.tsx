import { useEffect, useMemo, useState } from 'react';
import { getAudit } from '../lib/api';

type AuditEntry = {
  ts?: string;
  actor?: string;
  action?: string;
  fields?: any;
};

function parseAudit(text: string): AuditEntry[] {
  const lines = text
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean);
  const out: AuditEntry[] = [];
  for (const l of lines) {
    try {
      out.push(JSON.parse(l));
    } catch {
      out.push({ action: l });
    }
  }
  return out;
}

export default function AuditLogs() {
  const [raw, setRaw] = useState<string>('');
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        setRaw(await getAudit());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  const entries = useMemo(() => parseAudit(raw), [raw]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Audit Logs</h1>
      <p style={{ color: '#6b7280' }}>Filters coming next. Detail view will expose request id, task id, and agent id context.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden' }}>
        <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
          ts / actor / action
        </div>
        {entries.length === 0 ? (
          <div style={{ padding: 12 }}>No audit entries yet. Try: <code>company status</code> or <code>company serve</code>.</div>
        ) : (
          entries.map((e, idx) => (
            <div key={idx} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12 }}>
              <span style={{ color: '#6b7280' }}>{e.ts ?? '-'}</span> / <span>{e.actor ?? '-'}</span> / <span>{e.action ?? '-'}</span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
