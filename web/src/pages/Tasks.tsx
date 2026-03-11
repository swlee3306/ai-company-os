import { useEffect, useState } from 'react';
import { createTask, listTasks, transitionTask, type Task } from '../lib/api';

const states = ['draft', 'planned', 'assigned', 'running', 'reviewing', 'done', 'blocked'];

export default function Tasks() {
  const [data, setData] = useState<Task[] | null>(null);
  const [error, setError] = useState('');
  const [title, setTitle] = useState('');
  const [desc, setDesc] = useState('');

  async function refresh() {
    try {
      setError('');
      setData(await listTasks());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function onCreate() {
    if (!title.trim()) return;
    try {
      await createTask(title.trim(), desc);
      setTitle('');
      setDesc('');
      await refresh();
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  async function onTransition(id: string, to: string) {
    try {
      await transitionTask(id, to);
      await refresh();
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Tasks</h1>
      <p style={{ color: '#6b7280' }}>Create, list, and transition tasks (MVP).</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white', marginBottom: 12 }}>
        <div style={{ fontWeight: 700, marginBottom: 10 }}>Create task</div>
        <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Title" style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', marginBottom: 8 }} />
        <textarea value={desc} onChange={(e) => setDesc(e.target.value)} placeholder="Description (optional)" rows={3} style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', resize: 'vertical', marginBottom: 8 }} />
        <button onClick={onCreate} style={{ padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: '#111827', color: 'white' }}>Create</button>
      </div>

      {data === null ? (
        <div>Loading…</div>
      ) : data.length === 0 ? (
        <div>No tasks yet.</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden', background: 'white' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            id / title / state / actions
          </div>
          {data.map((t) => (
            <div key={t.id} style={{ padding: 12, borderBottom: '1px solid #f3f4f6', display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <div>
                <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{t.id}</div>
                <div style={{ fontWeight: 700 }}>{t.title}</div>
                <div style={{ color: '#6b7280', fontSize: 12 }}>state: {t.state}</div>
              </div>
              <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                <select value={t.state} onChange={(e) => onTransition(t.id, e.target.value)} style={{ padding: 8, borderRadius: 10, border: '1px solid #e5e7eb' }}>
                  {states.map((s) => (
                    <option key={s} value={s}>{s}</option>
                  ))}
                </select>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
