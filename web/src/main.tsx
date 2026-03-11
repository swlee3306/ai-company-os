import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import AppShell from './layout/AppShell';
import './index.css';

import Dashboard from './pages/Dashboard';
import Projects from './pages/Projects';
import Agents from './pages/Agents';
import Tasks from './pages/Tasks';
import Workflows from './pages/Workflows';
import Approvals from './pages/Approvals';
import Artifacts from './pages/Artifacts';
import AuditLogs from './pages/AuditLogs';
import Settings from './pages/Settings';
import ArtifactDetail from './pages/ArtifactDetail';

const router = createBrowserRouter([
  {
    element: <AppShell />,
    children: [
      { path: '/', element: <Dashboard /> },
      { path: '/projects', element: <Projects /> },
      { path: '/agents', element: <Agents /> },
      { path: '/tasks', element: <Tasks /> },
      { path: '/workflows', element: <Workflows /> },
      { path: '/approvals', element: <Approvals /> },
      { path: '/artifacts', element: <Artifacts /> },
      { path: '/artifacts/:id', element: <ArtifactDetail /> },
      { path: '/audit-logs', element: <AuditLogs /> },
      { path: '/settings', element: <Settings /> },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
