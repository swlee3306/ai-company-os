import { useEffect, useMemo, useState } from 'react';
import { listApprovals, type ApprovalItem } from '../lib/api';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

type Decision = 'approve' | 'reject';

type Evidence = {
  approval: any;
  agent?: any;
  project?: any;
  audit_recent?: any[];
};

export default function Approvals() {
  const [data, setData] = useState<ApprovalItem[] | null>(null);
  const [error, setError] = useState<string>('');

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const selected = useMemo(() => (data || []).find((x) => x.id === selectedId) || null, [data, selectedId]);

  const [decision, setDecision] = useState<Decision>('approve');
  const [reason, setReason] = useState<string>('');
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState('');

  const [evidence, setEvidence] = useState<Evidence | null>(null);
  const [evidenceError, setEvidenceError] = useState('');

  async function refresh() {
    try {
      setError('');
      setData(await listApprovals());
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  async function loadEvidence(id: string) {
    setEvidence(null);
    setEvidenceError('');
    try {
      const res = await fetch(`${API_BASE}/api/approvals/${id}/evidence`);
      if (!res.ok) throw new Error(`evidence fetch failed: ${res.status}`);
      setEvidence(await res.json());
    } catch (e: any) {
      setEvidenceError(e?.message ?? String(e));
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
      await loadEvidence(selected.id);
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

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 420px', gap: 16 }}>
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
                    loadEvidence(it.id);
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
          <h2 style={{ marginTop: 0 }}>Review</h2>
          {!selected ? (
            <p style={{ color: '#6b7280' }}>Select an approval item to review evidence and decide.</p>
          ) : (
            <>
              <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace', fontSize: 12, marginBottom: 12 }}>
                <div><strong>type</strong>: {selected.type}</div>
                <div><strong>requester</strong>: {selected.requester}</div>
                <div><strong>target</strong>: {selected.target}</div>
                <div><strong>risk</strong>: {selected.risk}</div>
                {selected.task_id ? <div><strong>task</strong>: {selected.task_id}</div> : null}
              </div>

              <div style={{ marginBottom: 12 }}>
                <div style={{ fontWeight: 700, marginBottom: 6 }}>Evidence</div>
                {evidenceError ? <div style={{ color: '#b91c1c', fontSize: 12 }}>{evidenceError}</div> : null}
                {!evidence ? (
                  <div style={{ color: '#6b7280', fontSize: 12 }}>Loading evidence…</div>
                ) : (
                  <div style={{ fontSize: 12, color: '#374151' }}>
                    {evidence.project ? (
                      <div style={{ marginBottom: 6 }}>
                        <div style={{ color: '#6b7280' }}>Project</div>
                        <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{evidence.project.id} — {evidence.project.name}</div>
                      </div>
                    ) : null}
                    {evidence.agent ? (
                      <div style={{ marginBottom: 6 }}>
                        <div style={{ color: '#6b7280' }}>Agent</div>
                        <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{evidence.agent.id} ({evidence.agent.persona_role}{evidence.agent.ops_specialty ? `/${evidence.agent.ops_specialty}` : ''})</div>
                      </div>
                    ) : null}
                    {evidence.audit_recent?.length ? (
                      <div>
                        <div style={{ color: '#6b7280' }}>Recent audit</div>
                        <div style={{ border: '1px solid #f3f4f6', borderRadius: 10, padding: 10, maxHeight: 160, overflow: 'auto' }}>
                          {evidence.audit_recent.map((a, i) => (
                            <div key={i} style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>
                              {a.ts ?? '-'} / {a.actor ?? '-'} / {a.action ?? '-'}
                            </div>
                          ))}
                        </div>
                      </div>
                    ) : (
                      <div style={{ color: '#6b7280' }}>No recent audit entries found.</div>
                    )}
                  </div>
                )}
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
            </>
          )}
        </aside>
      </div>
    </div>
  );
}
