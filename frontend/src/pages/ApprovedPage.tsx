import { FetchApprovedSelection } from "../../wailsjs/go/main/App";
import { CallDataTable } from "../components";

const ApprovedPage = () => {
  const handleGetData = async () => {
    const response: any = FetchApprovedSelection();

    const { data } = await response;
    const { Applications } = await data;
    return await Applications;
  };

  return (
    <section className="px-4 w-full  h-[calc(100vh-82px)] overflow-auto lg:rounded-tl-3xl">
      <h1 className="text-3xl font-bold border-b-2 border-black py-4">
        Aprovados
      </h1>

      <CallDataTable handleGetData={handleGetData} hasSelector={false} />
    </section>
  );
};

export default ApprovedPage;
