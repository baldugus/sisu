import { useEffect, useState } from "react";
import {
  EmailPDF,
  EnrollmentPDF,
  FetchPeriods,
  FetchRollCalls,
  TeacherPDF,
  WebsitePDF,
} from "../../wailsjs/go/main/App";
import Divider from "../components/CallDataTable/Divider";

interface IPeriods {
  ID: number;
  Period: "Matutino" | "Noturno";
}
interface ICalls {
  ID: number;
  Status: string;
  Number: number;
}

const ReportsPage = () => {
  const [periods, setPeriods] = useState<Array<IPeriods>>();
  const [rowsCalls, setRowsCalls] = useState<Array<ICalls>>();
  useEffect(() => {
    FetchPeriods().then((res) => setPeriods(res.data));
    FetchRollCalls().then((res) => setRowsCalls(res.data));
  }, []);

  const handleExportWebsite = (callId: number, periodId: number) => {
    WebsitePDF(callId, periodId);
  };
  const handleEnrollmentPDF = (callId: number, periodId: number) => {
    EnrollmentPDF(callId, periodId);
  };
  const handleTeacherPDF = (callId: number) => {
    TeacherPDF(callId);
  };
  const handleEmailPDF = (callId: number, periodId: number) => {
    EmailPDF(callId, periodId);
  };

  return (
    <section className="px-4 w-full  h-[calc(100vh-82px)] overflow-auto lg:rounded-tl-3xl">
      <h1 className="text-3xl font-bold border-b-2 border-black py-4">
        Relatórios
      </h1>
      <Divider>Exportar para o Website</Divider>
      {rowsCalls &&
        rowsCalls.map(
          (callElement) =>
            periods &&
            periods.map((periodElement) => (
              <div
                key={periodElement.ID}
                className="px-4 border justify-between hover:bg-blue-gray-50 border-blue-600 cursor-pointer shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 rounded-xl my-4"
                onClick={() =>
                  handleExportWebsite(callElement.ID, periodElement.ID)
                }
              >
                <div className="flex gap-4 items-center">
                  Exportar Arquivo Website {periodElement.Period}{" "}
                  {callElement.ID}ª Chamada
                </div>
                <div className="flex gap-4 items-center"></div>
              </div>
            ))
        )}
      <Divider>Exportar Convocados</Divider>
      {rowsCalls &&
        rowsCalls.map(
          (callElement) =>
            periods &&
            periods.map((periodElement) => (
              <div
                key={periodElement.ID}
                className="px-4 border justify-between hover:bg-blue-gray-50 border-blue-600 cursor-pointer shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 rounded-xl my-4"
                onClick={() =>
                  handleEnrollmentPDF(callElement.ID, periodElement.ID)
                }
              >
                <div className="flex gap-4 items-center">
                  Exportar Arquivo Convocados {periodElement.Period}{" "}
                  {callElement.ID}ª Chamada
                </div>
                <div className="flex gap-4 items-center"></div>
              </div>
            ))
        )}
      <Divider>Exportar Emails</Divider>
      {rowsCalls &&
        rowsCalls.map(
          (callElement) =>
            periods &&
            periods.map((periodElement) => (
              <div
                key={periodElement.ID}
                className="px-4 border justify-between hover:bg-blue-gray-50 border-blue-600 cursor-pointer shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 rounded-xl my-4"
                onClick={() =>
                  handleEmailPDF(callElement.ID, periodElement.ID)
                }
              >
                <div className="flex gap-4 items-center">
                  Exportar Arquivo Emails {periodElement.Period}{" "}
                  {callElement.ID}ª Chamada
                </div>
                <div className="flex gap-4 items-center"></div>
              </div>
            ))
        )}
      <Divider>Exportar Pdf do Professor</Divider>
      {periods &&
        periods.map((callElement) => (
          <div
            key={callElement.ID}
            className="px-4 border justify-between hover:bg-blue-gray-50 border-blue-600 cursor-pointer shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 rounded-xl my-4"
            onClick={() => handleTeacherPDF(callElement.ID)}
          >
            <div className="flex gap-4 items-center">
              Exportar Arquivo Convocados {callElement.Period} Turno
            </div>
            <div className="flex gap-4 items-center"></div>
          </div>
        ))}
    </section>
  );
};

export default ReportsPage;
