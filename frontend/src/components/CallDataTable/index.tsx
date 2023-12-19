import { AiFillCloseCircle, AiOutlineSearch } from "react-icons/ai";
import { RxAvatar } from "react-icons/rx";
import {
  CardHeader,
  Input,
  Typography,
  CardBody,
} from "@material-tailwind/react";

import { useEffect, useState } from "react";
import {
  AbsentApplication,
  ClearApplicationStatus,
  EnrollApplication,
  FetchApplicationsByRollCall,
} from "../../../wailsjs/go/main/App";
import InfoCard from "../InfoCard";
import TableHeader from "./TableHeader";
import CheckBox from "./CheckBox";
import TabsShow from "./TabsShow";
import Divider from "./Divider";
import DialogDataShow from "./DialogDataShow";

const TABS = [
  {
    label: "Todos",
    value: "all",
  },
  {
    label: "Matutino",
    value: "Matutino",
  },
  {
    label: "Noturno",
    value: "Noturno",
  },
];

const TABS_2 = [
  {
    label: "Todos",
    value: "all",
  },
  {
    label: "Convocado(a)",
    value: "APPROVED",
  },
  {
    label: "Faltoso(a)",
    value: "ABSENT",
  },
  {
    label: "Matriculado(a)",
    value: "ENROLLED",
  },
];
const TABS_3 = [
  {
    label: "Todos",
    value: "all",
  },
  {
    label: "AC",
    value: "Ampla concorrência",
  },
  {
    label: "C1",
    value:
      "Candidatos Negros ou Indígenas com comprovação de carência socioeconômica",
  },
  {
    label: "C2",
    value:
      "Candidatos com deficiência ou filhos de policiais militares, bombeiros militares, inspetores de segurança e administração penitenciária, mortos ou incapacitados em razão do serviço, com comprovação de carência socioeconômica",
  },
  {
    label: "C3",
    value:
      "Candidatos que tenham cursado na rede pública os últimos quatro anos do ensino fundamental e todo o ensino médio e com comprovação de carência socioeconômica",
  },
];

const status: any = {
  APPROVED: {
    color: "bg-yellow-400",
    name: "Convocado(a)",
  },
  ABSENT: {
    color: "bg-red-400",
    name: "Faltoso(a)",
  },
  ENROLLED: {
    color: "bg-green-400",
    name: "Matriculado(a)",
  },
  WAITING: {
    color: "bg-cyan-400",
    name: "Esperando...",
  },
};

const TABLE_HEAD = ["Selecionar: ", "Inscrito", "Turno", "Cota", "Status", ""];

interface ITableProps {
  handleGetData: () => Promise<any>;
  hasSelector?: boolean;
}

const CallDataTable = ({ handleGetData, hasSelector = true }: ITableProps) => {
  const [rowsTable, setRowsTable] = useState<any>([]);
  const [activeTab, setActiveTab] = useState("all");
  const [activeTab2, setActiveTab2] = useState("all");
  const [activeTab3, setActiveTab3] = useState("all");
  const [currentIdDataShow, setCurrentIdData] = useState(1);
  const [searchInput, setSearchInput] = useState<any>("");
  // Estados dos status
  const [dataCurrentStatus, setDataCurrentStatus] = useState("");

  //   ----------Caixa de dialogo--------------
  const [open, setOpen] = useState(false);

  const handleOpen = () => setOpen(!open);
  //   ----------Caixa de dialogo--------------
  useEffect(() => {
    handleGetData().then((res: any) => setRowsTable(res));
  }, []);

  const filteredRows =
    rowsTable &&
    rowsTable.filter((el: any) => {
      if (
        activeTab === "all" &&
        activeTab2 === "all" &&
        activeTab3 === "all" &&
        searchInput === ""
      ) {
        return true; // Show all records when both tabs are "Todos"
      } else {
        // Filter based on both tabs' values
        return (
          (activeTab === "all" || el.Period === activeTab) &&
          (activeTab2 === "all" || el.Status === activeTab2) &&
          (activeTab3 === "all" || el.Quota === activeTab3) &&
          (searchInput === "" ||
            el.Name.toLowerCase().includes(searchInput.toLowerCase()) ||
            el.CPF.toLowerCase().includes(searchInput.toLowerCase()))
        );
      }
    });
  //-----------------------------INICIO LOGICA SELECTED ROWS----------------------------------------
  const [selectedRows, setSelectedRows] = useState<any>([]);
  const [selectedOption, setselectedOption] = useState<any>("APPROVED");

  //Seleciona os clicados
  const handleSelectedRow = (index: any) => {
    const selectedRow = filteredRows[index];
    const isSelected = selectedRows.includes(selectedRow);
    if (isSelected) {
      // If already selected, remove it from the list
      setSelectedRows(selectedRows.filter((row: any) => row !== selectedRow));
    } else {
      // If not selected, add it to the list
      setSelectedRows([...selectedRows, selectedRow]);
    }
  };

  //-----------------------------FIM LOGICA SELECTED ROWS----------------------------------------
  useEffect(() => {
    console.log(selectedRows);
  }, [selectedRows]);

  // const delay = (ms: any) => new Promise((resolve) => setTimeout(resolve, ms));

  const handleChangeMultiple = async (valueKey: any) => {
    switch (valueKey) {
      case "APPROVED":
        for (const el of selectedRows) {
          await ClearApplicationStatus(el.ID);
        }
        break;
      case "ABSENT":
        for (const el of selectedRows) {
          await AbsentApplication(el.ID);
        }
        break;
      case "ENROLLED":
        for (const el of selectedRows) {
          await EnrollApplication(el.ID);
        }
        break;
      default:
        alert("No change has been applied");
        break;
    }

    // The rest of your code remains unchanged
    handleGetData().then((res: any) => {
      setRowsTable(res);
      setSelectedRows([]);
      setselectedOption("APPROVED");
    });
  };

  return (
    <div>
      <CardHeader floated={false} shadow={false} className="rounded-none">
        <div className="mb-8 flex items-center justify-between gap-8"></div>
        <Divider>Visualização</Divider>
        <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
          <div className="flex gap-4">
            <TabsShow dataFilter={TABS} setActiveTab={setActiveTab} />
            <TabsShow dataFilter={TABS_3} setActiveTab={setActiveTab3} />
          </div>
        </div>
        <div className="flex flex-col items-center  gap-4 md:flex-row mt-4">
          <div className="flex gap-4">
            <TabsShow dataFilter={TABS_2} setActiveTab={setActiveTab2} />
          </div>

          <div className="w-full md:w-72">
            <Input
              label="Procurar por NOME ou CPF"
              id="Procurar por NOME ou CPF"
              name="Procurar por NOME ou CPF"
              icon={
                searchInput ? (
                  <button onClick={() => setSearchInput("")}>
                    <AiFillCloseCircle className="h-5 w-5" />
                  </button>
                ) : (
                  <AiOutlineSearch className="h-5 w-5" />
                )
              }
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
        </div>

        {/* -------------------------INICIO AREA DE SELECTION----------------------------------- */}
        {hasSelector && selectedRows.length > 0 && (
          <div className="mt-4 ">
            <div className="relative flex py-5 items-center">
              <div className="flex-grow border-t border-gray-400"></div>
              <span className="flex-shrink mx-4 text-gray-400">
                Multipla Edição
              </span>
              <div className="flex-grow border-t border-gray-400"></div>
            </div>
            <label htmlFor="select">Modificar: </label>
            <select
              name="select"
              className="w-40 h-10 border rounded-md"
              onChange={(e) => setselectedOption(e.target.value)}
            >
              {TABS_2.map((val: any, index) =>
                index != 0 ? (
                  <option value={`${val.value}`}>{val.label}</option>
                ) : (
                  ""
                )
              )}
            </select>
            <button
              className="ms-4 border p-2 rounded-full hover:bg-blue-200 bg-blue-400 text-white"
              onClick={() => handleChangeMultiple(selectedOption)}
            >
              Apply Selection to {selectedRows.length} Applicants
            </button>
          </div>
        )}
        {/* -------------------------FIM AREA DE SELECTION----------------------------------- */}
      </CardHeader>
      <CardBody className="overflow-scroll px-0">
        <table className="mt-4 w-full min-w-max table-auto text-left">
          <TableHeader
            headings={TABLE_HEAD}
            selectedRows={selectedRows}
            setSelectedRows={setSelectedRows}
            filteredRows={filteredRows}
            hasSelect={hasSelector}
          />

          <tbody>
            {rowsTable &&
              filteredRows.map((el: any, index: any) => {
                const isLast = index === rowsTable.length - 1;

                const classes = isLast
                  ? "p-4 opacity-100 transition-opacity transition-transform duration-300"
                  : "p-4 border-b border-blue-gray-50 opacity-100 transition-opacity transition-transform duration-300 delay-100";

                return (
                  <tr
                    key={el.CPF}
                    className={`fade-in ${classes} hover:bg-blue-gray-100 cursor-pointer ${
                      selectedRows.includes(filteredRows[index])
                        ? "bg-cyan-50"
                        : ""
                    }`}
                    onClick={() => {
                      handleOpen();
                      setCurrentIdData(index);
                      setDataCurrentStatus(el.Status);
                    }}
                  >
                    {/* -----------------------------------INICIO Coluna de selection------------------------------------- */}
                    <td className={classes}>
                      <div className="inline-flex items-center ">
                        {hasSelector && (
                          <CheckBox
                            onClick={(e) => {
                              e.stopPropagation();
                              handleSelectedRow(index);
                            }}
                            checked={selectedRows.includes(filteredRows[index])}
                            id={`checkbox-${el.Name}`}
                          />
                        )}
                      </div>
                    </td>
                    {/* -----------------------------------FIM Coluna de selection------------------------------------- */}
                    <td className={classes}>
                      <div className="flex items-center gap-3">
                        <RxAvatar size="40" />
                        <div className="flex flex-col">
                          <Typography
                            variant="small"
                            color="blue-gray"
                            className="font-normal"
                          >
                            {el.Name}
                          </Typography>
                        </div>
                      </div>
                    </td>
                    <td className={classes}>
                      <div className="flex flex-col">
                        <Typography
                          variant="small"
                          color="blue-gray"
                          className="font-normal"
                        >
                          {el.Period}
                        </Typography>
                      </div>
                    </td>

                    <td className={classes}>
                      <div className="flex flex-col">
                        <Typography
                          variant="small"
                          color="blue-gray"
                          className="font-normal w-10 truncate md:text-clip"
                        >
                          {el.Quota}
                        </Typography>
                      </div>
                    </td>
                    <td className={classes}>
                      <div className="flex flex-col">
                        <Typography
                          variant="small"
                          color="blue-gray"
                          className="font-normal flex gap-4"
                        >
                          <span
                            className={`w-4 h-4 rounded-full ${
                              status[el.Status].color
                            }`}
                          ></span>
                          {status[el.Status].name}
                        </Typography>
                      </div>
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>
      </CardBody>
      {/* ----------------------------------Parte da caixa de dialogo------------------------------------------------ */}
      {filteredRows && (
        <DialogDataShow
          open={open}
          handleOpen={handleOpen}
          filteredData={filteredRows}
          currentDataIdShow={currentIdDataShow}
          setRowsTable={setRowsTable}
          handleGetData={handleGetData}
          setDataCurrentStatus={setDataCurrentStatus}
          dataCurrentStatus={dataCurrentStatus}
          hasSelector={hasSelector}
        />
      )}

      {!rowsTable && (
        <InfoCard
          text="Não há dados disponíveis, retorne à página inicial ."
          type="sad"
        />
      )}
    </div>
  );
};

export default CallDataTable;
