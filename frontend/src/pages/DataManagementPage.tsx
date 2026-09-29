import { toast } from "sonner";
import { useState, useEffect } from "react";
import {
  Backup,
  DeleteApprovedSelection,
  DeleteInterestedSelection,
  Destroy,
  ExportCSV,
  FetchApprovedSelection,
  FetchInterestedSelection,
  OpenFileDialog,
  Restore,
  SaveFileDialog,
} from "@/lib/backend";
import ApprovedImportModal from "./components/ApprovedImportModal";
import InterestedImportModal from "./components/InterestedImportModal";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Upload,
  Trash2,
  Database,
  Download,
  FileSpreadsheet,
  AlertTriangle,
} from "lucide-react";

const DataManagementPage = () => {
  const [openApprovedModal, setOpenApprovedModal] = useState(false);
  const [openInterestedModal, setOpenInterestedModal] = useState(false);
  const [hasApprovedData, setHasApprovedData] = useState(false);
  const [hasInterestedData, setHasInterestedData] = useState(false);

  const checkDataExists = async () => {
    try {
      const approved = await FetchApprovedSelection();
      setHasApprovedData(approved != null);
    } catch {
      setHasApprovedData(false);
    }
    try {
      const interested = await FetchInterestedSelection();
      setHasInterestedData(interested != null);
    } catch {
      setHasInterestedData(false);
    }
  };

  useEffect(() => {
    checkDataExists();
  }, []);

  const handleDeleteApproved = async () => {
    try {
      await DeleteApprovedSelection();
      toast.success("Aprovados removidos com sucesso.");
      setHasApprovedData(false);
    } catch (e: any) {
      toast.error(e?.message || "Falha ao remover aprovados.");
    }
  };

  const handleDeleteInterested = async () => {
    try {
      await DeleteInterestedSelection();
      toast.success("Em espera removidos com sucesso.");
      setHasInterestedData(false);
    } catch (e: any) {
      toast.error(e?.message || "Falha ao remover em espera.");
    }
  };

  const handleRestore = async () => {
    try {
      const filePath = await OpenFileDialog(
        "Selecionar backup",
        "*.db;*.sqlite",
        "Banco de dados (*.db, *.sqlite)"
      );
      if (!filePath) return;
      await Restore(filePath);
      toast.success("Banco de dados restaurado com sucesso!");
      checkDataExists();
    } catch (e: any) {
      toast.error(e?.message || "Falha ao restaurar banco de dados.");
    }
  };

  const handleBackup = async () => {
    try {
      const filePath = await SaveFileDialog(
        "Salvar backup",
        "backup.db",
        "*.db",
        "Banco de dados (*.db)"
      );
      if (!filePath) return;
      await Backup(filePath);
      toast.success("Backup realizado com sucesso!");
    } catch (e: any) {
      toast.error(e?.message || "Falha ao realizar backup.");
    }
  };

  const handleExportCSV = async () => {
    try {
      const filePath = await SaveFileDialog(
        "Exportar CSV",
        "matriculados.csv",
        "*.csv",
        "CSV (*.csv)"
      );
      if (!filePath) return;
      await ExportCSV(filePath);
      toast.success("CSV exportado com sucesso!");
    } catch (e: any) {
      toast.error(e?.message || "Falha ao exportar CSV.");
    }
  };

  const handleDestroy = async () => {
    try {
      await Destroy();
      toast.success("Todos os dados foram destruídos.");
      setHasApprovedData(false);
      setHasInterestedData(false);
    } catch (e: any) {
      toast.error(e?.message || "Falha ao destruir dados.");
    }
  };

  return (
    <div className="flex-1 space-y-6 p-8 pt-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-foreground">
          Gerenciamento de Dados
        </h2>
      </div>

      {/* Card 1: CSV Imports */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Importação de Dados</CardTitle>
          <CardDescription>
            Importe os arquivos CSV do SiSU para o sistema.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Aprovados row */}
          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="flex items-center gap-3">
              <span
                className={`h-2.5 w-2.5 rounded-full ${
                  hasApprovedData ? "bg-green-500" : "bg-gray-300"
                }`}
              />
              <div>
                <p className="text-sm font-medium">Aprovados</p>
                <p className="text-xs text-muted-foreground">
                  {hasApprovedData
                    ? "Dados importados"
                    : "Nenhum dado importado"}
                </p>
              </div>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setOpenApprovedModal(true)}
              >
                <Upload className="mr-2 h-4 w-4" />
                Importar
              </Button>
              <Button
                variant="destructive"
                size="sm"
                disabled={!hasApprovedData}
                onClick={handleDeleteApproved}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Remover
              </Button>
            </div>
          </div>

          {/* Em Espera row */}
          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="flex items-center gap-3">
              <span
                className={`h-2.5 w-2.5 rounded-full ${
                  hasInterestedData ? "bg-green-500" : "bg-gray-300"
                }`}
              />
              <div>
                <p className="text-sm font-medium">Em Espera</p>
                <p className="text-xs text-muted-foreground">
                  {hasInterestedData
                    ? "Dados importados"
                    : "Nenhum dado importado"}
                </p>
              </div>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setOpenInterestedModal(true)}
              >
                <Upload className="mr-2 h-4 w-4" />
                Importar
              </Button>
              <Button
                variant="destructive"
                size="sm"
                disabled={!hasInterestedData}
                onClick={handleDeleteInterested}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Remover
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Card 2: Database & Export */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Banco de Dados & Exportação</CardTitle>
          <CardDescription>
            Gerencie backups do banco de dados e exporte relatórios.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="flex items-center gap-3">
              <Database className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Importar Banco de Dados</p>
                <p className="text-xs text-muted-foreground">
                  Restaurar um backup do banco de dados existente.
                </p>
              </div>
            </div>
            <Button variant="outline" size="sm" onClick={handleRestore}>
              <Upload className="mr-2 h-4 w-4" />
              Restaurar
            </Button>
          </div>

          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="flex items-center gap-3">
              <Download className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Exportar Banco de Dados</p>
                <p className="text-xs text-muted-foreground">
                  Criar um backup completo do banco de dados.
                </p>
              </div>
            </div>
            <Button variant="outline" size="sm" onClick={handleBackup}>
              <Download className="mr-2 h-4 w-4" />
              Backup
            </Button>
          </div>

          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="flex items-center gap-3">
              <FileSpreadsheet className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">
                  Exportar CSV de Matriculados
                </p>
                <p className="text-xs text-muted-foreground">
                  Gerar um arquivo CSV com os candidatos matriculados.
                </p>
              </div>
            </div>
            <Button variant="outline" size="sm" onClick={handleExportCSV}>
              <Download className="mr-2 h-4 w-4" />
              Exportar
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Card 3: Danger Zone */}
      <Card className="border-destructive/50">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            Zona de Perigo
          </CardTitle>
          <CardDescription>
            Ações irreversíveis. Prossiga com cuidado.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between rounded-lg border border-destructive/30 bg-destructive/5 p-4">
            <div>
              <p className="text-sm font-medium">Destruir Todos os Dados</p>
              <p className="text-xs text-muted-foreground">
                Remove permanentemente todos os dados do sistema. Esta ação não
                pode ser desfeita.
              </p>
            </div>
            <Button variant="destructive" size="sm" onClick={handleDestroy}>
              <Trash2 className="mr-2 h-4 w-4" />
              Destruir
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Modals */}
      {openApprovedModal && (
        <ApprovedImportModal
          onClose={() => setOpenApprovedModal(false)}
          onSuccess={(msg: string) => {
            toast.success(msg || "Importado!");
            setOpenApprovedModal(false);
            setHasApprovedData(true);
          }}
          onError={(msg: string) => toast.error(msg || "Falha na importação.")}
        />
      )}
      {openInterestedModal && (
        <InterestedImportModal
          onClose={() => setOpenInterestedModal(false)}
          onSuccess={(msg: string) => {
            toast.success(msg || "Importado!");
            setOpenInterestedModal(false);
            setHasInterestedData(true);
          }}
          onError={(msg: string) => toast.error(msg || "Falha na importação.")}
        />
      )}
    </div>
  );
};

export default DataManagementPage;
