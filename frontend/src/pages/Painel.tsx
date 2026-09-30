import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowRight, Upload, Users, ListOrdered, CheckCircle2, Clock } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import {
  FetchApprovedSelection,
  FetchInterestedSelection,
  FetchRegistrationsBySelectionID,
  FetchRollCalls,
} from '@/lib/backend';

interface SelectionInfo {
  year: number;
  totalApproved: number;
  totalWaitlisted: number;
  totalEnrolled: number;
}

interface RollCall {
  ID: number;
  Number: number;
  Status: string;
}

function StatCard({
  label,
  value,
  sub,
  accent,
  progress,
}: {
  label: string;
  value: string | number;
  sub?: string;
  accent?: boolean;
  progress?: number;
}) {
  return (
    <div
      className={cn(
        'rounded-2xl border border-border p-5 flex flex-col gap-1',
        accent ? 'bg-primary text-primary-foreground border-primary' : 'bg-card'
      )}
    >
      <span className={cn('text-xs font-semibold uppercase tracking-widest', accent ? 'text-primary-foreground/70' : 'text-muted-foreground')}>
        {label}
      </span>
      <span className={cn('font-heading font-black text-4xl tabular-nums', accent && 'text-primary-foreground')}>
        {value}
      </span>
      {sub && (
        <span className={cn('text-xs', accent ? 'text-primary-foreground/70' : 'text-muted-foreground')}>
          {sub}
        </span>
      )}
      {progress != null && (
        <div className="h-1.5 rounded-full bg-muted overflow-hidden mt-1">
          <div
            className="h-full bg-primary"
            style={{ width: `${Math.min(Math.max(progress, 0), 100)}%` }}
          />
        </div>
      )}
    </div>
  );
}

function CallNode({ call, isLast, active }: { call: RollCall; isLast: boolean; active: boolean }) {
  const navigate = useNavigate();
  const isCalling = call.Status === 'CALLING';
  const isDone = call.Status === 'DONE';

  return (
    <div className="flex gap-4">
      {/* Spine */}
      <div className="flex flex-col items-center w-10 shrink-0">
        <button
          onClick={() => navigate(`/chamadas/${call.ID}`)}
          className={cn(
            'w-10 h-10 rounded-full border-2 flex items-center justify-center font-heading font-black text-base transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
            isCalling
              ? 'border-primary bg-primary text-primary-foreground'
              : isDone
              ? 'border-border bg-muted text-muted-foreground'
              : 'border-border bg-background text-foreground'
          )}
        >
          {call.Number}
        </button>
        {!isLast && (
          <div className={cn('flex-1 w-0.5 my-1 min-h-[20px]', isCalling ? 'bg-primary/30' : 'bg-border')} />
        )}
      </div>

      {/* Label */}
      <div className="pb-5 flex items-start gap-2 pt-2">
        <span className="text-sm font-semibold">{call.Number}ª Chamada</span>
        <span
          className={cn(
            'text-xs px-2 py-0.5 rounded-full font-medium',
            isCalling ? 'bg-[#FFF8E1] text-[#7A4500]' : 'bg-muted text-muted-foreground'
          )}
        >
          {isCalling ? 'Aberta' : isDone ? 'Fechada' : call.Status}
        </span>
      </div>
    </div>
  );
}

export default function Painel() {
  const navigate = useNavigate();
  const [info, setInfo] = useState<SelectionInfo | null>(null);
  const [calls, setCalls] = useState<RollCall[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const [approvedRes, waitlistRes, callsRes] = await Promise.all([
          FetchApprovedSelection(),
          FetchInterestedSelection(),
          FetchRollCalls(),
        ]);

        const approved = approvedRes;
        const waitlisted = waitlistRes;

        if (!approved) {
          setInfo(null);
        } else {
          const [approvedRegs, waitlistRegs] = await Promise.all([
            FetchRegistrationsBySelectionID(approved.ID),
            waitlisted ? FetchRegistrationsBySelectionID(waitlisted.ID) : Promise.resolve([]),
          ]);
          const enrolled = (approvedRegs ?? []).filter(
            (r) => r.Status?.toUpperCase?.() === 'ENROLLED'
          ).length;
          setInfo({
            year: approved.Year,
            totalApproved: approvedRegs?.length ?? 0,
            totalWaitlisted: waitlistRegs?.length ?? 0,
            totalEnrolled: enrolled,
          });
        }

        const rawCalls = callsRes ?? [];
        setCalls(
          rawCalls.map((c) => ({
            ID: c.ID,
            Number: c.Number,
            Status: c.Status?.toUpperCase?.() ?? 'DONE',
          }))
        );
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const activeCall = calls.find((c) => c.Status === 'CALLING');

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
        Carregando painel…
      </div>
    );
  }

  if (!info) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-6 px-6 text-center">
        <div className="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
          <Upload className="size-8 text-primary" />
        </div>
        <div>
          <h1 className="font-heading font-black text-3xl mb-2">Iniciar ciclo SISU</h1>
          <p className="text-muted-foreground text-sm max-w-xs">
            Importe o arquivo CSV dos convocados para começar a gerenciar o processo seletivo.
          </p>
        </div>
        <Button onClick={() => navigate('/dados')} className="gap-2">
          Importar dados
          <ArrowRight className="size-4" />
        </Button>
      </div>
    );
  }

  return (
    <div className="px-6 py-6 h-full overflow-auto">
      {/* Header */}
      <div className="mb-6 border-b-2 border-foreground pb-4">
        <h1 className="font-heading font-black text-3xl">Painel</h1>
        <p className="text-muted-foreground text-sm mt-0.5">
          SISU{' '}
          <span className="font-mono font-bold text-foreground">{info.year}</span>
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard label="Convocados" value={info.totalApproved} sub="aprovados no SISU" accent />
        <StatCard label="Em espera" value={info.totalWaitlisted} sub="candidatos em espera" />
        <StatCard
          label="Chamadas"
          value={calls.length}
          sub={activeCall ? `${activeCall.Number}ª chamada aberta` : 'nenhuma chamada aberta'}
        />
        <StatCard
          label="Vagas preenchidas"
          value={`${info.totalEnrolled} / ${info.totalApproved}`}
          sub={
            info.totalApproved
              ? `${((info.totalEnrolled / info.totalApproved) * 100).toFixed(1)}% das vagas`
              : 'sem vagas'
          }
          progress={info.totalApproved ? (info.totalEnrolled / info.totalApproved) * 100 : 0}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Spine */}
        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-heading font-bold text-lg">Chamadas</h2>
            <Button variant="outline" size="sm" className="gap-1.5" onClick={() => navigate('/chamadas')}>
              <ListOrdered className="size-3.5" />
              Gerenciar
            </Button>
          </div>

          {calls.length === 0 ? (
            <div className="text-sm text-muted-foreground py-6 flex items-center gap-2">
              <Clock className="size-4" />
              Nenhuma chamada criada ainda.
            </div>
          ) : (
            <div>
              {calls.map((call, i) => (
                <CallNode
                  key={call.ID}
                  call={call}
                  isLast={i === calls.length - 1}
                  active={call.Status === 'CALLING'}
                />
              ))}
            </div>
          )}
        </div>

        {/* Quick actions */}
        <div>
          <h2 className="font-heading font-bold text-lg mb-4">Próxima ação</h2>
          {activeCall ? (
            <button
              onClick={() => navigate(`/chamadas/${activeCall.ID}`)}
              className="w-full text-left rounded-2xl border-2 border-primary p-5 hover:bg-accent/40 transition-colors group focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <div className="flex items-center gap-3 mb-2">
                <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-heading font-black text-lg">
                  {activeCall.Number}
                </div>
                <div>
                  <p className="font-semibold text-sm">{activeCall.Number}ª Chamada — aberta</p>
                  <p className="text-xs text-muted-foreground">Clique para gerenciar matrículas</p>
                </div>
              </div>
              <div className="flex items-center gap-1 text-primary text-xs font-semibold group-hover:gap-2 transition-all">
                Abrir chamada <ArrowRight className="size-3.5" />
              </div>
            </button>
          ) : (
            <div className="rounded-2xl border border-border p-5 text-sm text-muted-foreground flex items-center gap-3">
              <CheckCircle2 className="size-5 text-[#2F9E6B]" />
              <span>
                {calls.length === 0
                  ? 'Vá até "Chamadas" para criar a primeira chamada.'
                  : 'Todas as chamadas estão fechadas.'}
              </span>
            </div>
          )}

          <div className="mt-4 flex gap-2">
            <Button variant="outline" size="sm" className="gap-1.5" onClick={() => navigate('/candidatos')}>
              <Users className="size-3.5" />
              Ver candidatos
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
