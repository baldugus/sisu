export type StatusValue =
  | 'APPROVED'
  | 'ABSENT'
  | 'ENROLLED'
  | 'WAITLISTED'
  | 'DECLINED';

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
  {
    // Only for promotion offers: the student stays in semester 2.
    value: 'DECLINED',
    label: 'Recusou',
    color: 'bg-[#7C6FB0]',
    textColor: 'text-[#3A2F6B]',
    badgeBg: 'bg-[#F3F0FF]',
  },
];

export const STATUS_MAP = Object.fromEntries(
  STATUSES.map((s) => [s.value, s])
) as Record<StatusValue, StatusDef>;

export function getStatus(raw: string): StatusDef {
  const key = raw?.toUpperCase() as StatusValue;
  return STATUS_MAP[key] ?? STATUS_MAP['APPROVED'];
}

// ── Call entries ────────────────────────────────────────────────────────────
// A call entry's outcome ("pending" | "enrolled" | "absent" | "declined") maps
// onto the same badges as a registration's status.

export type EntryKind = 'initial' | 'waitlist' | 'promotion';

const OUTCOME_TO_STATUS: Record<string, StatusValue> = {
  pending: 'APPROVED',
  enrolled: 'ENROLLED',
  absent: 'ABSENT',
  declined: 'DECLINED',
};

const STATUS_TO_OUTCOME: Record<StatusValue, string | undefined> = {
  APPROVED: 'pending',
  ENROLLED: 'enrolled',
  ABSENT: 'absent',
  DECLINED: 'declined',
  WAITLISTED: undefined,
};

export function outcomeToStatus(outcome: string): StatusValue {
  return OUTCOME_TO_STATUS[outcome] ?? 'APPROVED';
}

export function statusToOutcome(status: string): string | undefined {
  return STATUS_TO_OUTCOME[status as StatusValue];
}

/** Outcomes the operator can set, per entry kind. */
export function statusesForKind(kind?: string): StatusDef[] {
  const allowed: StatusValue[] =
    kind === 'promotion' ? ['APPROVED', 'ENROLLED', 'DECLINED'] : ['APPROVED', 'ENROLLED', 'ABSENT'];
  return allowed.map((v) => statusLabelForKind(v, kind));
}

/** Promotion offers read as accepted/declined rather than enrolled. */
export function statusLabelForKind(value: StatusValue, kind?: string): StatusDef {
  const def = STATUS_MAP[value];
  if (kind === 'promotion' && value === 'ENROLLED') return { ...def, label: 'Aceitou' };
  return def;
}

export const KIND_LABELS: Record<EntryKind, string> = {
  initial: 'Inicial',
  waitlist: 'Lista de espera',
  promotion: 'Promoção',
};

export function semesterLabel(semester?: number | null): string {
  return semester ? `${semester}º` : '—';
}
