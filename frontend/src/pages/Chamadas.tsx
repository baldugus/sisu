import { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Play, Square, Trash2, ArrowRight, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import {
  FetchRollCalls,
  FetchSemesters,
  CreateRollCall,
  OpenRollCall,
  CloseRollCall,
  DeleteRollcall,
} from '@/lib/backend';

interface RollCall {
  ID: number;
  Number: number;
  Status: string;
}

function CallCard({
  call,
  isLast,
  onOpen,
  onClose,
  onDetail,
  canDelete,
  onDelete,
}: {
  call: RollCall;
  isLast: boolean;
  onOpen: () => void;
  onClose: () => void;
  onDetail: () => void;
  canDelete: boolean;
  onDelete: () => void;
}) {
  const isCalling = call.Status === 'CALLING';
  const isDone = call.Status === 'DONE';

  return (
    <div className="flex gap-5">
      {/* Spine */}
      <div className="flex flex-col items-center w-12 shrink-0">
        <div
          className={cn(
            'w-12 h-12 rounded-full border-2 flex items-center justify-center font-heading font-black text-xl transition-colors',
            isCalling
              ? 'border-primary bg-primary text-primary-foreground'
              : 'border-border bg-muted text-muted-foreground'
          )}
        >
          {call.Number}
        </div>
        {!isLast && (
          <div className={cn('flex-1 w-0.5 my-2 min-h-[32px]', isCalling ? 'bg-primary/30' : 'bg-border')} />
        )}
      </div>

      {/* Card */}
      <div className={cn(
        'flex-1 rounded-2xl border p-5 mb-5 transition-colors',
        isCalling ? 'border-primary bg-accent/20' : 'border-border bg-card'
      )}>
        <div className="flex items-start justify-between gap-3">
          <div>
            <h3 className="font-heading font-bold text-xl mb-1">
              {call.Number}ª Chamada
            </h3>
            <span
              className={cn(
                'inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-0.5 rounded-full',
                isCalling ? 'bg-[#FFF8E1] text-[#7A4500]' : 'bg-muted text-muted-foreground'
              )}
            >
              <span className={cn('w-1.5 h-1.5 rounded-full', isCalling ? 'bg-[#E0A100]' : 'bg-muted-foreground')} />
              {isCalling ? 'Aberta' : isDone ? 'Fechada' : call.Status}
            </span>
          </div>

          {/* Actions */}
          <div className="flex items-center gap-2">
            {isCalling && (
              <Button variant="outline" size="sm" className="h-8 gap-1.5" onClick={onClose}>
                <Square className="size-3.5" /> Fechar
              </Button>
            )}
            {isDone && (
              <Button variant="outline" size="sm" className="h-8 gap-1.5" onClick={onOpen}>
                <Play className="size-3.5" /> Reabrir
              </Button>
            )}
            {canDelete && isLast && (
              <Button variant="ghost" size="sm" className="h-8 text-destructive hover:text-destructive hover:bg-destructive/10" onClick={onDelete}>
                <Trash2 className="size-3.5" />
              </Button>
            )}
            <Button size="sm" className="h-8 gap-1.5" onClick={onDetail}>
              Gerenciar <ArrowRight className="size-3.5" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function Chamadas() {
  const navigate = useNavigate();
  const [calls, setCalls] = useState<RollCall[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const raw = await FetchRollCalls() ?? [];
      setCalls(raw.map((c) => ({
        ID: c.ID,
        Number: c.Number,
        Status: c.Status?.toUpperCase?.() ?? 'DONE',
      })));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  async function act(fn: (...args: any[]) => Promise<any>, successMsg: string) {
    setBusy(true);
    try {
      await fn();
      await load();
      toast.success(successMsg);
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setBusy(false);
    }
  }

  async function handleCreateRollCall() {
    setBusy(true);
    try {
      const semesters = await FetchSemesters() ?? [];
      // SemesterStatus serializes via MarshalText → runtime value is the string "open"
      const openSemester = semesters.find((s) => (s.Status as unknown as string) === 'open');
      if (!openSemester) {
        toast.error('Nenhum semestre aberto encontrado.');
        return;
      }
      await CreateRollCall(openSemester.ID);
      await load();
      toast.success('Nova chamada criada.');
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setBusy(false);
    }
  }

  const hasOpenCall = calls.some((c) => c.Status === 'CALLING');
  const lastCall = calls[calls.length - 1];

  return (
    <div className="px-6 py-6 h-full overflow-auto">
      {/* Header */}
      <div className="flex items-start justify-between mb-6 border-b-2 border-foreground pb-4">
        <div>
          <h1 className="font-heading font-black text-3xl">Chamadas</h1>
          <p className="text-muted-foreground text-sm mt-0.5">
            {calls.length} chamada{calls.length !== 1 ? 's' : ''} no ciclo atual
          </p>
        </div>
        <Button
          className="gap-2"
          disabled={busy || hasOpenCall}
          onClick={handleCreateRollCall}
        >
          {busy ? <Loader2 className="size-4 animate-spin" /> : <Plus className="size-4" />}
          Nova chamada
        </Button>
      </div>

      {loading && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-12 justify-center">
          <Loader2 className="size-4 animate-spin" /> Carregando chamadas…
        </div>
      )}

      {!loading && calls.length === 0 && (
        <div className="text-center py-16 text-muted-foreground text-sm">
          Nenhuma chamada criada. Clique em "Nova chamada" para começar.
        </div>
      )}

      {!loading && calls.length > 0 && (
        <div>
          {calls.map((call, i) => (
            <CallCard
              key={call.ID}
              call={call}
              isLast={i === calls.length - 1}
              onDetail={() => navigate(`/chamadas/${call.ID}`)}
              onOpen={() => act(() => OpenRollCall(call.ID), `${call.Number}ª chamada reaberta.`)}
              onClose={() => act(() => CloseRollCall(call.ID), `${call.Number}ª chamada fechada.`)}
              canDelete={calls.length > 1}
              onDelete={() => act(() => DeleteRollcall(call.ID), `${call.Number}ª chamada removida.`)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
