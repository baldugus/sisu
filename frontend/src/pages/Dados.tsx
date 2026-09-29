import { useEffect, useState } from 'react';
import {
  Upload, Download, Trash2, Loader2, FileArchive, AlertTriangle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import {
  FetchApprovedSelection,
  FetchInterestedSelection,
  OpenFileDialog,
  SaveFileDialog,
  LoadApprovedSelection,
  LoadInterestedSelection,
  DeleteApprovedSelection,
  DeleteInterestedSelection,
  Backup,
  Restore,
  Destroy,
} from '@/lib/backend';

interface SelectionInfo { year: number }

function SectionTitle({ children }: { children: React.ReactNode }) {
  return (
    <h2 className="font-heading font-bold text-lg mb-4">{children}</h2>
  );
}

function DataRow({
  label,
  desc,
  value,
  onAction,
  actionLabel,
  actionIcon: Icon,
  actionVariant,
  onDelete,
  busy,
}: {
  label: string;
  desc?: string;
  value?: string;
  onAction: () => void;
  actionLabel: string;
  actionIcon: React.ElementType;
  actionVariant?: 'default' | 'outline' | 'destructive' | 'ghost';
  onDelete?: () => void;
  busy?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-4 py-4 border-b border-border last:border-0">
      <div className="flex-1 min-w-0">
        <p className="font-medium text-sm">{label}</p>
        {desc && <p className="text-xs text-muted-foreground mt-0.5">{desc}</p>}
        {value && (
          <p className="text-xs font-mono text-foreground mt-0.5 truncate">{value}</p>
        )}
      </div>
      <div className="flex items-center gap-2 shrink-0">
        {onDelete && (
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onClick={onDelete}
            disabled={busy}
          >
            <Trash2 className="size-3.5" />
          </Button>
        )}
        <Button
          variant={actionVariant ?? 'outline'}
          size="sm"
          className="h-8 gap-1.5"
          onClick={onAction}
          disabled={busy}
        >
          {busy ? <Loader2 className="size-3.5 animate-spin" /> : <Icon className="size-3.5" />}
          {actionLabel}
        </Button>
      </div>
    </div>
  );
}

/* ─── Import modal ─────────────────────────────────────────────────── */
interface ImportModalProps {
  type: 'approved' | 'waitlisted';
  open: boolean;
  onClose: () => void;
  approvedInfo: SelectionInfo | null;
  onSuccess: () => void;
}

function ImportModal({ type, open, onClose, approvedInfo, onSuccess }: ImportModalProps) {
  const [year, setYear] = useState('');
  const [filePath, setFilePath] = useState('');
  const [busy, setBusy] = useState(false);

  const isWaitlisted = type === 'waitlisted';
  const resolvedYear = isWaitlisted ? String(approvedInfo?.year ?? '') : year;

  async function pickFile() {
    const path = await OpenFileDialog('Selecionar CSV', '*.csv', 'CSV (*.csv)');
    if (path) setFilePath(path);
  }

  async function submit() {
    if (!filePath || !resolvedYear) {
      toast.error('Preencha todos os campos antes de importar.');
      return;
    }
    setBusy(true);
    try {
      const y = Number(resolvedYear);
      if (isWaitlisted) await LoadInterestedSelection(y, filePath);
      else await LoadApprovedSelection(y, filePath);
      toast.success('Dados importados com sucesso.');
      onSuccess();
      onClose();
    } catch (e: any) {
      toast.error(e?.message ?? 'Falha ao importar.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={(o) => { if (!o) onClose(); }}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="font-heading font-bold">
            Importar {isWaitlisted ? 'em espera' : 'aprovados'}
          </DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-4 py-2">
          {isWaitlisted && approvedInfo ? (
            <div className="text-xs text-muted-foreground rounded-lg bg-muted px-3 py-2">
              Usando ano <span className="font-mono font-bold text-foreground">{approvedInfo.year}</span>{' '}
              da seleção de aprovados.
            </div>
          ) : (
            <div>
              <Label className="text-xs mb-1 block">Ano</Label>
              <Input
                placeholder="2025"
                className="h-8 text-sm font-mono"
                value={year}
                onChange={(e) => setYear(e.target.value)}
                maxLength={4}
              />
            </div>
          )}

          <div>
            <Label className="text-xs mb-1 block">Arquivo CSV</Label>
            <div className="flex gap-2">
              <Input
                readOnly
                placeholder="Nenhum arquivo selecionado…"
                className="h-8 text-xs flex-1 truncate"
                value={filePath}
              />
              <Button variant="outline" size="sm" className="h-8 shrink-0" onClick={pickFile}>
                Selecionar
              </Button>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="ghost" size="sm" onClick={onClose} disabled={busy}>Cancelar</Button>
          <Button size="sm" onClick={submit} disabled={busy || !filePath}>
            {busy && <Loader2 className="size-3.5 animate-spin mr-1" />}
            Importar
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

/* ─── Main page ─────────────────────────────────────────────────────── */
export default function Dados() {
  const [approvedInfo, setApprovedInfo] = useState<SelectionInfo | null>(null);
  const [waitlistedInfo, setWaitlistedInfo] = useState<SelectionInfo | null>(null);
  const [busy, setBusy] = useState<string | null>(null);
  const [importModal, setImportModal] = useState<'approved' | 'waitlisted' | null>(null);
  const [destroyOpen, setDestroyOpen] = useState(false);
  const [destroyConfirm, setDestroyConfirm] = useState('');

  async function loadInfo() {
    const [a, w] = await Promise.all([FetchApprovedSelection(), FetchInterestedSelection()]);
    setApprovedInfo(a ? { year: a.Year } : null);
    setWaitlistedInfo(w ? { year: w.Year } : null);
  }

  useEffect(() => { loadInfo(); }, []);

  async function act(key: string, fn: () => Promise<any>, successMsg: string) {
    setBusy(key);
    try {
      await fn();
      toast.success(successMsg);
      await loadInfo();
    } catch (e: any) {
      toast.error(e?.message ?? 'Ocorreu um erro.');
    } finally {
      setBusy(null);
    }
  }

  async function doBackup() {
    const path = await SaveFileDialog('Salvar backup', 'sisu_backup.sqlite', '*.sqlite;*.db', 'SQLite (*.sqlite, *.db)');
    if (!path) return;
    await act('backup', () => Backup(path), 'Backup salvo com sucesso.');
  }

  async function doRestore() {
    const path = await OpenFileDialog('Selecionar backup', '*.sqlite;*.db', 'SQLite (*.sqlite, *.db)');
    if (!path) return;
    await act('restore', () => Restore(path), 'Dados restaurados com sucesso.');
  }

  async function doDestroy() {
    await act('destroy', Destroy, 'Todos os dados foram apagados.');
    setDestroyOpen(false);
    setDestroyConfirm('');
  }

  const QUOTAS = [
    { code: 'AC', label: 'Ampla Concorrência' },
    { code: 'C1', label: 'Escola Pública, renda familiar bruta per capita ≤ 1,5 SM' },
    { code: 'C2', label: 'Escola Pública, renda ≤ 1,5 SM, autodeclarado PPI' },
    { code: 'C3', label: 'Escola Pública, autodeclarado PPI' },
  ];

  function selectionLabel(info: SelectionInfo) {
    return `${info.year}`;
  }

  return (
    <div className="px-6 py-6 h-full overflow-auto">
      <div className="mb-6 border-b-2 border-foreground pb-4">
        <h1 className="font-heading font-black text-3xl">Dados</h1>
        <p className="text-muted-foreground text-sm mt-0.5">
          Importação de dados, backup e configurações do ciclo.
        </p>
      </div>

      <div className="max-w-2xl space-y-8">
        {/* ── Import ── */}
        <div>
          <SectionTitle>Importação de seleção</SectionTitle>
          <div className="rounded-2xl border border-border">
            <DataRow
              label="Aprovados"
              desc="CSV dos candidatos convocados no SISU"
              value={approvedInfo ? selectionLabel(approvedInfo) : undefined}
              actionLabel={approvedInfo ? 'Reimportar' : 'Importar'}
              actionIcon={Upload}
              onAction={() => setImportModal('approved')}
              onDelete={approvedInfo ? () => act('del-approved', DeleteApprovedSelection, 'Seleção de aprovados removida.') : undefined}
              busy={busy === 'del-approved'}
            />
            <DataRow
              label="Em espera"
              desc={approvedInfo ? 'CSV dos candidatos em lista de espera' : 'Importe aprovados primeiro'}
              value={waitlistedInfo ? selectionLabel(waitlistedInfo) : undefined}
              actionLabel={waitlistedInfo ? 'Reimportar' : 'Importar'}
              actionIcon={Upload}
              actionVariant={approvedInfo ? 'outline' : 'ghost'}
              onAction={() => {
                if (!approvedInfo) {
                  toast.error('Importe a seleção de aprovados antes de importar a lista de espera.');
                  return;
                }
                setImportModal('waitlisted');
              }}
              onDelete={waitlistedInfo ? () => act('del-waitlisted', DeleteInterestedSelection, 'Lista de espera removida.') : undefined}
              busy={busy === 'del-waitlisted'}
            />
          </div>
        </div>

        {/* ── Backup ── */}
        <div>
          <SectionTitle>Backup e restauração</SectionTitle>
          <div className="rounded-2xl border border-border">
            <DataRow
              label="Salvar backup"
              desc="Exporta o banco de dados local para um arquivo SQLite"
              actionLabel="Salvar"
              actionIcon={Download}
              onAction={doBackup}
              busy={busy === 'backup'}
            />
            <DataRow
              label="Restaurar backup"
              desc="Substitui os dados atuais pelos de um arquivo de backup"
              actionLabel="Restaurar"
              actionIcon={FileArchive}
              onAction={doRestore}
              busy={busy === 'restore'}
            />
          </div>
        </div>

        {/* ── Quota reference ── */}
        <div>
          <SectionTitle>Referência de cotas</SectionTitle>
          <div className="rounded-2xl border border-border divide-y divide-border">
            {QUOTAS.map(({ code, label }) => (
              <div key={code} className="flex gap-4 px-5 py-3">
                <span className="font-mono font-black text-sm w-8 shrink-0 text-primary">{code}</span>
                <span className="text-sm text-muted-foreground">{label}</span>
              </div>
            ))}
          </div>
        </div>

        {/* ── Danger zone ── */}
        <div>
          <SectionTitle>Zona de perigo</SectionTitle>
          <div className="rounded-2xl border-2 border-destructive/40 p-5">
            <div className="flex items-start gap-3">
              <AlertTriangle className="size-5 text-destructive shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="font-semibold text-sm">Apagar todos os dados</p>
                <p className="text-xs text-muted-foreground mt-0.5">
                  Remove permanentemente todo o banco de dados — seleções, chamadas e matrículas.
                  Essa ação não pode ser desfeita.
                </p>
              </div>
              <Button
                variant="destructive"
                size="sm"
                className="h-8 shrink-0"
                onClick={() => setDestroyOpen(true)}
              >
                Apagar tudo
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* ── Import modals ── */}
      {importModal && (
        <ImportModal
          type={importModal}
          open
          onClose={() => setImportModal(null)}
          approvedInfo={approvedInfo}
          onSuccess={loadInfo}
        />
      )}

      {/* ── Destroy confirm ── */}
      <AlertDialog open={destroyOpen} onOpenChange={setDestroyOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="font-heading font-bold text-destructive">
              Apagar todos os dados?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Esta ação é permanente e irrecuperável. Todos os dados do ciclo — seleções,
              chamadas e registros de matrícula — serão apagados.
              <br /><br />
              Digite <strong>APAGAR</strong> para confirmar:
            </AlertDialogDescription>
          </AlertDialogHeader>

          <Input
            className={cn('font-mono', destroyConfirm === 'APAGAR' && 'border-destructive')}
            placeholder="APAGAR"
            value={destroyConfirm}
            onChange={(e) => setDestroyConfirm(e.target.value)}
          />

          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => { setDestroyConfirm(''); }}>
              Cancelar
            </AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive hover:bg-destructive/90 text-white"
              disabled={destroyConfirm !== 'APAGAR' || busy === 'destroy'}
              onClick={doDestroy}
            >
              {busy === 'destroy' && <Loader2 className="size-3.5 animate-spin mr-1" />}
              Apagar tudo
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
