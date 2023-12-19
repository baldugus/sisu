import { useEffect, useState } from "react";
import { GiBugleCall } from "react-icons/gi";
import { AiFillPlusCircle } from "react-icons/ai";
import {
  CloseRollCall,
  CreateRollCall,
  FetchRollCalls,
  OpenRollCall,
  DeleteRollcall,
} from "../../wailsjs/go/main/App";
import { Link } from "react-router-dom";
import { InfoCard } from "../components";
import { Alert, Button } from "@material-tailwind/react";
import { BsFillTrashFill } from "react-icons/bs";

const status: any = {
  CALLING: {
    color: "bg-yellow-400",
    name: "Em andamento",
  },
  DONE: {
    color: "bg-red-400",
    name: "Fechado",
  },
  ENROLLED: {
    color: "bg-green-400",
    name: "Matriculado(a)",
  },
};

const CallsPage = () => {
  const [rowsCalls, setRowsCalls] = useState<any>([]);
  // -----------------------------INICIO ALERTA-------------------------------
  const [open, setOpen] = useState(false);
  const [dialogContent, setDialogContent] = useState({
    color: "",
    msg: "",
  });
  // -----------------------------FIM ALERTA-------------------------------

  useEffect(() => {
    FetchRollCalls().then(
      (res) =>
        // setRowsTable(res.data.map((el: any) => el))
        // setRowsTable(res.data)
        setRowsCalls(res.data)
      // setRowsTable(res)
    );
  }, []);

  return (
    <section className="px-4 w-full  h-[calc(100vh-82px)] overflow-auto lg:rounded-tl-3xl">
      
      <h1 className="text-3xl font-bold border-b-2 border-black py-4">
        Chamadas
      </h1>

      {rowsCalls &&
        rowsCalls.map((el: any) => (
          <div className="flex gap-4 items-center" key={el.ID}>
            <Link to={`/call-page/${el.ID}`} className="w-full">
              <div
                key={el.ID}
                className="px-4 border justify-between hover:bg-blue-gray-50 border-blue-600 cursor-pointer shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 rounded-xl my-4"
              >
                <div className="flex gap-4 items-center">
                  <GiBugleCall />
                  {el.Number}ª Chamada
                </div>

                <div className="flex gap-4 items-center">
                  <div
                    className={`w-4 h-4 ${
                      status[el.Status].color
                    } rounded-full`}
                  />
                  {status[el.Status].name}
                </div>
              </div>
            </Link>
            <div className="flex h-10 gap-4">
              <div className="flex gap-2">
                <Button
                  variant="gradient"
                  color="green"
                  disabled={el.Status == "CALLING" ? true : false}
                  onClick={() =>
                    OpenRollCall(el.ID).then((res) => {
                      if (res.status != 200) {
                        setDialogContent({ color: "bg-red-400", msg: res.msg });
                        setOpen(true);
                      } else {
                        setDialogContent({
                          color: "bg-green-400",
                          msg: res.msg,
                        });
                        setOpen(true);
                      }
                      FetchRollCalls().then(
                        (res) =>
                          // setRowsTable(res.data.map((el: any) => el))
                          // setRowsTable(res.data)
                          setRowsCalls(res.data)
                        // setRowsTable(res)
                      );
                    })
                  }
                  className="w-48"
                >
                  Abrir Chamada
                </Button>
              </div>
              <Button
                variant="gradient"
                color="red"
                disabled={el.Status == "CALLING" ? false : true}
                className="w-48"
                onClick={() =>
                  CloseRollCall(el.ID).then((res) => {
                    if (res.status != 200) {
                      setDialogContent({ color: "bg-red-400", msg: res.msg });
                      setOpen(true);
                    } else {
                      setDialogContent({ color: "bg-green-400", msg: res.msg });
                      setOpen(true);
                    }
                    FetchRollCalls().then(
                      (res) =>
                        // setRowsTable(res.data.map((el: any) => el))
                        // setRowsTable(res.data)
                        setRowsCalls(res.data)
                      // setRowsTable(res)
                    );
                  })
                }
              >
                Fechar Chamada
              </Button>
            </div>
          </div>
        ))}
      {!rowsCalls && (
        <InfoCard
          text="Não há dados disponíveis, retorne à página inicial e importe o(s) .csv(s) de inscritos e/ou aprovados."
          type="sad"
        />
      )}

      {/* --------------------------------ALERTA------------------------------------------------ */}
      <div className="flex gap-2">
        <Button
          variant="gradient"
          className="w-48 flex gap-4 items-center"
          onClick={() => {
            CreateRollCall().then((res) => {
              if (res.status != 200) {
                setDialogContent({ color: "bg-red-400", msg: res.msg });
                setOpen(true);
              } else {
                setDialogContent({ color: "bg-green-400", msg: res.msg });
                setOpen(true);
              }
              FetchRollCalls().then(
                (res) =>
                  // setRowsTable(res.data.map((el: any) => el))
                  // setRowsTable(res.data)
                  setRowsCalls(res.data)
                // setRowsTable(res)
              );
            });
          }}
        >
          <AiFillPlusCircle /> Abrir Nova Chamada
        </Button>

        {rowsCalls.length > 1 && (
          <Button
            color="red"
            variant="gradient"
            className="w-48 flex gap-4 items-center "
            onClick={() => {
              DeleteRollcall(rowsCalls[rowsCalls.length - 1].ID).then((res) => {
                if (res.status != 200) {
                  setDialogContent({ color: "bg-red-400", msg: res.msg });
                  setOpen(true);
                } else {
                  setDialogContent({ color: "bg-green-400", msg: res.msg });
                  setOpen(true);
                }
                FetchRollCalls().then(
                  (res) =>
                    // setRowsTable(res.data.map((el: any) => el))
                    // setRowsTable(res.data)
                    setRowsCalls(res.data)
                  // setRowsTable(res)
                );
              });
            }}
          >
            <BsFillTrashFill /> Deletar{" "}
            {rowsCalls[rowsCalls.length - 1] &&
              rowsCalls[rowsCalls.length - 1].ID}
            ª Chamada
          </Button>
        )}
      </div>
      <Alert
        open={open}
        onClose={() => setOpen(false)}
        className={`${dialogContent.color} mt-4`}
        animate={{
          mount: { y: 0 },
          unmount: { y: 100 },
        }}
      >
        {dialogContent.msg}
      </Alert>
    </section>
  );
};

export default CallsPage;
