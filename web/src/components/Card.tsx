export function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section style={{ border: '1px solid #e5e7eb', borderRadius: 12, padding: 16, background: 'white', marginBottom: 12 }}>
      <h2 style={{ marginTop: 0, marginBottom: 10, fontSize: 16 }}>{title}</h2>
      {children}
    </section>
  );
}
