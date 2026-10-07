import { useCallback, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Loader2, Square } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { RosterTable } from '@/components/RosterTable';
import { useRollCallRows } from '@/hooks/useRollCallRows';
import { CloseCall, FetchCalls } from '@/lib/backend';
import type { types } from '../../wailsjs/go/models';

export default function Chamada() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const callId = Number(id);

  const { rows, loading, error, refresh } = useRollCallRows(callId);
  const [call, setCall] = useState<types.CallSummary | null>(null);
  const [closing, setClosing] = useState(false);

  const loadCall = useCallback(async () => {
    const calls = await FetchCalls() ?? [];
    setCall(calls.find((c) => c.ID === callId) ?? null);
  }, [callId]);

  useEffect(() => { loadCall(); }, [loadCall, rows]);

  const isOpen = call?.Status === 'calling';
  const pending = rows.filter((r) => r.Status === 'APPROVED').length;

  async function handleClose() {
    setClosing(true);
    try {
      await CloseCall(callId);
      toast.success(`${call?.Number}ª chamada fechada.`);
      await loadCall();
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setClosing(false);
    }
  }

  return (
    <div className="flex flex-col h-full">
      <div className="px-6 pt-6 pb-3 border-b-2 border-foreground flex items-center gap-4 shrink-0">
        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => navigate('/chamadas')}>
          <ArrowLeft className="size-4" />
        </Button>
        <div className="flex-1">
          <h1 className="font-heading font-black text-3xl">{call?.Number ?? ''}ª Chamada</h1>
          <p className="text-xs text-muted-foreground mt-0.5">
            {isOpen
              ? pending > 0
                ? `${pending} pendente${pending !== 1 ? 's' : ''} — registre matrícula, falta ou resposta à promoção.`
                : 'Nenhuma pendência — a chamada pode ser fechada.'
              : 'Chamada fechada — somente leitura. Reabra em "Chamadas" para editar.'}
          </p>
          {error && <p className="text-xs text-destructive mt-0.5">{error}</p>}
        </div>
        {isOpen && (
          <Button
            variant="outline"
            size="sm"
            className="gap-1.5"
            disabled={closing || pending > 0}
            title={pending > 0 ? 'Resolva as pendências antes de fechar.' : undefined}
            onClick={handleClose}
          >
            {closing ? <Loader2 className="size-3.5 animate-spin" /> : <Square className="size-3.5" />}
            Fechar chamada
          </Button>
        )}
      </div>

      <div className="flex-1 min-h-0">
        <RosterTable
          rows={rows}
          loading={loading}
          hasSelector
          callId={callId}
          readOnly={!isOpen}
          onRefresh={refresh}
          emptyMessage="Nenhum candidato nesta chamada."
        />
      </div>
    </div>
  );
}
