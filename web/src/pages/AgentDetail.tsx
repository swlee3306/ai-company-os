import { Link, useParams } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import { type Agent } from '../lib/api';
import { Card } from '../components/Card';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

function kv(label: string, value: any) {
  return (
    <div style={{ display: 'flex', justifyContent: 'space-between', gap: 16, padding: '6px 0', borderBottom: '1px solid #f3f4f6' }}>
      <div style={{ color: '#6b7280' }}>{label}</div>
      <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{String(value ?? '-')}</div>
    </div>
  );
}

export default function AgentDetail() {
  const { id } = useParams();
  const [data, setData] = useState<Agent | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        const res = await fetch(`${API_BASE}/api/agents/${id}`);
        if (!res.ok) throw new Error(`agent fetch failed: ${res.status}`);
        setData(await res.json());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, [id]);

  const danger = useMemo(() => (data?.risk_scope || []).join(', '), [data]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Agent Detail {data ? ` / ${data.name}` : ''}</h1>
      <p style={{ color: '#6b7280' }}>Policy, permissions, and execution history for an infrastructure agent (approval-gated).</p>
      <p>
        <Link to="/agents">← Back to agents</Link>
      </p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {!data ? (
        <div>Loading…</div>
      ) : (
        <>
          <Card title="Configuration and policy">
            {kv('Name', data.name)}
            {kv('Persona role', data.persona_role)}
            {data.ops_specialty ? kv('Ops specialty', data.ops_specialty) : null}
            {kv('Scope', (data.scope || []).join(', '))}
            {kv('Version', data.version)}
            {kv('Heartbeat', `${data.heartbeat_seconds}s ago`)}
            {kv('Approval required', data.approval_required ? 'yes' : 'no')}
          </Card>

          <Card title="Health and recent executions">
            {kv('Status', data.status)}
            <div style={{ marginTop: 10, color: '#6b7280', fontSize: 12 }}>
              Links: <span style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>View logs | View audit</span>
            </div>
            {data.status === 'blocked' || data.approval_required ? (
              <div style={{ marginTop: 10 }}>
                <Link to="/approvals">Open Approval Center →</Link>
              </div>
            ) : null}
          </Card>

          <Card title="Danger permissions">
            <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{danger || '-'}</div>
            <div style={{ marginTop: 10 }}>
              <Link to="/approvals">Action: Open Approval Center →</Link>
            </div>
          </Card>
        </>
      )}
    </div>
  );
}
