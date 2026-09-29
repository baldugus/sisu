import { useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { OpenFileDialog, LoadApprovedSelection } from "@/lib/backend";

interface Props {
  onClose: () => void;
  onSuccess: (msg: string) => void;
  onError: (msg: string) => void;
}

export default function ApprovedImportModal({ onClose, onSuccess, onError }: Props) {
  const [year, setYear] = useState("");
  const [filePath, setFilePath] = useState("");
  const [busy, setBusy] = useState(false);

  async function pickFile() {
    const path = await OpenFileDialog("Selecionar CSV", "*.csv", "CSV (*.csv)");
    if (path) setFilePath(path);
  }

  async function submit() {
    if (!filePath || !year) return;
    setBusy(true);
    try {
      await LoadApprovedSelection(Number(year), filePath);
      onSuccess("Aprovados importados com sucesso.");
    } catch (e: any) {
      onError(e?.message ?? "Falha ao importar aprovados.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <Dialog open onOpenChange={(o) => { if (!o) onClose(); }}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Importar aprovados</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-4 py-2">
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
          <Button size="sm" onClick={submit} disabled={busy || !filePath || !year}>
            {busy && <Loader2 className="size-3.5 animate-spin mr-1" />}
            Importar
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
