export type StatusValue =
  | 'APPROVED'
  | 'ABSENT'
  | 'ENROLLED'
  | 'WAITLISTED';

export interface StatusDef {
  value: StatusValue;
  label: string;
  color: string;       // Tailwind bg class for dot
  textColor: string;   // Tailwind text class for badge text
  badgeBg: string;     // Tailwind bg class for badge background
}

export const STATUSES: StatusDef[] = [
  {
    value: 'APPROVED',
    label: 'Convocado(a)',
    color: 'bg-[#E0A100]',
    textColor: 'text-[#7A4500]',
    badgeBg: 'bg-[#FFF8E1]',
  },
  {
    value: 'ABSENT',
    label: 'Faltoso(a)',
    color: 'bg-[#D64545]',
    textColor: 'text-[#7A1010]',
    badgeBg: 'bg-[#FFF0F0]',
  },
  {
    value: 'ENROLLED',
    label: 'Matriculado(a)',
    color: 'bg-[#2F9E6B]',
    textColor: 'text-[#0D5236]',
    badgeBg: 'bg-[#EDFAF4]',
  },
  {
    value: 'WAITLISTED',
    label: 'Em espera',
    color: 'bg-[#3B82C4]',
    textColor: 'text-[#0D3B6B]',
    badgeBg: 'bg-[#EEF5FF]',
  },
];

export const STATUS_MAP = Object.fromEntries(
  STATUSES.map((s) => [s.value, s])
) as Record<StatusValue, StatusDef>;

export function getStatus(raw: string): StatusDef {
  const key = raw?.toUpperCase() as StatusValue;
  return STATUS_MAP[key] ?? STATUS_MAP['APPROVED'];
}
