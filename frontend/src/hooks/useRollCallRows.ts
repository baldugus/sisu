import { useState, useEffect, useCallback } from 'react';
import type { RowData } from '@/components/RosterTable';
import { FetchCallEntries } from '@/lib/backend';
import { outcomeToStatus, type EntryKind } from '@/lib/status';

function periodLabel(p: string) {
  if (p === 'morning') return 'Matutino';
  if (p === 'evening') return 'Noturno';
  return p;
}

/** Rows of a call: one per call entry, with the entry's semester, kind and outcome. */
export function useRollCallRows(callId: number) {
  const [rows, setRows] = useState<RowData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetch = useCallback(async () => {
    if (!callId) return;
    setLoading(true);
    setError(null);
    try {
      const entries = await FetchCallEntries(callId) ?? [];

      const data: RowData[] = entries.flatMap((d) => {
        const reg = d.Registration;
        const entry = d.Entry;
        if (!reg || !entry) return [];
        const candidate = reg.Candidate;
        return [{
          ID: reg.ID,
          Name: candidate?.Name ?? '',
          CPF: candidate?.CPF ?? '',
          Email: candidate?.Email ?? '',
          Period: periodLabel(d.Course?.Period ?? ''),
          Quota: d.Course?.Quota ?? '',
          Status: outcomeToStatus(entry.Outcome),
          EnrollmentID: reg.EnrollmentID,
          Ranking: reg.Ranking,
          Semester: entry.Semester,
          Kind: entry.Kind as EntryKind,
        }];
      });

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
