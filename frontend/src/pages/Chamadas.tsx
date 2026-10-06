import { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Play, Square, Trash2, ArrowRight, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import {
  FetchCalls,
  FetchSemesters,
  CreateCall,
  OpenCall,
  CloseCall,
  DeleteCall,
} from '@/lib/backend';
import type { types } from '../../wailsjs/go/models';

/** A button that explains, on hover, why it is disabled. */
function GuardedButton({
  reason,
  children,
  ...props
}: React.ComponentProps<typeof Button> & { reason?: string | null }) {
  const button = <Button {...props} disabled={props.disabled || !!reason}>{children}</Button>;
  if (!reason) return button;
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span tabIndex={0}>{button}</span>
      </TooltipTrigger>
      <TooltipContent className="max-w-xs">{reason}</TooltipContent>
    </Tooltip>
  );
}

function SemesterCard({ semester }: { semester: types.Semester }) {
  const pct = semester.Seats ? (semester.Occupied / semester.Seats) * 100 : 0;

  return (
    <div className="rounded-2xl border p-4 flex flex-col gap-2 bg-card border-border">
      <h3 className="font-heading font-bold text-lg">{semester.Number}º semestre</h3>
      <div className="flex items-baseline gap-1">
        <span className="font-heading font-black text-2xl tabular-nums">{semester.Occupied}</span>
        <span className="text-sm text-muted-foreground">/ {semester.Seats} vagas ocupadas</span>
      </div>
      <div className="h-1.5 rounded-full bg-muted overflow-hidden">
        <div className="h-full bg-primary" style={{ width: `${Math.min(Math.max(pct, 0), 100)}%` }} />
      </div>
    </div>
  );
}

function SemesterSummary({ s }: { s: types.CallSemesterSummary }) {
  const parts: string[] = [];
  if (s.Initial) parts.push(`${s.Initial} da lista de aprovados`);
  if (s.Waitlist) parts.push(`${s.Waitlist} da lista de espera`);

  return (
    <div className="text-xs">
      <span className="font-semibold text-foreground">{s.Semester}º sem.:</span>{' '}
      <span className="text-muted-foreground">
        {parts.length === 0 ? 'ninguém' : parts.join(' · ')}
      </span>
    </div>
  );
}

function CallCard({
  call,
  isLast,
  onOpen,
  onClose,
  onDetail,
  onDelete,
}: {
  call: types.CallSummary;
  isLast: boolean;
  onOpen: () => void;
  onClose: () => void;
  onDetail: () => void;
  onDelete: () => void;
}) {
  const isCalling = call.Status === 'calling';

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
            <div className="flex items-center gap-2">
              <span
                className={cn(
                  'inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-0.5 rounded-full',
                  isCalling ? 'bg-[#FFF8E1] text-[#7A4500]' : 'bg-muted text-muted-foreground'
                )}
              >
                <span className={cn('w-1.5 h-1.5 rounded-full', isCalling ? 'bg-[#E0A100]' : 'bg-muted-foreground')} />
                {isCalling ? 'Aberta' : 'Fechada'}
              </span>
              {isCalling && call.Pending > 0 && (
                <span className="text-xs text-muted-foreground">
                  {call.Pending} pendente{call.Pending !== 1 ? 's' : ''}
                </span>
              )}
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center gap-2">
            {isCalling && (
              <GuardedButton
                variant="outline"
                size="sm"
                className="h-8 gap-1.5"
                reason={call.Pending > 0 ? 'Registre matrícula ou falta de todos antes de fechar.' : null}
                onClick={onClose}
              >
                <Square className="size-3.5" /> Fechar
              </GuardedButton>
            )}
            {!isCalling && isLast && (
              <Button variant="outline" size="sm" className="h-8 gap-1.5" onClick={onOpen}>
                <Play className="size-3.5" /> Reabrir
              </Button>
            )}
            {isCalling && isLast && call.Number > 1 && (
              <Button
                variant="ghost"
                size="sm"
                className="h-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                aria-label="Excluir chamada"
                onClick={onDelete}
              >
                <Trash2 className="size-3.5" />
              </Button>
            )}
            <Button size="sm" className="h-8 gap-1.5" onClick={onDetail}>
              Gerenciar <ArrowRight className="size-3.5" />
            </Button>
          </div>
        </div>

        <div className="mt-3 flex flex-col gap-0.5">
          {(call.Semesters ?? []).map((s) => <SemesterSummary key={s.Semester} s={s} />)}
        </div>
      </div>
    </div>
  );
}

export default function Chamadas() {
  const navigate = useNavigate();
  const [calls, setCalls] = useState<types.CallSummary[]>([]);
  const [semesters, setSemesters] = useState<types.Semester[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [c, s] = await Promise.all([FetchCalls(), FetchSemesters()]);
      setCalls(c ?? []);
      setSemesters(s ?? []);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  async function act(fn: () => Promise<any>, successMsg: string) {
    try {
      await fn();
      toast.success(successMsg);
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      await load();
    }
  }

  const hasOpenCall = calls.some((c) => c.Status === 'calling');
  const lastNumber = calls.length ? calls[calls.length - 1].Number : 0;

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
        <GuardedButton
          className="gap-2"
          reason={
            calls.length === 0
              ? 'Importe a lista de aprovados primeiro.'
              : hasOpenCall
                ? 'Feche a chamada aberta antes de criar outra.'
                : null
          }
          onClick={() => act(() => CreateCall(), `${lastNumber + 1}ª chamada criada.`)}
        >
          <Plus className="size-4" />
          Nova chamada
        </GuardedButton>
      </div>

      {semesters.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
          {semesters.map((s) => (
            <SemesterCard key={s.Number} semester={s} />
          ))}
        </div>
      )}

      {loading && calls.length === 0 && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-12 justify-center">
          <Loader2 className="size-4 animate-spin" /> Carregando chamadas…
        </div>
      )}

      {!loading && calls.length === 0 && (
        <div className="text-center py-16 text-muted-foreground text-sm">
          Nenhuma chamada. A 1ª chamada é criada ao importar a lista de aprovados.
        </div>
      )}

      {calls.length > 0 && (
        <div>
          {calls.map((call, i) => (
            <CallCard
              key={call.ID}
              call={call}
              isLast={i === calls.length - 1}
              onDetail={() => navigate(`/chamadas/${call.ID}`)}
              onOpen={() => act(() => OpenCall(call.ID), `${call.Number}ª chamada reaberta.`)}
              onClose={() => act(() => CloseCall(call.ID), `${call.Number}ª chamada fechada.`)}
              onDelete={() => act(() => DeleteCall(call.ID), `${call.Number}ª chamada removida.`)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
