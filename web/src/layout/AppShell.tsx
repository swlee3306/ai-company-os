import { NavLink, Outlet } from 'react-router-dom';

const navItems = [
  { to: '/', label: 'Dashboard' },
  { to: '/projects', label: 'Projects' },
  { to: '/agents', label: 'Agents' },
  { to: '/tasks', label: 'Tasks' },
  { to: '/workflows', label: 'Workflows' },
  { to: '/approvals', label: 'Approvals' },
  { to: '/artifacts', label: 'Artifacts' },
  { to: '/audit-logs', label: 'Audit Logs' },
  { to: '/runs', label: 'Runs' },
  { to: '/settings', label: 'Settings' },
];

export default function AppShell() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <aside
        style={{
          width: 240,
          padding: 16,
          borderRight: '1px solid #e5e7eb',
          background: '#f9fafb',
        }}
      >
        <div style={{ fontWeight: 700, marginBottom: 16 }}>AI Company OS</div>
        <nav style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {navItems.map((it) => (
            <NavLink
              key={it.to}
              to={it.to}
              end={it.to === '/'}
              style={({ isActive }) => ({
                padding: '10px 12px',
                borderRadius: 10,
                textDecoration: 'none',
                color: isActive ? '#111827' : '#374151',
                background: isActive ? '#eef2ff' : 'transparent',
                border: isActive ? '1px solid #c7d2fe' : '1px solid transparent',
              })}
            >
              {it.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <main style={{ flex: 1, padding: 24 }}>
        <Outlet />
      </main>
    </div>
  );
}
