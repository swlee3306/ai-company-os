import { useEffect, useMemo, useState } from 'react';
import { listTasks, type Task } from '../lib/api';

const columns = ['draft', 'planned', 'assigned', 'running', 'reviewing', 'done', 'blocked'];

export default function Workflows() {
  const [data, setData] = useState<Task[] | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        setData(await listTasks());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  const byState = useMemo(() => {
    const m = new Map<string, Task[]>();
    for (const c of columns) m.set(c, []);
    for (const t of data || []) {
      m.get(t.state)?.push(t);
    }
    return m;
  }, [data]);

  return (
    <div style={{ maxWidth: 1300 }}>
      <h1 style={{ marginTop: 0 }}>Task / Workflow Board</h1>
      <p style={{ color: '#6b7280' }}>Kanban board prioritized for daily operational control.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {data === null ? (
        <div>Loading…</div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: `repeat(${columns.length}, minmax(180px, 1fr))`, gap: 10 }}>
          {columns.map((c) => (
            <div key={c} style={{ border: '1px solid #e5e7eb', borderRadius: 12, background: 'white', overflow: 'hidden' }}>
              <div style={{ padding: 10, fontWeight: 700, borderBottom: '1px solid #f3f4f6' }}>{c} / {(byState.get(c) || []).length}</div>
              <div style={{ padding: 10, display: 'flex', flexDirection: 'column', gap: 8 }}>
                {(byState.get(c) || []).map((t) => (
                  <div key={t.id} style={{ border: '1px solid #f3f4f6', borderRadius: 10, padding: 10 }}>
                    <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280' }}>{t.id}</div>
                    <div style={{ fontWeight: 700, fontSize: 13 }}>{t.title}</div>
                  </div>
                ))}
                {(byState.get(c) || []).length === 0 ? <div style={{ color: '#9ca3af', fontSize: 12 }}>—</div> : null}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
