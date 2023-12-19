import { BsFillTrashFill } from "react-icons/bs";

function BtnTrash() {
  return (
    <div className="shadow-md  p-2  gap-5 disabled:bg-gray-400 bg-white  text-black font-bold   rounded-full inline-flex items-center disabled:shadow-none disabled:text-gray-200">
      <BsFillTrashFill className="text-md" />
    </div>
  );
}

export default BtnTrash;
