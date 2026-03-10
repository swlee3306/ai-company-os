export default function PageScaffold({ title }: { title: string }) {
  return (
    <div>
      <h1 style={{ marginTop: 0 }}>{title}</h1>
      <p style={{ color: '#6b7280' }}>Empty / Loading / Error states will be wired next.</p>
      <div style={{ border: '1px dashed #d1d5db', padding: 16, borderRadius: 12 }}>
        <em>Placeholder</em>
      </div>
    </div>
  );
}
