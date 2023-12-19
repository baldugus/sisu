import { useParams } from "react-router-dom";
import { CallDataTable } from "../components";
import { Typography } from "@material-tailwind/react";
import { FetchApplicationsByRollCall } from "../../wailsjs/go/main/App";

const CallPage = () => {
  const { id } = useParams();
  const handleGetData = async () => {
    const response: any = id && FetchApplicationsByRollCall(parseInt(id));
    const { data } = await response;
    return await data;
  };

  return (
    <section className="px-4 w-full  h-[calc(100vh-82px)] overflow-auto lg:rounded-tl-3xl">
      <Typography variant="h2" className="border-b-2 border-black py-4">
        {`${id}ª Chamada`}
      </Typography>
      <CallDataTable handleGetData={handleGetData} />
    </section>
  );
};

export default CallPage;
