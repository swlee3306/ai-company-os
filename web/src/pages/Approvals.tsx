import { useEffect, useMemo, useState } from 'react';
import { listApprovals, type ApprovalItem } from '../lib/api';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Decision = 'approve' | 'reject';

export default function Approvals() {
  const [data, setData] = useState<ApprovalItem[] | null>(null);
  const [error, setError] = useState<string>('');

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const selected = useMemo(() => (data || []).find((x) => x.id === selectedId) || null, [data, selectedId]);

  const [decision, setDecision] = useState<Decision>('approve');
  const [reason, setReason] = useState<string>('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  async function refresh() {
    try {
      setError('');
      setData(await listApprovals());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  async function submitDecision() {
    if (!selected) return;
    setSaveError('');
    if (decision === 'reject' && !reason.trim()) {
      setSaveError('Reason is required for reject.');
      return;
    }
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/api/approvals/${selected.id}/decision`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ decision, reason: reason.trim() }),
      });
      if (!res.ok) {
        const t = await res.text();
        throw new Error(t || `save failed: ${res.status}`);
      }
      await refresh();
    } catch (e: any) {
      setSaveError(e?.message ?? String(e));
    } finally {
      setSaving(false);
    }
  }

  return (
    <div style={{ maxWidth: 1200 }}>
      <h1 style={{ marginTop: 0 }}>Approval Center</h1>
      <p style={{ color: '#6b7280' }}>Queue for production deploys, sensitive permission elevation, and new agent activation requests.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 360px', gap: 16 }}>
        <div style={{ border: '1px solid #e5e7eb', borderRadius: 12, overflow: 'hidden' }}>
          <div style={{ padding: 12, fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, color: '#6b7280', borderBottom: '1px solid #e5e7eb' }}>
            type / requester / target / risk / action
          </div>

          {data === null ? (
            <div style={{ padding: 12 }}>Loading…</div>
          ) : data.length === 0 ? (
            <div style={{ padding: 12 }}>No approvals yet. Run <code>company seed</code> to add demo approvals.</div>
          ) : (
            data.map((it) => {
              const isSelected = it.id === selectedId;
              const color = it.risk === 'HIGH' ? '#b91c1c' : it.risk === 'MEDIUM' ? '#b45309' : '#374151';
              return (
                <button
                  key={it.id}
                  onClick={() => {
                    setSelectedId(it.id);
                    setDecision('approve');
                    setReason('');
                    setSaveError('');
                  }}
                  style={{
                    width: '100%',
                    textAlign: 'left',
                    padding: 12,
                    border: 'none',
                    borderBottom: '1px solid #f3f4f6',
                    background: isSelected ? '#eef2ff' : 'white',
                    cursor: 'pointer',
                    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
                    fontSize: 12,
                    color,
                  }}
                >
                  {it.type} / {it.requester} / {it.target} / {it.risk} / {it.action}
                  {it.status && it.status !== 'pending' ? <span style={{ marginLeft: 8, color: '#111827' }}>({it.status})</span> : null}
                </button>
              );
            })
          )}
        </div>

        <aside style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white' }}>
          <h2 style={{ marginTop: 0 }}>Decision</h2>
          {!selected ? (
            <p style={{ color: '#6b7280' }}>Select an approval item to review evidence and decide.</p>
          ) : (
            <>
              <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, marginBottom: 10 }}>
                <div><strong>type</strong>: {selected.type}</div>
                <div><strong>requester</strong>: {selected.requester}</div>
                <div><strong>target</strong>: {selected.target}</div>
                <div><strong>risk</strong>: {selected.risk}</div>
              </div>

              <div style={{ marginBottom: 10 }}>
                <label style={{ display: 'block', fontSize: 12, color: '#6b7280', marginBottom: 6 }}>Decision</label>
                <select value={decision} onChange={(e) => setDecision(e.target.value as Decision)} style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb' }}>
                  <option value="approve">approve</option>
                  <option value="reject">reject</option>
                </select>
              </div>

              <div style={{ marginBottom: 10 }}>
                <label style={{ display: 'block', fontSize: 12, color: '#6b7280', marginBottom: 6 }}>
                  Reason {decision === 'reject' ? '(required)' : '(optional)'}
                </label>
                <textarea
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                  placeholder={decision === 'reject' ? 'Why is this rejected?' : 'Optional note'}
                  rows={4}
                  style={{ width: '100%', padding: 10, borderRadius: 10, border: '1px solid #e5e7eb', resize: 'vertical' }}
                />
              </div>

              {saveError ? (
                <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 10, borderRadius: 12, marginBottom: 10, fontSize: 12 }}>
                  {saveError}
                </div>
              ) : null}

              <button
                onClick={submitDecision}
                disabled={saving}
                style={{ width: '100%', padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: '#111827', color: 'white' }}
              >
                {saving ? 'Saving…' : 'Submit decision'}
              </button>

              <p style={{ marginTop: 12, color: '#6b7280', fontSize: 12 }}>
                Evidence panel coming next (logs, commits, tests, audit trail).
              </p>
            </>
          )}
        </aside>
      </div>
    </div>
  );
}
