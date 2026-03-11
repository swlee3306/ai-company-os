import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { listApprovals, type ApprovalItem, type Task } from '../lib/api';
import { color, radius, space } from '../ui/tokens';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

async function getTask(id: string): Promise<Task> {
  const res = await fetch(`${API_BASE}/api/tasks/${id}`);
  if (!res.ok) throw new Error(`task fetch failed: ${res.status}`);
  return res.json();
}

export default function TaskDetail() {
  const { id } = useParams();
  const [task, setTask] = useState<Task | null>(null);
  const [approvals, setApprovals] = useState<ApprovalItem[] | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!id) return;
    (async () => {
      try {
        setError('');
        setTask(await getTask(id));
        setApprovals(await listApprovals());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, [id]);

  const linked = (approvals || []).find((a) => a.task_id === id) || null;

  const [running, setRunning] = useState(false);
  const [runError, setRunError] = useState('');

  async function startRun(pipeline: string) {
    if (!id) return;
    setRunError('');
    setRunning(true);
    try {
      const res = await fetch(`${API_BASE}/api/tasks/${id}/run`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pipeline }),
      });
      if (!res.ok) throw new Error(await res.text());
    } catch (e: any) {
      setRunError(e?.message ?? String(e));
    } finally {
      setRunning(false);
    }
  }

  return (
    <div style={{ maxWidth: 900 }}>
      <h1 style={{ marginTop: 0 }}>Task</h1>
      {error ? (
        <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>{error}</div>
      ) : null}

      {!task ? (
        <div style={{ color: color.text.muted }}>Loading…</div>
      ) : (
        <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface }}>
          <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280' }}>{task.id}</div>
          <div style={{ fontWeight: 800, fontSize: 18, marginTop: 6 }}>{task.title}</div>
          {task.desc ? <p style={{ marginBottom: 8, color: '#374151' }}>{task.desc}</p> : null}
          <div style={{ color: '#6b7280', fontSize: 12 }}>state: {task.state}</div>

          <div style={{ marginTop: 12 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Runner</div>
            {runError ? <div style={{ color: color.text.danger, fontSize: 12, marginBottom: 8 }}>{runError}</div> : null}
            <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
              <button
                onClick={() => startRun('pm_only')}
                disabled={running}
                style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.text.primary, color: 'white' }}
              >
                {running ? 'Running…' : 'Run (PM only)'}
              </button>
              <button
                onClick={() => startRun('full')}
                disabled={running}
                style={{ padding: '10px 14px', borderRadius: 10, border: `1px solid ${color.border.default}`, background: color.bg.surface }}
              >
                Run (full)
              </button>
              <Link to="/runs" style={{ alignSelf: 'center', color: color.text.primary }}>View runs →</Link>
            </div>
          </div>

          <div style={{ marginTop: 12 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Linked approval</div>
            {linked ? (
              <div style={{ fontSize: 12 }}>
                <Link to="/approvals" style={{ color: color.text.primary }}>{linked.id}</Link>
                <span style={{ marginLeft: 8, color: color.text.muted }}>{linked.type} / {linked.risk} / {linked.status}</span>
                <div style={{ marginTop: 6, color: color.text.muted }}>Open Approval Center and select the item to review evidence/decide.</div>
              </div>
            ) : (
              <div style={{ color: color.text.muted, fontSize: 12 }}>No linked approval found for this task.</div>
            )}
          </div>
        </div>
      )}

      <div style={{ marginTop: 12 }}>
        <Link to="/tasks" style={{ color: '#111827' }}>← Back to Tasks</Link>
      </div>
    </div>
  );
}
