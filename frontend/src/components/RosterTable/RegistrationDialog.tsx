import { useEffect, useState, type ReactNode } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { formatCpf } from '@/lib/format';
import { getStatus, STATUSES } from '@/lib/status';
import {
  FetchRegistration,
  ClearApplicationStatus,
  AbsentApplication,
  EnrollApplication,
} from '@/lib/backend';

interface RegistrationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  id: number;
  initialStatus: string;
  hasSelector: boolean;
  onStatusChanged: () => void;
}

// ENEM scores arrive pre-formatted as pt-BR strings (e.g. "655,16"). Parsed
// only to size the magnitude bars below — display always uses the raw string.
const SCORE_MAX = 1000;
function parseScore(value?: string): number | null {
  if (value == null || value === '') return null;
  const n = Number(value.replace(/\./g, '').replace(',', '.'));
  return Number.isFinite(n) ? n : null;
}

function StatusPill({ status }: { status: string }) {
  const def = getStatus(status);
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium shrink-0',
        def.badgeBg,
        def.textColor
      )}
    >
      <span className={cn('w-1.5 h-1.5 rounded-full', def.color)} />
      {def.label}
    </span>
  );
}

function Stat({
  label, value, emphasis,
}: { label: string; value?: string; emphasis?: boolean }) {
  return (
    <div className="flex flex-col gap-0.5 px-3 py-1 min-w-0 first:pl-0 last:pr-0">
      <span className="text-[10px] uppercase tracking-wide text-muted-foreground font-semibold leading-tight">
        {label}
      </span>
      <span
        className={cn(
          'font-mono font-bold tabular-nums text-xl text-foreground',
          emphasis && 'text-primary'
        )}
      >
        {value ?? '—'}
      </span>
    </div>
  );
}

function Field({
  label, value, mono, full,
}: { label: string; value?: any; mono?: boolean; full?: boolean }) {
  if (value == null || value === '') return null;
  return (
    <div className={cn('flex flex-col gap-0.5 min-w-0', full && 'col-span-2')}>
      <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-semibold">
        {label}
      </span>
      <span className={cn('text-sm text-foreground truncate', mono && 'font-mono')}>
        {value}
      </span>
    </div>
  );
}

function ScoreBar({ label, value }: { label: string; value?: string }) {
  const score = parseScore(value);
  const pct = score == null ? 0 : Math.min(100, Math.max(0, (score / SCORE_MAX) * 100));
  return (
    <div className="flex items-center gap-3">
      <span className="text-xs text-muted-foreground w-24 shrink-0">{label}</span>
      <div className="flex-1 h-2 rounded-full bg-border overflow-hidden">
        <div
          className="h-full rounded-full bg-primary transition-[width] motion-reduce:transition-none"
          style={{ width: `${pct}%` }}
        />
      </div>
      <span className="font-mono font-bold tabular-nums text-xs text-foreground w-16 text-right shrink-0">
        {value ?? '—'}
      </span>
    </div>
  );
}

function Section({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="mb-4 last:mb-0">
      <p className="text-[10px] uppercase tracking-widest text-muted-foreground font-semibold mb-2">
        {label}
      </p>
      {children}
    </div>
  );
}

const SCORE_FIELDS: { key: string; label: string }[] = [
  { key: 'Nota Linguagens', label: 'Linguagens' },
  { key: 'Nota Humanas', label: 'Humanas' },
  { key: 'Nota Natureza', label: 'Natureza' },
  { key: 'Nota Matemática', label: 'Matemática' },
  { key: 'Nota Redação', label: 'Redação' },
];

const MUTABLE_STATUSES = STATUSES.filter(
  (s) => s.value === 'APPROVED' || s.value === 'ABSENT' || s.value === 'ENROLLED'
);

export function RegistrationDialog({
  open, onOpenChange, id, initialStatus, hasSelector, onStatusChanged,
}: RegistrationDialogProps) {
  const [detail, setDetail] = useState<Record<string, any>>({});
  const [loading, setLoading] = useState(false);
  const [pendingStatus, setPendingStatus] = useState(initialStatus);
  const [applying, setApplying] = useState(false);

  useEffect(() => {
    if (!open || !id) return;
    setPendingStatus(initialStatus);
    setLoading(true);
    FetchRegistration(id)
      .then((d) => {
        if (!d) return;
        const reg = d.Registration;
        const c = reg?.Candidate;
        const course = d.Course;
        const call = d.Call;

        setDetail({
          Nome: c?.Name,
          'Nome Social': c?.SocialName,
          CPF: formatCpf(c?.CPF),
          'Data de Nascimento': c?.BirthDate,
          Sexo: c?.Sex,
          'Nome da Mãe': c?.MotherName,
          Endereço: [c?.AddressLine, c?.AddressLine2, c?.HouseNumber].filter(Boolean).join(', '),
          Bairro: c?.Neighborhood,
          'Município / UF': c?.Municipality && c?.State ? `${c.Municipality} — ${c.State}` : c?.Municipality,
          CEP: c?.CEP,
          Email: c?.Email,
          'Telefone 1': c?.Phone1,
          'Telefone 2': c?.Phone2,
          'Inscrição ENEM': reg?.EnrollmentID,
          Opção: reg?.Option != null ? `${reg.Option}ª` : undefined,
          Classificação: reg?.Ranking != null ? `${reg.Ranking}º` : undefined,
          'Nota Linguagens': reg?.LanguagesScore,
          'Nota Humanas': reg?.HumanitiesScore,
          'Nota Natureza': reg?.NaturalSciencesScore,
          'Nota Matemática': reg?.MathematicsScore,
          'Nota Redação': reg?.EssayScore,
          'Nota Final': reg?.CompositeScore,
          Turno: course?.Period === 'morning' ? 'Matutino' : course?.Period === 'evening' ? 'Noturno' : course?.Period,
          Cota: course?.Quota,
          Vagas: course?.Seats,
          ...(call ? { Chamada: call.Number } : {}),
        });
      })
      .finally(() => setLoading(false));
  }, [open, id, initialStatus]);

  async function applyStatus() {
    setApplying(true);
    try {
      if (pendingStatus === 'APPROVED') await ClearApplicationStatus(id);
      else if (pendingStatus === 'ABSENT') await AbsentApplication(id);
      else if (pendingStatus === 'ENROLLED') await EnrollApplication(id);
      onStatusChanged();
      onOpenChange(false);
    } finally {
      setApplying(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl max-h-[85vh] flex flex-col gap-0 p-0 overflow-hidden">
        <DialogHeader className="sr-only">
          <DialogTitle>{loading ? 'Carregando candidato' : (detail['Nome'] ?? 'Candidato')}</DialogTitle>
        </DialogHeader>

        {loading ? (
          <div className="flex items-center justify-center h-40 gap-2 text-muted-foreground">
            <Loader2 className="size-4 animate-spin" />
            <span className="text-sm">Carregando...</span>
          </div>
        ) : (
          <>
            {/* Persistent identity + headline stats — does not scroll */}
            <div className="px-6 py-5 border-b border-border shrink-0">
              <div className="flex items-center pr-8">
                <StatusPill status={pendingStatus} />
              </div>
              <h2
                title={detail['Nome']}
                className="font-heading text-xl font-bold text-foreground break-words line-clamp-2 mt-2"
              >
                {detail['Nome']}
              </h2>
              {detail['Nome Social'] && (
                <p className="text-xs text-muted-foreground break-words">{detail['Nome Social']}</p>
              )}
              {detail['CPF'] && (
                <p className="text-xs font-mono text-muted-foreground mt-1">CPF {detail['CPF']}</p>
              )}

              <div className="grid grid-cols-3 divide-x divide-border rounded-xl bg-accent/50 p-3 mt-4">
                <Stat label="Nota Final" value={detail['Nota Final']} emphasis />
                <Stat label="Classificação" value={detail['Classificação']} />
                <Stat label="Opção" value={detail['Opção']} />
              </div>
            </div>

            <div className="flex-1 overflow-y-auto px-6 py-4">
              <Section label="Desempenho por área">
                <div className="flex flex-col gap-1.5">
                  {SCORE_FIELDS.map(({ key, label }) => (
                    <ScoreBar key={key} label={label} value={detail[key]} />
                  ))}
                </div>
              </Section>

              <Separator className="mb-4" />

              <Section label="Documento">
                <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                  <Field label="Data de Nascimento" value={detail['Data de Nascimento']} />
                  <Field label="Sexo" value={detail['Sexo']} />
                  <Field label="Nome da Mãe" value={detail['Nome da Mãe']} full />
                  <Field label="Inscrição ENEM" value={detail['Inscrição ENEM']} mono full />
                </div>
              </Section>

              <Separator className="mb-4" />

              <Section label="Contato">
                <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                  <Field label="Email" value={detail['Email']} full />
                  <Field label="Telefone 1" value={detail['Telefone 1']} />
                  <Field label="Telefone 2" value={detail['Telefone 2']} />
                </div>
              </Section>

              <Separator className="mb-4" />

              <Section label="Endereço">
                <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                  <Field label="Logradouro" value={detail['Endereço']} full />
                  <Field label="Bairro" value={detail['Bairro']} />
                  <Field label="Município / UF" value={detail['Município / UF']} />
                  <Field label="CEP" value={detail['CEP']} mono />
                </div>
              </Section>

              <Separator className="mb-4" />

              <Section label="Curso">
                <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                  <Field label="Turno" value={detail['Turno']} />
                  <Field label="Cota" value={detail['Cota']} />
                  <Field label="Vagas" value={detail['Vagas']} />
                  <Field label="Chamada" value={detail['Chamada']} />
                </div>
              </Section>
            </div>
          </>
        )}

        {hasSelector && !loading && (
          <div className="px-6 py-4 border-t border-border flex items-center gap-3 shrink-0">
            <Select value={pendingStatus} onValueChange={setPendingStatus}>
              <SelectTrigger className="h-8 text-xs flex-1">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {MUTABLE_STATUSES.map((s) => (
                  <SelectItem key={s.value} value={s.value} className="text-xs">
                    {s.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              size="sm"
              className="h-8"
              disabled={pendingStatus === initialStatus || applying}
              onClick={applyStatus}
            >
              {applying ? <Loader2 className="size-3 animate-spin" /> : 'Confirmar'}
            </Button>
          </div>
        )}

        {!hasSelector && (
          <div className="px-6 py-4 border-t border-border flex items-center justify-end shrink-0">
            <Button variant="ghost" size="sm" onClick={() => onOpenChange(false)}>
              Fechar
            </Button>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
