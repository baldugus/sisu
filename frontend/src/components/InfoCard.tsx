import { AiOutlineInfo } from "react-icons/ai";
import { CiWarning } from "react-icons/ci";
import { FaRegSadCry } from "react-icons/fa";

interface IInfoCard {
  type: "info" | "warning" | "sad";
  text: any;
}

const InfoCard = ({ type, text }: IInfoCard) => {
  const handleIcon = () => {
    switch (type) {
      case "info":
        return <AiOutlineInfo className="animate-ping text-3xl text-blue-600" />;
      case "warning":
        return <CiWarning className="animate-ping text-3xl text-orange-600" />;
      case "sad":
        return <FaRegSadCry className="animate-ping text-3xl text-yellow-600" />;
      default:
        return "";
    }
  };

  return (
    <div className="shadow-lg flex gap-2 text-1xl  py-4 items-center bg-gray-100 px-2 rounded-xl my-4">
      {handleIcon()}
      {text}
    </div>
  );
};

export default InfoCard;
