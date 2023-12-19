import { MdPeopleAlt } from "react-icons/md";
import { BsPersonFillCheck } from "react-icons/bs";
import { BiSolidReport } from "react-icons/bi";
import { GiBugleCall } from "react-icons/gi";
import { AiFillHome } from "react-icons/ai";
import { CiImport, CiExport } from "react-icons/ci";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
// import {
//   ImportApprovedApplications,
//   ListApprovedApplications,
// } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";

const SideBar = () => {
  const [active, setActive] = useState("home");

  const handleActive = (text: string) => {
    if (active == text) {
      return "text-white";
    }
    return "text-gray-700";
  };

  return (
    <aside className="py-2 px-5 flex flex-col justify-between w-1/4 xl:w-2/12 2xl:w-60 bg-black rounded-se-3xl ">
      <div className="flex flex-col justify-center ">
        <Link
          onClick={() => setActive("home")}
          to="/"
          className={`${handleActive(
            "home"
          )}  gap-5 cursor-pointer bg-grey-light hover:text-gray-500 font-bold py-2 px-4 rounded inline-flex items-center `}
        >
          <AiFillHome className="text-3xl" /> Início
        </Link>
        <Link
          onClick={() => setActive("subscribe")}
          to="/subscribe-page"
          className={`${handleActive(
            "subscribe"
          )}  gap-5 cursor-pointer bg-grey-light hover:text-gray-500 font-bold py-2 px-4 rounded inline-flex items-center `}
        >
          <MdPeopleAlt className="text-3xl" /> Em espera
        </Link>
        <Link
          onClick={() => setActive("approved")}
          to="/approved-page"
          className={`${handleActive(
            "approved"
          )} gap-5 cursor-pointer bg-grey-light hover:text-gray-500 font-bold py-2 px-4 rounded inline-flex items-center`}
        >
          <BsPersonFillCheck className="text-3xl" /> Aprovados
        </Link>
        <Link
          onClick={() => setActive("calls")}
          to="/calls-page"
          className={`${handleActive(
            "calls"
          )} gap-5 cursor-pointer bg-grey-light hover:text-gray-500 font-bold py-2 px-4 rounded inline-flex items-center`}
        >
          <GiBugleCall className="text-3xl" /> Chamadas
        </Link>
        <Link
          onClick={() => setActive("reports")}
          to="/reports-page"
          className={`${handleActive(
            "reports"
          )} gap-5 cursor-pointer bg-grey-light hover:text-gray-500 font-bold py-2 px-4 rounded inline-flex items-center`}
        >
          <BiSolidReport className="text-3xl" /> Relatórios
        </Link>
      </div>
      <div className="flex flex-col gap-2 px-4 py-4 border-t">
        <h1 className="text-sm text-center text-white">SISU 2023</h1>

        {/* <h1
          // onClick={async () =>
          //   console.log(
          //     ImportApprovedApplications().then((res) => console.log(res))
          //   )
          // }
          className="py-1 gap-5 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center"
        >
          <CiImport className="text-2xl" /> Importar
        </h1>
        <h1
          // onClick={() =>
          //   ListApprovedApplications().then((res) => alert(JSON.stringify(res)))
          // }
          className=" py-1 gap-5 bg-white cursor-pointer text-black font-bold px-4 rounded-full inline-flex items-center"
        >
          <CiExport className="text-2xl" /> Exportar
        </h1> */}
      </div>
    </aside>
  );
};

export default SideBar;
