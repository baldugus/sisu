import { useEffect, useState, useCallback, type ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Play, Square, Trash2, ArrowRight, Loader2, Lock, LockOpen, ArrowUp } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import {
  FetchCalls,
  FetchSemesters,
  PreviewCall,
  CreateCall,
  OpenCall,
  CloseCall,
  DeleteCall,
  CloseSemester,
  ReopenSemester,
} from '@/lib/backend';
import type { types } from '../../wailsjs/go/models';

function periodLabel(p?: string) {
  if (p === 'morning') return 'Matutino';
  if (p === 'evening') return 'Noturno';
  return p ?? '';
}

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

function SemesterCard({
  semester,
  closeReason,
  reopenReason,
  onClose,
  onReopen,
}: {
  semester: types.Semester;
  closeReason: string | null;
  reopenReason: string | null;
  onClose: () => void;
  onReopen: () => void;
}) {
  const closed = semester.Status === 'closed';
  const pct = semester.Seats ? (semester.Occupied / semester.Seats) * 100 : 0;

  return (
    <div className={cn('rounded-2xl border p-4 flex flex-col gap-2', closed ? 'bg-muted/50 border-border' : 'bg-card border-border')}>
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <h3 className="font-heading font-bold text-lg">{semester.Number}º semestre</h3>
          <span
            className={cn(
              'inline-flex items-center gap-1 text-xs font-medium px-2 py-0.5 rounded-full',
              closed ? 'bg-muted text-muted-foreground' : 'bg-[#EDFAF4] text-[#0D5236]'
            )}
          >
            {closed ? <Lock className="size-3" /> : <LockOpen className="size-3" />}
            {closed ? `Fechado após a ${semester.ClosedAfterCall}ª chamada` : 'Aberto'}
          </span>
        </div>
        {closed ? (
          <GuardedButton variant="outline" size="sm" className="h-7 text-xs" reason={reopenReason} onClick={onReopen}>
            Reabrir
          </GuardedButton>
        ) : (
          <GuardedButton variant="outline" size="sm" className="h-7 text-xs" reason={closeReason} onClick={onClose}>
            Fechar
          </GuardedButton>
        )}
      </div>
      <div className="flex items-baseline gap-1">
        <span className="font-heading font-black text-2xl tabular-nums">{semester.Occupied}</span>
        <span className="text-sm text-muted-foreground">/ {semester.Seats} vagas ocupadas</span>
      </div>
      <div className="h-1.5 rounded-full bg-muted overflow-hidden">
        <div className="h-full bg-primary" style={{ width: `${Math.min(Math.max(pct, 0), 100)}%` }} />
      </div>
      {closed && (
        <p className="text-xs text-muted-foreground">Não recebe ninguém nas próximas chamadas.</p>
      )}
    </div>
  );
}

function SemesterSummary({ s }: { s: types.CallSemesterSummary }) {
  const parts: ReactNode[] = [];
  if (s.Initial) parts.push(`${s.Initial} da lista de aprovados`);
  if (s.Promotion) parts.push(
    <span key="p" className="inline-flex items-center gap-0.5">
      <ArrowUp className="size-3" /> {s.Promotion} promoç{s.Promotion !== 1 ? 'ões' : 'ão'}
    </span>
  );
  if (s.Waitlist) parts.push(`${s.Waitlist} da lista de espera`);

  return (
    <div className="text-xs">
      <span className="font-semibold text-foreground">{s.Semester}º sem.:</span>{' '}
      <span className="text-muted-foreground">
        {parts.length === 0
          ? 'ninguém'
          : parts.map((p, i) => (
            <span key={i}>{i > 0 && ' · '}{p}</span>
          ))}
      </span>
    </div>
  );
}

function CallCard({
  call,
  isLast,
  reopenReason,
  deleteReason,
  onOpen,
  onClose,
  onDetail,
  onDelete,
}: {
  call: types.CallSummary;
  isLast: boolean;
  reopenReason: string | null;
  deleteReason: string | null;
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
                reason={call.Pending > 0 ? 'Registre matrícula, falta ou resposta de todos antes de fechar.' : null}
                onClick={onClose}
              >
                <Square className="size-3.5" /> Fechar
              </GuardedButton>
            )}
            {!isCalling && isLast && (
              <GuardedButton variant="outline" size="sm" className="h-8 gap-1.5" reason={reopenReason} onClick={onOpen}>
                <Play className="size-3.5" /> Reabrir
              </GuardedButton>
            )}
            {isCalling && isLast && call.Number > 1 && (
              <GuardedButton
                variant="ghost"
                size="sm"
                className="h-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                reason={deleteReason}
                aria-label="Excluir chamada"
                onClick={onDelete}
              >
                <Trash2 className="size-3.5" />
              </GuardedButton>
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

function NamesCell({ regs }: { regs?: types.Registration[] }) {
  if (!regs || regs.length === 0) return <span className="text-muted-foreground">—</span>;
  return (
    <ul className="flex flex-col gap-0.5">
      {regs.map((r) => (
        <li key={r.ID} className="truncate">
          <span className="font-mono text-xs text-muted-foreground mr-1.5">{r.Ranking || '—'}</span>
          {r.Candidate?.Name}
        </li>
      ))}
    </ul>
  );
}

function PreviewDialog({
  open,
  onOpenChange,
  onCreated,
}: {
  open: boolean;
  onOpenChange: (o: boolean) => void;
  onCreated: () => void;
}) {
  const [plan, setPlan] = useState<types.CallPlan | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    if (!open) return;
    setPlan(null);
    setError(null);
    setLoading(true);
    PreviewCall()
      .then(setPlan)
      .catch((e: any) => setError(e?.message ?? 'Ocorreu um erro.'))
      .finally(() => setLoading(false));
  }, [open]);

  const total = plan ? plan.Promoted + plan.Waitlist1 + plan.Waitlist2 : 0;
  const courses = (plan?.Courses ?? []).filter(
    (c) => c.Vacancies1 + c.Vacancies2 > 0 || (c.Promoted?.length ?? 0) > 0
  );
  const closed = (plan?.Semesters ?? []).filter((s) => s.Status === 'closed');

  async function confirm() {
    setCreating(true);
    try {
      await CreateCall();
      toast.success(`${plan?.Number}ª chamada criada.`);
      onOpenChange(false);
      onCreated();
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setCreating(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl max-h-[85vh] flex flex-col gap-4">
        <DialogHeader>
          <DialogTitle>{plan ? `Prévia da ${plan.Number}ª chamada` : 'Nova chamada'}</DialogTitle>
          <DialogDescription>
            Primeiro as vagas do 1º semestre vão para quem pediu para adiantar; o que sobra, nos dois
            semestres, vai para a lista de espera, por classificação em cada turno e cota.
          </DialogDescription>
        </DialogHeader>

        {loading && (
          <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
            <Loader2 className="size-4 animate-spin" /> Calculando…
          </div>
        )}

        {error && <p className="text-sm text-destructive">{error}</p>}

        {plan && (
          <>
            <div className="grid grid-cols-3 gap-3 text-center">
              <div className="rounded-xl bg-accent/50 p-3">
                <p className="font-heading font-black text-2xl tabular-nums">{plan.Promoted}</p>
                <p className="text-xs text-muted-foreground">promoções para o 1º</p>
              </div>
              <div className="rounded-xl bg-accent/50 p-3">
                <p className="font-heading font-black text-2xl tabular-nums">{plan.Waitlist1}</p>
                <p className="text-xs text-muted-foreground">da espera para o 1º</p>
              </div>
              <div className="rounded-xl bg-accent/50 p-3">
                <p className="font-heading font-black text-2xl tabular-nums">{plan.Waitlist2}</p>
                <p className="text-xs text-muted-foreground">da espera para o 2º</p>
              </div>
            </div>

            {closed.length > 0 && (
              <p className="text-xs text-muted-foreground">
                {closed.map((s) => `${s.Number}º`).join(' e ')} semestre fechado — não recebe ninguém.
              </p>
            )}

            <div className="flex-1 overflow-auto -mx-1 px-1">
              {courses.length === 0 ? (
                <p className="text-sm text-muted-foreground py-6 text-center">Não há vagas abertas.</p>
              ) : (
                <table className="w-full text-sm border-collapse">
                  <thead className="sticky top-0 bg-background">
                    <tr className="text-xs uppercase tracking-wide text-muted-foreground text-left">
                      <th className="py-2 pr-3 font-semibold">Turno / cota</th>
                      <th className="py-2 pr-3 font-semibold">Vagas 1º / 2º</th>
                      <th className="py-2 pr-3 font-semibold">Promoção → 1º</th>
                      <th className="py-2 pr-3 font-semibold">Espera → 1º</th>
                      <th className="py-2 font-semibold">Espera → 2º</th>
                    </tr>
                  </thead>
                  <tbody>
                    {courses.map((c) => (
                      <tr key={c.Course?.ID} className="border-t border-border align-top">
                        <td className="py-2 pr-3">
                          <p className="font-medium">{periodLabel(c.Course?.Period)}</p>
                          <p className="text-xs text-muted-foreground line-clamp-2">{c.Course?.Quota}</p>
                        </td>
                        <td className="py-2 pr-3 tabular-nums">{c.Vacancies1} / {c.Vacancies2}</td>
                        <td className="py-2 pr-3 max-w-[160px]"><NamesCell regs={c.Promoted} /></td>
                        <td className="py-2 pr-3 max-w-[160px]"><NamesCell regs={c.Waitlist1} /></td>
                        <td className="py-2 max-w-[160px]"><NamesCell regs={c.Waitlist2} /></td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </>
        )}

        <div className="flex items-center justify-end gap-2 pt-2 border-t border-border">
          {plan && total === 0 && (
            <span className="text-xs text-muted-foreground mr-auto">Ninguém a convocar.</span>
          )}
          <Button variant="ghost" onClick={() => onOpenChange(false)}>Cancelar</Button>
          <Button disabled={!plan || total === 0 || creating} onClick={confirm} className="gap-2">
            {creating && <Loader2 className="size-4 animate-spin" />}
            Criar chamada{plan && total > 0 ? ` (${total})` : ''}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default function Chamadas() {
  const navigate = useNavigate();
  const [calls, setCalls] = useState<types.CallSummary[]>([]);
  const [semesters, setSemesters] = useState<types.Semester[]>([]);
  const [loading, setLoading] = useState(true);
  const [previewOpen, setPreviewOpen] = useState(false);

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
  const bothClosed = semesters.length > 0 && semesters.every((s) => s.Status === 'closed');

  // A semester closed after call N blocks undoing call N (undo runs in reverse order).
  const closedSince = (n: number) =>
    semesters.find((s) => s.ClosedAfterCall != null && s.ClosedAfterCall >= n);

  function closeSemesterReason(): string | null {
    if (calls.length === 0) return 'Importe a lista de aprovados primeiro.';
    if (hasOpenCall) return 'Feche a chamada aberta antes de fechar o semestre.';
    return null;
  }

  function reopenSemesterReason(s: types.Semester): string | null {
    if (s.ClosedAfterCall != null && lastNumber > s.ClosedAfterCall) {
      return `Exclua as chamadas criadas depois do fechamento (a partir da ${s.ClosedAfterCall + 1}ª) antes.`;
    }
    return null;
  }

  function undoReason(call: types.CallSummary): string | null {
    const s = closedSince(call.Number);
    return s ? `O ${s.Number}º semestre foi fechado depois desta chamada. Reabra o semestre antes.` : null;
  }

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
                : bothClosed
                  ? 'Os dois semestres estão fechados.'
                  : null
          }
          onClick={() => setPreviewOpen(true)}
        >
          <Plus className="size-4" />
          Nova chamada
        </GuardedButton>
      </div>

      {semesters.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
          {semesters.map((s) => (
            <SemesterCard
              key={s.Number}
              semester={s}
              closeReason={closeSemesterReason()}
              reopenReason={reopenSemesterReason(s)}
              onClose={() => act(() => CloseSemester(s.Number), `${s.Number}º semestre fechado.`)}
              onReopen={() => act(() => ReopenSemester(s.Number), `${s.Number}º semestre reaberto.`)}
            />
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
              reopenReason={undoReason(call)}
              deleteReason={undoReason(call)}
              onDetail={() => navigate(`/chamadas/${call.ID}`)}
              onOpen={() => act(() => OpenCall(call.ID), `${call.Number}ª chamada reaberta.`)}
              onClose={() => act(() => CloseCall(call.ID), `${call.Number}ª chamada fechada.`)}
              onDelete={() => act(() => DeleteCall(call.ID), `${call.Number}ª chamada removida.`)}
            />
          ))}
        </div>
      )}

      <PreviewDialog open={previewOpen} onOpenChange={setPreviewOpen} onCreated={load} />
    </div>
  );
}
