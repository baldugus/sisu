import { useState, useEffect, useCallback } from 'react';
import type { RowData } from '@/components/RosterTable';
import {
  FetchApplicationsByRollCall,
  FetchRegistration,
} from '@/lib/backend';

function periodLabel(p: string) {
  if (p === 'morning') return 'Matutino';
  if (p === 'evening') return 'Noturno';
  return p;
}

export function useRollCallRows(callId: number) {
  const [rows, setRows] = useState<RowData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetch = useCallback(async () => {
    if (!callId) return;
    setLoading(true);
    setError(null);
    try {
      const apps = await FetchApplicationsByRollCall(callId) ?? [];
      const details = await Promise.all(apps.map((a) => FetchRegistration(a.ID)));

      const mapped: (RowData | null)[] = details.map((d) => {
        const reg = d?.Registration;
        const course = d?.Course;
        const candidate = reg?.Candidate;
        if (!reg) return null;
        return {
          ID: reg.ID,
          Name: candidate?.Name ?? '',
          CPF: candidate?.CPF ?? '',
          Period: periodLabel(course?.Period ?? ''),
          Quota: course?.Quota ?? '',
          Status: reg.Status?.toUpperCase() ?? 'APPROVED',
          EnrollmentID: reg.EnrollmentID,
          Ranking: reg.Ranking,
        };
      });
      const data: RowData[] = mapped.filter((r): r is RowData => r !== null);

      setRows(data);
    } catch (e: any) {
      setError(e?.message ?? 'Erro ao carregar chamada.');
      setRows([]);
    } finally {
      setLoading(false);
    }
  }, [callId]);

  useEffect(() => { fetch(); }, [fetch]);

  return { rows, loading, error, refresh: fetch };
}
