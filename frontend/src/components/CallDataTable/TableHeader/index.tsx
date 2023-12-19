import { Typography } from "@material-tailwind/react";
import React from "react";
import CheckBox from "../CheckBox";

const TableHeader = ({
  headings,
  selectedRows,
  setSelectedRows,
  filteredRows,
  hasSelect,
}: any) => {
  const isSelectedRowsEmpty = selectedRows.length > 0;
  const isFistHeader = (index: number) => index == 0;

  const handleAlternateSelection = (
    e: React.MouseEvent<HTMLInputElement, MouseEvent>
  ) => {
    e.stopPropagation();
    isSelectedRowsEmpty
      ? setSelectedRows([])
      : setSelectedRows([...filteredRows]);
  };

  const showSelectedNumber = isSelectedRowsEmpty && (
    <div className="bg-white p-2 rounded-full font-bold">
      {selectedRows.length}
    </div>
  );

  return (
    <thead>
      <tr>
        {headings.map((head: any, index: any) => (
          <th
            key={head}
            className="cursor-pointer border-y border-blue-gray-100 bg-blue-gray-50/50 p-4 transition-colors hover:bg-blue-gray-50"
          >
            <div className="flex">
              <Typography
                variant="small"
                color="blue-gray"
                className="flex items-center justify-between gap-2 font-normal leading-none opacity-70"
              >
                {head}{" "}
              </Typography>
              {isFistHeader(index) && hasSelect && (
                <div className="inline-flex items-center ">
                  <CheckBox
                    checked={isSelectedRowsEmpty}
                    onClick={handleAlternateSelection}
                    id="SelectorAll"
                  />
                  {showSelectedNumber}
                </div>
              )}
            </div>
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
