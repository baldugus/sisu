import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Users, UserCheck, GraduationCap, PhoneOutgoing } from "lucide-react";

const mockTableData = [
  {
    turno: "Matutino",
    cota: "AC",
    vagas: 20,
    matriculados: 15,
    restantes: 5,
    espera: 320,
  },
  {
    turno: "Matutino",
    cota: "C1",
    vagas: 5,
    matriculados: 5,
    restantes: 0,
    espera: 85,
  },
  {
    turno: "Noturno",
    cota: "AC",
    vagas: 20,
    matriculados: 12,
    restantes: 8,
    espera: 410,
  },
  {
    turno: "Noturno",
    cota: "L2",
    vagas: 10,
    matriculados: 8,
    restantes: 2,
    espera: 150,
  },
];

export default function DashboardPage() {
  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight text-foreground">Dashboard</h2>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Aprovados
            </CardTitle>
            <UserCheck className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">92</div>
            <p className="text-xs text-muted-foreground">
              candidatos aprovados
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total em Espera
            </CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1,540</div>
            <p className="text-xs text-muted-foreground">
              na fila de espera geral
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Vagas Preenchidas
            </CardTitle>
            <GraduationCap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">45 / 92</div>
            <Progress value={48.9} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              48.9% das vagas totais
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Status da Chamada
            </CardTitle>
            <PhoneOutgoing className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">Chamada 2</div>
            <p className="text-xs text-muted-foreground">
              Ativa no momento
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-1 mt-4">
        <Card className="col-span-1">
          <CardHeader>
            <CardTitle>Resumo por Cota e Turno</CardTitle>
            <CardDescription>
              Detalhes do preenchimento de vagas para o curso de Análise e Desenvolvimento de Sistemas.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Turno</TableHead>
                  <TableHead>Cota</TableHead>
                  <TableHead className="text-right">Vagas</TableHead>
                  <TableHead className="text-right">Matriculados</TableHead>
                  <TableHead className="text-right">Restantes</TableHead>
                  <TableHead className="text-right">Fila de Espera</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mockTableData.map((row, idx) => (
                  <TableRow key={idx}>
                    <TableCell className="font-medium">{row.turno}</TableCell>
                    <TableCell>{row.cota}</TableCell>
                    <TableCell className="text-right">{row.vagas}</TableCell>
                    <TableCell className="text-right">{row.matriculados}</TableCell>
                    <TableCell className="text-right">{row.restantes}</TableCell>
                    <TableCell className="text-right">{row.espera}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
