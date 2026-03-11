import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { createTask, listTasks, transitionTask, type Task } from '../lib/api';
import { color, font, radius, space } from '../ui/tokens';

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
      <p style={{ color: color.text.muted }}>Create, list, and transition tasks (MVP).</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface, marginBottom: 12 }}>
        <div style={{ fontWeight: 700, marginBottom: 10 }}>Create task</div>
        <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Title" style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', marginBottom: 8 }} />
        <textarea value={desc} onChange={(e) => setDesc(e.target.value)} placeholder="Description (optional)" rows={3} style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', resize: 'vertical', marginBottom: 8 }} />
        <button onClick={onCreate} style={{ padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: '#111827', color: 'white' }}>Create</button>
      </div>

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No tasks yet.</div>
      ) : (
        <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, overflow: 'hidden', background: color.bg.surface }}>
          <div style={{ padding: 12, fontFamily: font.mono, fontSize: 12, color: color.text.muted, borderBottom: `1px solid ${color.border.default}` }}>
            id / title / state / actions
          </div>
          {data.map((t) => (
            <div key={t.id} style={{ padding: 12, borderBottom: `1px solid ${color.border.subtle}`, display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <div>
                <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>
                  <Link to={`/tasks/${t.id}`} style={{ color: '#111827', textDecoration: 'none' }}>{t.id}</Link>
                </div>
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
