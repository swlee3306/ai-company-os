import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { listApprovals, type ApprovalItem, type Task } from '../lib/api';

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

  return (
    <div style={{ maxWidth: 900 }}>
      <h1 style={{ marginTop: 0 }}>Task</h1>
      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {!task ? (
        <div>Loading…</div>
      ) : (
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white' }}>
          <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280' }}>{task.id}</div>
          <div style={{ fontWeight: 800, fontSize: 18, marginTop: 6 }}>{task.title}</div>
          {task.desc ? <p style={{ marginBottom: 8, color: '#374151' }}>{task.desc}</p> : null}
          <div style={{ color: '#6b7280', fontSize: 12 }}>state: {task.state}</div>

          <div style={{ marginTop: 12 }}>
            <div style={{ fontWeight: 700, marginBottom: 6 }}>Linked approval</div>
            {linked ? (
              <div style={{ fontSize: 12 }}>
                <Link to="/approvals" style={{ color: '#111827' }}>{linked.id}</Link>
                <span style={{ marginLeft: 8, color: '#6b7280' }}>{linked.type} / {linked.risk} / {linked.status}</span>
                <div style={{ marginTop: 6, color: '#6b7280' }}>Open Approval Center and select the item to review evidence/decide.</div>
              </div>
            ) : (
              <div style={{ color: '#6b7280', fontSize: 12 }}>No linked approval found for this task.</div>
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
