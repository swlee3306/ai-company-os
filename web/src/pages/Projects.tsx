import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listProjects, type Project } from '../lib/api';

export default function Projects() {
  const [data, setData] = useState<Project[] | null>(null);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    (async () => {
      try {
        setError('');
        setData(await listProjects());
      } catch (e: any) {
        setError(e?.message ?? String(e));
      }
    })();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <h1 style={{ marginTop: 0 }}>Projects</h1>
      <p style={{ color: '#6b7280' }}>Project overview, participating agents, memory, and evidence bundles.</p>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>{error}</div>
      ) : null}

      {data === null ? (
        <div>Loading…</div>
      ) : data.length === 0 ? (
        <div>
          No projects yet. Run <code>company seed</code> to add demo projects.
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: 12 }}>
          {data.map((p) => (
            <div key={p.id} style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8 }}>
                <div>
                  <div style={{ fontWeight: 700 }}>
                    <Link to={`/projects/${p.id}`} style={{ color: '#111827', textDecoration: 'none' }}>
                      {p.name}
                    </Link>
                  </div>
                  <div style={{ color: '#6b7280', fontSize: 12 }}>{p.id}</div>
                </div>
                <div style={{ fontSize: 12, color: '#6b7280' }}>{p.status}</div>
              </div>

              <div style={{ marginTop: 10, color: '#374151', fontSize: 13 }}>{p.summary}</div>
              <div style={{ marginTop: 10, color: '#6b7280', fontSize: 12 }}>Phase: {p.phase}</div>
              <div style={{ marginTop: 6, color: '#6b7280', fontSize: 12 }}>
                Owner: {p.owner_ceo} · Team Lead: {p.team_lead} {p.due ? `· Due: ${p.due}` : ''}
              </div>

              {p.evidence_bundle?.length ? (
                <div style={{ marginTop: 10, fontSize: 12 }}>
                  <div style={{ color: '#6b7280' }}>Evidence bundle</div>
                  <div style={{ fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace' }}>{p.evidence_bundle.join(', ')}</div>
                </div>
              ) : null}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
