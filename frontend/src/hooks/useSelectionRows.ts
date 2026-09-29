import { useState, useEffect, useCallback } from 'react';
import type { RowData } from '@/components/RosterTable';
import {
  FetchApprovedSelection,
  FetchInterestedSelection,
  FetchRegistrationsBySelectionID,
  FetchRegistration,
} from '@/lib/backend';

function periodLabel(p: string) {
  if (p === 'morning') return 'Matutino';
  if (p === 'evening') return 'Noturno';
  return p;
}

async function loadRowsForSelection(
  getSelection: () => Promise<any>,
  defaultStatus: string
): Promise<RowData[]> {
  const sel = await getSelection();
  if (!sel) return [];

  const regs = await FetchRegistrationsBySelectionID(sel.ID) ?? [];

  const details = await Promise.all(regs.map((r) => FetchRegistration(r.ID)));

  const mapped: (RowData | null)[] = details.map((d) => {
    const reg = d?.Registration;
    const course = d?.Course;
    const candidate = reg?.Candidate;
    if (!reg) return null;
    return {
      ID: reg.ID,
      Name: candidate?.Name ?? '',
      CPF: candidate?.CPF ?? '',
      Email: candidate?.Email ?? '',
      Period: periodLabel(course?.Period ?? ''),
      Quota: course?.Quota ?? '',
      Status: reg.Status?.toUpperCase() ?? defaultStatus,
      EnrollmentID: reg.EnrollmentID,
      Ranking: reg.Ranking,
    };
  });
  return mapped.filter((r): r is RowData => r !== null);
}

export function useSelectionRows(kind: 'approved' | 'waitlisted') {
  const [rows, setRows] = useState<RowData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetch = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const getter = kind === 'approved' ? FetchApprovedSelection : FetchInterestedSelection;
      const defaultStatus = kind === 'approved' ? 'APPROVED' : 'WAITLISTED';
      const data = await loadRowsForSelection(getter, defaultStatus);
      setRows(data);
    } catch (e: any) {
      setError(e?.message ?? 'Erro ao carregar dados.');
      setRows([]);
    } finally {
      setLoading(false);
    }
  }, [kind]);

  useEffect(() => { fetch(); }, [fetch]);

  return { rows, loading, error, refresh: fetch };
}
