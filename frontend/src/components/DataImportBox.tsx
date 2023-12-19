import { MdPeopleAlt } from "react-icons/md";
import { CiImport, CiExport } from "react-icons/ci";
import { BsFillTrashFill } from "react-icons/bs";
import CardWrapper from "./CardWrapper";

interface IDataImportBox {
  dataType: "inscritos" | "aprovados";
  text: string;
  mockHasData: boolean;
  actionfunction: any;
}

const DataImportBox = ({ dataType, text, mockHasData, actionfunction }: IDataImportBox) => {
  return (
    <CardWrapper>
        <div className="flex gap-2">
          <span className="bg-gray-300 p-1 rounded-full text-white ">
            <MdPeopleAlt className="text-2xl" />
          </span>
          <h4 className="text-lg font-bold">{text}</h4>
        </div>
        <div className="flex mt-4  gap-2 ">
          <button
            disabled={mockHasData}
            onClick={actionfunction}
            className="shadow-md  py-1 gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold  px-4 rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
          >
            <CiImport className="text-2xl" /> Importar
          </button>
          <button
            disabled={!mockHasData}
            className="shadow-md  px-2  gap-5 disabled:bg-gray-400 bg-white cursor-pointer text-black font-bold   rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200"
          >
            <BsFillTrashFill className="text-md" />
          </button>
        </div>
    </CardWrapper>
  );
};

export default DataImportBox;
