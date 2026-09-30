import { useEffect, useState } from 'react';
import { Download, Loader2, FileText } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import {
  FetchRollCalls,
  FetchApprovedSelection,
  SaveFileDialog,
  WebsitePDF,
  EnrollmentPDF,
  EmailPDF,
  TeacherPDF,
  ExportCSV,
} from '@/lib/backend';

interface RollCall { ID: number; Number: number; Status: string }
interface SelectionInfo { year: number }

const PERIODS = [
  { label: 'Matutino', value: 'morning' },
  { label: 'Noturno',  value: 'evening' },
] as const;

const REPORT_TYPES = [
  { key: 'website',    label: 'Website',    desc: 'Lista de convocados para publicação' },
  { key: 'enrollment', label: 'Convocados', desc: 'Relatório de convocação oficial' },
  { key: 'email',      label: 'E-mails',    desc: 'Lista de contatos dos convocados' },
] as const;

type ReportKey = typeof REPORT_TYPES[number]['key'];

function ExportCard({
  label,
  desc,
  busy,
  onExport,
}: {
  label: string;
  desc: string;
  busy: boolean;
  onExport: () => void;
}) {
  return (
    <button
      disabled={busy}
      onClick={onExport}
      className={cn(
        'w-full text-left rounded-xl border border-border p-4 flex items-center gap-4 transition-colors',
        'hover:bg-accent/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
        busy && 'opacity-60 cursor-not-allowed'
      )}
    >
      <div className="w-9 h-9 rounded-lg bg-muted flex items-center justify-center shrink-0">
        {busy
          ? <Loader2 className="size-4 animate-spin text-muted-foreground" />
          : <Download className="size-4 text-muted-foreground" />
        }
      </div>
      <div>
        <p className="font-semibold text-sm">{label}</p>
        <p className="text-xs text-muted-foreground">{desc}</p>
      </div>
    </button>
  );
}

async function pickSave(title: string, name: string): Promise<string | null> {
  const path = await SaveFileDialog(title, name, '*.pdf', 'PDF (*.pdf)');
  return path || null;
}

export default function Relatorios() {
  const [calls, setCalls] = useState<RollCall[]>([]);
  const [info, setInfo] = useState<SelectionInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [busyKey, setBusyKey] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      const [callsRes, approvedRes] = await Promise.all([FetchRollCalls(), FetchApprovedSelection()]);
      const raw = callsRes ?? [];
      setCalls(raw.map((c) => ({ ID: c.ID, Number: c.Number, Status: c.Status?.toUpperCase?.() ?? 'DONE' })));
      const sel = approvedRes;
      if (sel) setInfo({ year: sel.Year });
      setLoading(false);
    }
    load();
  }, []);

  async function exportReport(
    key: string,
    fn: (path: string) => Promise<any>,
    title: string,
    defaultName: string
  ) {
    const path = await pickSave(title, defaultName);
    if (!path) return;
    setBusyKey(key);
    try {
      await fn(path);
      toast.success(`${title} exportado com sucesso.`);
    } catch (e: any) {
      toast.error(e?.message ?? 'Erro ao exportar.');
    } finally {
      setBusyKey(null);
    }
  }

  async function exportCallReport(type: ReportKey, call: RollCall, period: typeof PERIODS[number]) {
    const key = `${type}-${call.ID}-${period.value}`;
    const label = `${type}_chamada${call.Number}_${period.label.toLowerCase()}_${info?.year ?? ''}`;
    await exportReport(
      key,
      (path) => {
        if (type === 'website') return WebsitePDF(call.ID, period.value, path);
        if (type === 'enrollment') return EnrollmentPDF(call.ID, period.value, path);
        return EmailPDF(call.ID, period.value, path);
      },
      `${REPORT_TYPES.find((r) => r.key === type)?.label} — ${call.Number}ª Chamada (${period.label})`,
      `${label}.pdf`
    );
  }

  async function exportTeacher(period: typeof PERIODS[number]) {
    const key = `teacher-${period.value}`;
    await exportReport(
      key,
      (path) => TeacherPDF(period.value, path),
      `Professor — ${period.label}`,
      `professor_${period.label.toLowerCase()}_${info?.year ?? ''}.pdf`
    );
  }

  async function exportCSV() {
    const path = await SaveFileDialog('Exportar CSV', 'matriculados.csv', '*.csv', 'CSV (*.csv)');
    if (!path) return;
    setBusyKey('csv');
    try {
      await ExportCSV(path);
      toast.success('CSV exportado com sucesso.');
    } catch (e: any) {
      toast.error(e?.message ?? 'Erro ao exportar CSV.');
    } finally {
      setBusyKey(null);
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full text-muted-foreground text-sm gap-2">
        <Loader2 className="size-4 animate-spin" /> Carregando relatórios…
      </div>
    );
  }

  return (
    <div className="px-6 py-6 h-full overflow-auto">
      <div className="mb-6 border-b-2 border-foreground pb-4">
        <h1 className="font-heading font-black text-3xl">Relatórios</h1>
        {info && (
          <p className="text-muted-foreground text-sm mt-0.5">
            SISU <span className="font-mono font-bold text-foreground">{info.year}</span>
          </p>
        )}
      </div>

      {calls.length === 0 ? (
        <div className="flex items-center gap-3 text-muted-foreground text-sm py-8">
          <FileText className="size-5" />
          Nenhuma chamada criada. Os relatórios ficam disponíveis após a primeira chamada.
        </div>
      ) : (
        <div className="space-y-8">
          {/* Per-call reports */}
          {calls.map((call) => (
            <div key={call.ID}>
              <h2 className="font-heading font-bold text-lg mb-3">
                {call.Number}ª Chamada
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {PERIODS.map((period) => (
                  <div key={period.value}>
                    <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
                      {period.label}
                    </p>
                    <div className="space-y-2">
                      {REPORT_TYPES.map((type) => (
                        <ExportCard
                          key={type.key}
                          label={type.label}
                          desc={type.desc}
                          busy={busyKey === `${type.key}-${call.ID}-${period.value}`}
                          onExport={() => exportCallReport(type.key, call, period)}
                        />
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}

          {/* Period-level reports */}
          <div>
            <h2 className="font-heading font-bold text-lg mb-3">Geral</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {PERIODS.map((period) => (
                <div key={period.value}>
                  <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
                    {period.label}
                  </p>
                  <ExportCard
                    label="Professor"
                    desc="Relatório para professores — todos os matriculados"
                    busy={busyKey === `teacher-${period.value}`}
                    onExport={() => exportTeacher(period)}
                  />
                </div>
              ))}
            </div>
          </div>

          {/* CSV */}
          <div>
            <h2 className="font-heading font-bold text-lg mb-3">Exportação de dados</h2>
            <div className="max-w-sm">
              <ExportCard
                label="CSV de matriculados"
                desc="Planilha completa dos candidatos matriculados"
                busy={busyKey === 'csv'}
                onExport={exportCSV}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
