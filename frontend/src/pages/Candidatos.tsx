import { useState } from 'react';
import { RosterTable } from '@/components/RosterTable';
import { useSelectionRows } from '@/hooks/useSelectionRows';

export default function Candidatos() {
  const [kind, setKind] = useState<'approved' | 'waitlisted'>('approved');
  const { rows, loading, error } = useSelectionRows(kind);

  return (
    <div className="flex flex-col h-full">
      <div className="px-6 pt-6 pb-3 border-b-2 border-foreground shrink-0">
        <h1 className="font-heading font-black text-3xl">Candidatos</h1>
        {error && (
          <p className="text-xs text-destructive mt-1">{error}</p>
        )}
      </div>

      <div className="flex-1 min-h-0">
        <RosterTable
          rows={rows}
          loading={loading}
          hasSelector={false}
          showKindFilter
          kind={kind}
          onKindChange={setKind}
          showContact
          emptyMessage="Nenhum candidato nesta seleção."
        />
      </div>
    </div>
  );
}
