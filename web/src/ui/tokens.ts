export const space = {
  page: 32,
  section: 24,
  card: 16,
  row: 12,
} as const;

export const radius = {
  card: 12,
  pill: 999,
} as const;

export const color = {
  text: {
    primary: '#111827',
    muted: '#6b7280',
    danger: '#991b1b',
  },
  border: {
    default: '#e5e7eb',
    subtle: '#f3f4f6',
    danger: '#fecaca',
  },
  bg: {
    surface: 'white',
    subtle: '#f9fafb',
    danger: '#fef2f2',
    selected: '#eef2ff',
  },
} as const;

export const font = {
  mono: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
} as const;
