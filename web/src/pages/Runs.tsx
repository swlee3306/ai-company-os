import { useEffect, useState } from 'react';
import { color, font, radius, space } from '../ui/tokens';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Run = {
  id: string;
  task_id: string;
  runner_type: string;
  pipeline: string;
  status?: string;
};

export default function Runs() {
  const [data, setData] = useState<Run[] | null>(null);
  const [error, setError] = useState('');

  async function refresh() {
    try {
      setError('');
      const res = await fetch(`${API_BASE}/api/runs`);
      if (!res.ok) throw new Error(await res.text());
      setData(await res.json());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Runs</h1>
      <p style={{ color: color.text.muted }}>Execution history (MVP). Shows run folders under the local store.</p>

      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      {data === null ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : data.length === 0 ? (
        <div style={{ color: color.text.muted }}>No runs yet. Start one from a Task.</div>
      ) : (
        <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, overflow: 'hidden', background: color.bg.surface }}>
          <div style={{ padding: 12, fontFamily: font.mono, fontSize: 12, color: color.text.muted, borderBottom: `1px solid ${color.border.default}` }}>
            id / task / runner / pipeline
          </div>
          {data.map((r) => (
            <div key={r.id} style={{ padding: 12, borderBottom: `1px solid ${color.border.subtle}`, fontFamily: font.mono, fontSize: 12 }}>
              {r.id} / {r.task_id} / {r.runner_type} / {r.pipeline}
            </div>
          ))}
        </div>
      )}

      <div style={{ marginTop: 12 }}>
        <button onClick={refresh} style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.bg.surface }}>
          Refresh
        </button>
      </div>
    </div>
  );
}
