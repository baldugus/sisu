import { Tab, Tabs, TabsHeader } from "@material-tailwind/react";

interface ITabFilter {
  dataFilter: {
    label: string;
    value: string;
  }[];
  setActiveTab: React.Dispatch<React.SetStateAction<string>>;
}

const TabsShow = ({ dataFilter, setActiveTab }: ITabFilter) => (
  <Tabs value="all" className="w-full md:w-max">
    <TabsHeader className="bg-gray-400">
      {dataFilter.map(({ label, value }: any) => (
        <Tab key={value} value={value} onClick={() => setActiveTab(value)}>
          &nbsp;&nbsp;{label}&nbsp;&nbsp;
        </Tab>
      ))}
    </TabsHeader>
  </Tabs>
);
export default TabsShow;
