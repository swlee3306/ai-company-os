import { useEffect, useState } from 'react';
import { getAudit, getDoctor, getStatus, type Status } from '../lib/api';

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, marginBottom: 16 }}>
      <h2 style={{ margin: '0 0 12px 0' }}>{title}</h2>
      {children}
    </section>
  );
}

export default function Dashboard() {
  const [status, setStatus] = useState<Status | null>(null);
  const [doctor, setDoctor] = useState<any>(null);
  const [audit, setAudit] = useState<string>('');
  const [error, setError] = useState<string>('');

  async function refresh() {
    setError('');
    try {
      const [s, d, a] = await Promise.all([getStatus(), getDoctor(), getAudit()]);
      setStatus(s);
      setDoctor(d);
      setAudit(a);
    } catch (e: any) {
      setError(e?.message ?? String(e));
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <div style={{ maxWidth: 1100 }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <div>
          <h1 style={{ margin: 0 }}>Operational Dashboard</h1>
          <p style={{ margin: '6px 0 0 0', color: '#6b7280' }}>
            Monitor project load, agent health, approvals, and execution cost across the company.
          </p>
        </div>
        <button onClick={refresh} style={{ padding: '10px 14px', borderRadius: 10, border: '1px solid #e5e7eb', background: 'white' }}>
          Refresh
        </button>
      </header>

      {error ? (
        <div style={{ background: '#fef2f2', border: '1px solid #fecaca', color: '#991b1b', padding: 12, borderRadius: 12, marginBottom: 16 }}>
          {error}
        </div>
      ) : null}

      <Section title="Status">
        <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{status ? JSON.stringify(status, null, 2) : 'loading...'}</pre>
      </Section>

      <Section title="Doctor">
        <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{doctor ? JSON.stringify(doctor, null, 2) : 'loading...'}</pre>
      </Section>

      <Section title="Audit log">
        <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{audit || 'no audit entries yet'}</pre>
      </Section>

      <footer style={{ marginTop: 20, color: '#6b7280', fontSize: 12 }}>
        API: {import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8787'}
      </footer>
    </div>
  );
}
