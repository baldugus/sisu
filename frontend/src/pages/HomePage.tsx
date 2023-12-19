import { MdPeopleAlt } from "react-icons/md";
import { BtnTrash, CardWrapper, DataImportBox, InfoCard } from "../components";
import {
  Backup,
  Destroy,
  ExportCSV,
  LoadApprovedSelection,
  LoadInterestedSelection,
  Restore,
} from "../../wailsjs/go/main/App";
import { CiImport,CiExport } from "react-icons/ci";
import { BsFillTrashFill } from "react-icons/bs";


const HomePage = () => {
  return (
    <section className="px-4 w-full  h-[calc(100vh-82px)] overflow-auto lg:rounded-tl-3xl">
      <h1 className="text-3xl font-bold border-b-2 border-black py-4">
        Início
      </h1>
      <InfoCard
        text="Como passo inicial, adicione os CSVs, aos respectivos campos abaixo."
        type="info"
      />
      <div className=" flex gap-4">
        <DataImportBox
          actionfunction={LoadApprovedSelection}
          dataType="aprovados"
          text="Aprovados"
          mockHasData={false}
        />
        <DataImportBox
          actionfunction={LoadInterestedSelection}
          dataType="inscritos"
          text="Em Espera"
          mockHasData={false}
        />
      </div>
      <InfoCard
        text={
          <>
            Atenção! O botão <BtnTrash /> apagará os dados armazenados em sua
            máquina, utilize-o caso não necessite mais dos dados armazenados.
          </>
        }
        type="warning"
      />
      <div className=" flex gap-4">
        {/* <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">Ampla concorrência</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L1</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L2</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L5</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L6</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L9</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L10</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L13</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <MdPeopleAlt className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">L14</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <div
              // disabled={mockHasData}
              className=" py-1 gap-5 disabled:bg-gray-400  text-black font-bold  px-4  inline-flex items-center  disabled:text-gray-200"
            >
              20 - vagas
            </div>
          </div>
        </CardWrapper> */}
      </div>
      <InfoCard text={<>Banco de dados & CSV</>} type="info" />
      <div className=" flex gap-4">
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <CiImport className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">Importar banco de dados:</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <button
              onClick={Restore}
              className="shadow-md  py-1 gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
            >
              <CiImport className="text-md" /> Importar
            </button>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <CiExport className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">Exportar banco de dados:</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <button
              onClick={Backup}
              className="shadow-md  py-1 gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
            >
              <CiExport className="text-md" /> Exportar
            </button>
          </div>
        </CardWrapper>
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <CiExport className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">Exportar CSV de matriculados:</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <button
              onClick={ExportCSV}
              className="shadow-md  py-1 gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
            >
              <CiExport className="text-md" /> Exportar
            </button>
          </div>
        </CardWrapper>
      </div>
      
      <InfoCard
        text={
          <>
            Atenção! O botão Abaixo apagará os dados armazenados em sua máquina,
            utilize-o caso não necessite mais dos dados armazenados.
          </>
        }
        type="warning"
      />
      <div className=" flex gap-4">
        <CardWrapper>
          <div className="flex gap-2">
            <span className="bg-gray-300 p-1 rounded-full text-white ">
              <BsFillTrashFill className="text-2xl" />
            </span>
            <h4 className="text-lg font-bold">DESTRUIR TODOS OS DADOS:</h4>
          </div>
          <div className="flex mt-4  gap-2 ">
            <button
              onClick={Destroy}
              className="shadow-md  py-1 gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
            >
              <BsFillTrashFill className="text-md" /> Destruir
            </button>
          </div>
        </CardWrapper>
      </div>

      {/* <InfoCard text="Não há dados disponíveis, retorne à página inicial e importe o .csv de aprovados." type="sad"/> */}
    </section>
  );
};

export default HomePage;
