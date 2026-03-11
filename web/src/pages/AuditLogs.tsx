import { useEffect, useMemo, useState } from 'react';
import { getAudit } from '../lib/api';
import { color, font, radius, space } from '../ui/tokens';

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
  const [q, setQ] = useState('');

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

  const entries = useMemo(() => {
    const parsed = parseAudit(raw);
    // newest first
    parsed.reverse();
    if (!q.trim()) return parsed;
    const s = q.trim().toLowerCase();
    return parsed.filter((e) => {
      const actor = (e.actor || '').toLowerCase();
      const action = (e.action || '').toLowerCase();
      const fields = JSON.stringify(e.fields || {}).toLowerCase();
      return actor.includes(s) || action.includes(s) || fields.includes(s);
    });
  }, [raw, q]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Audit Logs</h1>
      <p style={{ color: color.text.muted }}>Filter by actor/action (MVP). Detail view comes later.</p>

      <div style={{ marginBottom: 12 }}>
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Search actor/action/fields…"
          style={{ width: '100%', maxWidth: 420, padding: 10, borderRadius: 10, border: `1px solid ${color.border.default}` }}
        />
      </div>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, overflow: 'hidden' }}>
        <div style={{ padding: 12, fontFamily: font.mono, fontSize: 12, color: color.text.muted, borderBottom: `1px solid ${color.border.default}` }}>
          ts / actor / action
        </div>
        {entries.length === 0 ? (
          <div style={{ padding: 12, color: color.text.muted }}>No audit entries found.</div>
        ) : (
          entries.map((e, idx) => (
            <div key={idx} style={{ padding: 12, borderBottom: `1px solid ${color.border.subtle}`, fontFamily: font.mono, fontSize: 12 }}>
              <span style={{ color: color.text.muted }}>{e.ts ?? '-'}</span> / <span>{e.actor ?? '-'}</span> / <span>{e.action ?? '-'}</span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
