import { Link, useParams } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import { type Project } from '../lib/api';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787';

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white', marginBottom: 12 }}>
      <h2 style={{ marginTop: 0, marginBottom: 10, fontSize: 16 }}>{title}</h2>
      {children}
    </section>
  );
}

export default function ProjectDetail() {
  const { id } = useParams();
  const [data, setData] = useState<Project | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        const res = await fetch(`${API_BASE}/api/projects/${id}`);
        if (!res.ok) throw new Error(`project fetch failed: ${res.status}`);
        setData(await res.json());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, [id]);

  const evidence = useMemo(() => data?.evidence_bundle || [], [data]);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Project Detail {data ? ` / ${data.id}` : ''}</h1>
      <p style={{ color: '#6b7280' }}>Goal, current work, and evidence for a project.</p>
      <p>
        <Link to="/projects">← Back to projects</Link>
      </p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {!data ? (
        <div>Loading…</div>
      ) : (
        <>
          <Card title="Project overview">
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 12 }}>
              <div><span style={{ color: '#6b7280' }}>Status</span><div>{data.status}</div></div>
              <div><span style={{ color: '#6b7280' }}>Phase</span><div>{data.phase}</div></div>
              <div><span style={{ color: '#6b7280' }}>Owner (CEO)</span><div>{data.owner_ceo}</div></div>
              <div><span style={{ color: '#6b7280' }}>Team Lead</span><div>{data.team_lead}</div></div>
            </div>
            {data.summary ? <div style={{ marginTop: 10, color: '#374151' }}>{data.summary}</div> : null}
          </Card>

          <Card title="Goal / current work / evidence bundle">
            <div style={{ color: '#6b7280', fontSize: 12 }}>Evidence bundle</div>
            <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{evidence.length ? evidence.join(', ') : '-'}</div>
          </Card>

          <Card title="Participating agents and recent work">
            <div style={{ color: '#6b7280', fontSize: 12 }}>Agents</div>
            <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{(data.agents || []).join(', ') || '-'}</div>
          </Card>

          <Card title="Project memory">
            <div style={{ color: '#6b7280' }}>(MVP placeholder)</div>
          </Card>

          <Card title="Key decisions">
            <div style={{ color: '#6b7280' }}>(MVP placeholder)</div>
          </Card>
        </>
      )}
    </div>
  );
}
