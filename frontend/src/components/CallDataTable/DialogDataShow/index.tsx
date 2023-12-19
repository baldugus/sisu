import {
  Card,
  CardHeader,
  Input,
  Typography,
  Button,
  CardBody,
  Chip,
  CardFooter,
  Tabs,
  TabsHeader,
  Tab,
  Avatar,
  IconButton,
  Tooltip,
  Dialog,
  DialogHeader,
  DialogBody,
  DialogFooter,
  Select,
  Option,
} from "@material-tailwind/react";
import { useState } from "react";
import {
  AbsentApplication,
  ClearApplicationStatus,
  EnrollApplication,
  FetchApplicationsByRollCall,
} from "../../../../wailsjs/go/main/App";

const DialogDataShow = ({
  open,
  handleOpen,
  filteredData,
  currentDataIdShow,
  setRowsTable,
  setDataCurrentStatus,
  dataCurrentStatus,
  handleGetData,
  hasSelector,
}: any) => {
  const handleAlterStatus = () => {
    switch (dataCurrentStatus) {
      case "APPROVED":
        ClearApplicationStatus(filteredData[currentDataIdShow].ID).then(() =>
          handleGetData().then((res: any) => setRowsTable(res))
        );
        break;
      case "ABSENT":
        AbsentApplication(filteredData[currentDataIdShow].ID).then(() =>
          handleGetData().then((res: any) => setRowsTable(res))
        );

        break;
      case "ENROLLED":
        EnrollApplication(filteredData[currentDataIdShow].ID).then(() =>
          handleGetData().then((res: any) => setRowsTable(res))
        );

        break;
      default:
        alert("Nenhum status selecionado");
    }
  };

  return (
    <Dialog open={open} handler={handleOpen}>
      {filteredData[currentDataIdShow] && (
        <DialogHeader>{filteredData[currentDataIdShow].Name}</DialogHeader>
      )}
      {filteredData[currentDataIdShow] && (
        <DialogBody className="h-[22rem] overflow-scroll">
          {Object.entries(filteredData[currentDataIdShow]).map(
            ([key, value]) => (
              <div key={key} className="flex gap-2 border-b-2">
                <Typography variant="small" className="font-normal">
                  {key}:
                </Typography>
                <Typography variant="small">{value}</Typography>
              </div>
            )
          )}
        </DialogBody>
      )}
      <DialogFooter className="space-x-2 flex flex-col justify-around">
        <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
          {hasSelector && "Alterar Status:"}
          {filteredData[currentDataIdShow] &&
            filteredData[currentDataIdShow].Status}
        </label>
        {hasSelector && filteredData[currentDataIdShow] && (
          <select
            id="countries"
            className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-1/2 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            onChange={(e: any) => setDataCurrentStatus(e.target.value)}
            defaultValue={filteredData[currentDataIdShow].Status}
          >
            <option value="APPROVED">Convocado</option>
            <option value="ENROLLED">Matriculado</option>
            <option value="ABSENT">Faltoso</option>
          </select>
        )}
        <div className="mt-4 gap-4 flex">
          <Button variant="text" color="blue-gray" onClick={handleOpen}>
            {hasSelector ? "Cancelar" : "Voltar"}
          </Button>
          {hasSelector && (
            <Button
              variant="gradient"
              color="green"
              onClick={() => {
                handleOpen();
                handleAlterStatus();
              }}
              disabled={
                filteredData[currentDataIdShow] &&
                filteredData[currentDataIdShow].Status == dataCurrentStatus
                  ? true
                  : false
              }
            >
              Confirmar {dataCurrentStatus}
            </Button>
          )}
        </div>
      </DialogFooter>
    </Dialog>
  );
};

export default DialogDataShow;
