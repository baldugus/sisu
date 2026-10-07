import { useState, useMemo, useCallback } from 'react';
import { Search, ChevronsUpDown, Loader2, ArrowUpDown, Copy } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';
import {
  getStatus,
  STATUSES,
  ENTRY_STATUSES,
  KIND_LABELS,
  semesterLabel,
  statusToOutcome,
  type EntryKind,
} from '@/lib/status';
import { copyText } from '@/lib/clipboard';
import { toast } from 'sonner';
import { RegistrationDialog } from './RegistrationDialog';
import { CopyEmailButton } from './CopyEmailButton';
import { SetCallEntryOutcome } from '@/lib/backend';

export interface RowData {
  ID: number;
  Name: string;
  CPF: string;
  Email?: string;
  Period: string;
  Quota: string;
  Status: string;
  EnrollmentID?: string;
  Ranking?: number;
  /** Semester the row refers to: the entry's semester in a call, or the
   *  registration's current semester elsewhere. */
  Semester?: number;
  /** Call entry kind (call view only). */
  Kind?: EntryKind;
}

interface RosterTableProps {
  rows: RowData[];
  loading?: boolean;
  hasSelector?: boolean;
  /** Call the rows belong to; enables outcome editing and the call columns. */
  callId?: number;
  /** Closed call: rows can be inspected but not changed. */
  readOnly?: boolean;
  onRefresh?: () => void;
  emptyMessage?: string;
  showKindFilter?: boolean;
  kind?: 'approved' | 'waitlisted';
  onKindChange?: (kind: 'approved' | 'waitlisted') => void;
  showContact?: boolean;
}

type SortKey = 'Ranking' | 'Name';

function FilterPills<T extends string>({
  options,
  value,
  onChange,
}: {
  options: { label: string; value: T }[];
  value: T;
  onChange: (v: T) => void;
}) {
  return (
    <div className="flex gap-1 flex-wrap">
      {options.map((opt) => (
        <button
          key={opt.value}
          onClick={() => onChange(opt.value)}
          className={cn(
            'px-2.5 py-1 rounded-full text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
            value === opt.value
              ? 'bg-primary text-primary-foreground'
              : 'bg-secondary text-secondary-foreground hover:bg-accent hover:text-accent-foreground'
          )}
        >
          {opt.label}
        </button>
      ))}
    </div>
  );
}

const PERIOD_OPTIONS = [
  { label: 'Todos os turnos', value: '' as const },
  { label: 'Matutino', value: 'Matutino' as const },
  { label: 'Noturno', value: 'Noturno' as const },
] as { label: string; value: string }[];

const QUOTA_ALL = '__all__';

const SEMESTER_OPTIONS = [
  { label: 'Todos os semestres', value: '' },
  { label: '1º semestre', value: '1' },
  { label: '2º semestre', value: '2' },
];


export function RosterTable({
  rows,
  loading = false,
  hasSelector = false,
  callId,
  readOnly = false,
  onRefresh,
  emptyMessage = 'Nenhum candidato encontrado.',
  showKindFilter = false,
  kind,
  onKindChange,
  showContact = false,
}: RosterTableProps) {
  const [search, setSearch] = useState('');
  const [filterPeriod, setFilterPeriod] = useState('');
  const [filterQuota, setFilterQuota] = useState('');
  const [filterStatus, setFilterStatus] = useState('');
  const [filterSemester, setFilterSemester] = useState('');
  const [sortKey, setSortKey] = useState<SortKey>('Ranking');
  const [sortAsc, setSortAsc] = useState(true);

  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [bulkStatus, setBulkStatus] = useState('APPROVED');
  const [bulkApplying, setBulkApplying] = useState(false);

  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogRow, setDialogRow] = useState<RowData | null>(null);

  const quotaOptions = useMemo(() => {
    const set = new Set<string>();
    for (const r of rows) if (r.Quota) set.add(r.Quota);
    return Array.from(set).sort((a, b) => a.localeCompare(b, 'pt-BR'));
  }, [rows]);

  const filtered = useMemo(() => {
    let out = rows;
    const q = search.toLowerCase();
    if (q) out = out.filter((r) => r.Name.toLowerCase().includes(q) || r.CPF.includes(q));
    if (filterPeriod) out = out.filter((r) => r.Period === filterPeriod);
    if (filterQuota) out = out.filter((r) => r.Quota === filterQuota);
    if (filterStatus) out = out.filter((r) => r.Status === filterStatus);
    if (filterSemester) out = out.filter((r) => String(r.Semester ?? '') === filterSemester);
    out = [...out].sort((a, b) => {
      if (sortKey === 'Ranking') {
        const av = a.Ranking ?? Infinity;
        const bv = b.Ranking ?? Infinity;
        return sortAsc ? av - bv : bv - av;
      }
      return sortAsc
        ? a.Name.localeCompare(b.Name, 'pt-BR')
        : b.Name.localeCompare(a.Name, 'pt-BR');
    });
    return out;
  }, [rows, search, filterPeriod, filterQuota, filterStatus, filterSemester, sortKey, sortAsc]);

  const toggleSort = useCallback(
    (key: SortKey) => {
      if (sortKey === key) setSortAsc((a) => !a);
      else { setSortKey(key); setSortAsc(true); }
    },
    [sortKey]
  );

  const inCall = callId != null;
  const editable = inCall && hasSelector && !readOnly;
  const colCount = 7 + (editable ? 1 : 0) + (showContact ? 1 : 0) + (inCall ? 1 : 0);

  const statusOptions = useMemo(() => {
    const present = new Set(rows.map((r) => r.Status));
    return [
      { label: 'Todos os status', value: '' },
      ...STATUSES.filter((s) => present.has(s.value)).map((s) => ({ label: s.label, value: s.value as string })),
    ];
  }, [rows]);

  const filteredEmails = useMemo(() => {
    const seen = new Set<string>();
    const out: string[] = [];
    for (const r of filtered) {
      const email = r.Email?.trim();
      if (!email) continue;
      const key = email.toLowerCase();
      if (seen.has(key)) continue;
      seen.add(key);
      out.push(email);
    }
    return out;
  }, [filtered]);

  async function copyAllEmails() {
    if (filteredEmails.length === 0) {
      toast.error('Nenhum e-mail para copiar');
      return;
    }
    const ok = await copyText(filteredEmails.join(', '));
    if (!ok) {
      toast.error('Não foi possível copiar os e-mails');
      return;
    }
    const skipped = filtered.length - filteredEmails.length;
    const n = filteredEmails.length;
    const base = `${n} e-mail${n !== 1 ? 's' : ''} copiado${n !== 1 ? 's' : ''}`;
    toast.success(skipped > 0 ? `${base} (${skipped} sem e-mail)` : base);
  }

  const allSelected = filtered.length > 0 && filtered.every((r) => selectedIds.has(r.ID));
  const someSelected = selectedIds.size > 0;

  function toggleAll() {
    if (allSelected) setSelectedIds(new Set());
    else setSelectedIds(new Set(filtered.map((r) => r.ID)));
  }

  function toggleRow(id: number) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  }

  async function applyBulk() {
    if (callId == null) return;
    const selected = filtered.filter((r) => selectedIds.has(r.ID));
    const outcome = statusToOutcome(bulkStatus);
    if (!outcome) return;
    setBulkApplying(true);
    try {
      for (const row of selected) {
        if (row.Status === bulkStatus) continue;
        await SetCallEntryOutcome(callId, row.ID, outcome);
      }
      setSelectedIds(new Set());
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setBulkApplying(false);
      onRefresh?.();
    }
  }

  function SortBtn({ col }: { col: SortKey }) {
    const active = sortKey === col;
    return (
      <button
        onClick={() => toggleSort(col)}
        className="inline-flex items-center gap-1 hover:text-foreground transition-colors"
      >
        {active ? (
          <ChevronsUpDown className="size-3.5" />
        ) : (
          <ArrowUpDown className="size-3.5 opacity-40" />
        )}
      </button>
    );
  }

  return (
    <div className="flex flex-col gap-3 h-full min-h-0">
      {/* ── Toolbar ── */}
      <div className="flex flex-col gap-2 px-6 pt-5 pb-3 border-b border-border">
        {/* Kind filter (Candidatos only) */}
        {showKindFilter && kind && onKindChange && (
          <div className="flex gap-2 mb-1">
            <button
              onClick={() => onKindChange('approved')}
              className={cn(
                'px-4 py-1.5 rounded-full text-sm font-semibold transition-colors',
                kind === 'approved'
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-secondary text-secondary-foreground hover:bg-accent hover:text-accent-foreground'
              )}
            >
              Aprovados
            </button>
            <button
              onClick={() => onKindChange('waitlisted')}
              className={cn(
                'px-4 py-1.5 rounded-full text-sm font-semibold transition-colors',
                kind === 'waitlisted'
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-secondary text-secondary-foreground hover:bg-accent hover:text-accent-foreground'
              )}
            >
              Em espera
            </button>
          </div>
        )}

        {/* Search */}
        <div className="relative max-w-xs">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground" />
          <Input
            placeholder="Buscar por nome ou CPF…"
            className="pl-8 h-8 text-sm"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {/* Filter pills */}
        <div className="flex flex-wrap items-center gap-y-1.5 gap-x-4">
          <FilterPills options={PERIOD_OPTIONS} value={filterPeriod} onChange={setFilterPeriod} />

          <Select
            value={filterQuota || QUOTA_ALL}
            onValueChange={(v) => setFilterQuota(v === QUOTA_ALL ? '' : v)}
          >
            <SelectTrigger className="h-7 text-xs w-[220px]">
              <SelectValue placeholder="Todas as cotas">
                {filterQuota || 'Todas as cotas'}
              </SelectValue>
            </SelectTrigger>
            <SelectContent className="max-w-sm">
              <SelectItem value={QUOTA_ALL} className="text-xs">
                Todas as cotas
              </SelectItem>
              {quotaOptions.map((q) => (
                <SelectItem
                  key={q}
                  value={q}
                  className="items-start whitespace-normal py-1.5 text-xs"
                >
                  <span className="line-clamp-2">{q}</span>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <FilterPills options={statusOptions} value={filterStatus} onChange={setFilterStatus} />

          <FilterPills options={SEMESTER_OPTIONS} value={filterSemester} onChange={setFilterSemester} />
        </div>
      </div>

      {/* ── Table ── */}
      <div className="flex-1 overflow-auto px-6">
        <table className="w-full border-collapse text-sm">
          <thead className="sticky top-0 bg-panel z-10">
            <tr>
              {editable && (
                <th className="w-8 px-2 py-3 text-left">
                  <Checkbox
                    checked={allSelected}
                    onCheckedChange={toggleAll}
                    aria-label="Selecionar todos"
                    className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                  />
                </th>
              )}
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                <span className="flex items-center gap-1">
                  # <SortBtn col="Ranking" />
                </span>
              </th>
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                <span className="flex items-center gap-1">
                  Nome <SortBtn col="Name" />
                </span>
              </th>
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                CPF
              </th>
              {showContact && (
                <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                  E-mail
                </th>
              )}
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                Turno
              </th>
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                Cota
              </th>
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                Sem.
              </th>
              {inCall && (
                <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                  Origem
                </th>
              )}
              <th className="px-2 py-3 text-left font-semibold text-xs uppercase tracking-wide text-muted-foreground">
                Status
              </th>
            </tr>
            <tr>
              <td colSpan={colCount} className="p-0">
                <div className="h-px bg-border" />
              </td>
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr>
                <td colSpan={colCount} className="py-16 text-center text-muted-foreground">
                  <div className="flex items-center justify-center gap-2">
                    <Loader2 className="size-4 animate-spin" />
                    <span className="text-sm">Carregando…</span>
                  </div>
                </td>
              </tr>
            )}
            {!loading && filtered.length === 0 && (
              <tr>
                <td colSpan={colCount} className="py-16 text-center text-muted-foreground text-sm">
                  {emptyMessage}
                </td>
              </tr>
            )}
            {!loading &&
              filtered.map((row, i) => {
                const status = getStatus(row.Status);
                const selected = selectedIds.has(row.ID);
                return (
                  <tr
                    key={row.ID}
                    onClick={() => { setDialogRow(row); setDialogOpen(true); }}
                    className={cn(
                      'fade-in border-b border-border cursor-pointer transition-colors group',
                      selected ? 'bg-accent/60' : 'hover:bg-muted/60'
                    )}
                    style={{ animationDelay: `${Math.min(i * 12, 200)}ms` }}
                  >
                    {editable && (
                      <td className="px-2 py-2.5" onClick={(e) => e.stopPropagation()}>
                        <Checkbox
                          checked={selected}
                          onCheckedChange={() => toggleRow(row.ID)}
                          className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                        />
                      </td>
                    )}
                    <td className="px-2 py-2.5 font-mono text-muted-foreground text-xs tabular-nums w-10">
                      {row.Ranking ?? '—'}
                    </td>
                    <td className="px-2 py-2.5 font-medium max-w-[200px] truncate">{row.Name}</td>
                    <td className="px-2 py-2.5 font-mono text-xs text-muted-foreground">{row.CPF}</td>
                    {showContact && (
                      <td className="px-2 py-2.5">
                        {row.Email ? (
                          <div className="flex items-center gap-1">
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <span
                                  className="block max-w-[200px] truncate font-mono text-xs text-foreground cursor-default"
                                  onClick={(e) => e.stopPropagation()}
                                >
                                  {row.Email}
                                </span>
                              </TooltipTrigger>
                              <TooltipContent className="max-w-sm">{row.Email}</TooltipContent>
                            </Tooltip>
                            <CopyEmailButton email={row.Email} />
                          </div>
                        ) : (
                          <span className="text-xs text-muted-foreground">—</span>
                        )}
                      </td>
                    )}
                    <td className="px-2 py-2.5 text-sm text-muted-foreground">{row.Period}</td>
                    <td className="px-2 py-2.5">
                      {row.Quota ? (
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span
                              className="block max-w-[220px] truncate text-xs text-foreground cursor-default"
                              onClick={(e) => e.stopPropagation()}
                            >
                              {row.Quota}
                            </span>
                          </TooltipTrigger>
                          <TooltipContent className="max-w-sm">{row.Quota}</TooltipContent>
                        </Tooltip>
                      ) : (
                        <span className="text-xs text-muted-foreground">—</span>
                      )}
                    </td>
                    <td className="px-2 py-2.5 text-sm text-muted-foreground tabular-nums">
                      {semesterLabel(row.Semester)}
                    </td>
                    {inCall && (
                      <td className="px-2 py-2.5 text-xs whitespace-nowrap">
                        <span className="text-muted-foreground">{row.Kind ? KIND_LABELS[row.Kind] : '—'}</span>
                      </td>
                    )}
                    <td className="px-2 py-2.5">
                      <span
                        className={cn(
                          'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium',
                          status.badgeBg,
                          status.textColor
                        )}
                      >
                        <span className={cn('w-1.5 h-1.5 rounded-full', status.color)} />
                        {status.label}
                      </span>
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>
      </div>

      {/* ── Footer: count + bulk action ── */}
      <div className="px-6 py-3 border-t border-border flex items-center justify-between shrink-0">
        <span className="text-xs text-muted-foreground">
          {filtered.length} candidato{filtered.length !== 1 ? 's' : ''}
          {rows.length !== filtered.length && ` de ${rows.length}`}
        </span>

        {editable && someSelected && (
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">
              {selectedIds.size} selecionado{selectedIds.size !== 1 ? 's' : ''}
            </span>
            <Select value={bulkStatus} onValueChange={setBulkStatus}>
              <SelectTrigger className="h-7 text-xs w-44">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {ENTRY_STATUSES.map((s) => (
                  <SelectItem key={s.value} value={s.value} className="text-xs">
                    {s.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              size="sm"
              className="h-7 text-xs"
              disabled={bulkApplying}
              onClick={applyBulk}
            >
              {bulkApplying ? <Loader2 className="size-3 animate-spin mr-1" /> : null}
              Aplicar
            </Button>
          </div>
        )}

        {showContact && (
          <Button
            variant="outline"
            size="sm"
            className="h-7 text-xs"
            onClick={copyAllEmails}
          >
            <Copy className="size-3.5" />
            Copiar e-mails
          </Button>
        )}
      </div>

      {/* ── Registration detail dialog ── */}
      {dialogRow && (
        <RegistrationDialog
          open={dialogOpen}
          onOpenChange={(o) => { if (!o) setDialogRow(null); setDialogOpen(o); }}
          id={dialogRow.ID}
          initialStatus={dialogRow.Status}
          hasSelector={editable}
          callId={callId}
          onStatusChanged={() => onRefresh?.()}
        />
      )}
    </div>
  );
}
