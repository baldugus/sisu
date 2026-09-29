import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { RosterTable } from '@/components/RosterTable';
import { useRollCallRows } from '@/hooks/useRollCallRows';

export default function Chamada() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const callId = Number(id);

  const { rows, loading, error, refresh } = useRollCallRows(callId);

  return (
    <div className="flex flex-col h-full">
      <div className="px-6 pt-6 pb-3 border-b-2 border-foreground flex items-center gap-4 shrink-0">
        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => navigate('/chamadas')}>
          <ArrowLeft className="size-4" />
        </Button>
        <div>
          <h1 className="font-heading font-black text-3xl">{id}ª Chamada</h1>
          {error && <p className="text-xs text-destructive mt-0.5">{error}</p>}
        </div>
      </div>

      <div className="flex-1 min-h-0">
        <RosterTable
          rows={rows}
          loading={loading}
          hasSelector
          onRefresh={refresh}
          emptyMessage="Nenhum candidato nesta chamada."
        />
      </div>
    </div>
  );
}
