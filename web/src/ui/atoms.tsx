import { color, font, radius, space } from './tokens';

export function Card(props: { title?: string; children: any }) {
  return (
    <div style={{ border: `1px solid ${color.border.default}`, borderRadius: radius.card, padding: space.card, background: color.bg.surface }}>
      {props.title ? <div style={{ fontWeight: 800, marginBottom: 10 }}>{props.title}</div> : null}
      {props.children}
    </div>
  );
}

export function MonoHeader(props: { children: any }) {
  return (
    <div style={{ padding: 12, fontFamily: font.mono, fontSize: 12, color: color.text.muted, borderBottom: `1px solid ${color.border.default}` }}>
      {props.children}
    </div>
  );
}

export function ErrorBox(props: { message: string }) {
  return (
    <div style={{ background: color.bg.danger, border: `1px solid ${color.border.danger}`, color: color.text.danger, padding: space.row, borderRadius: radius.card, marginBottom: space.card }}>
      {props.message}
    </div>
  );
}

export function Muted(props: { children: any }) {
  return <span style={{ color: color.text.muted }}>{props.children}</span>;
}

export function PillLink(props: { to: string; children: any }) {
  return (
    <a
      href={props.to}
      style={{
        display: 'inline-block',
        padding: '4px 8px',
        borderRadius: radius.pill,
        border: `1px solid ${color.border.default}`,
        background: color.bg.subtle,
        color: color.text.primary,
        textDecoration: 'none',
      }}
    >
      {props.children}
    </a>
  );
}
