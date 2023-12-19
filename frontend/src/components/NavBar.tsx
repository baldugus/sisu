import { useState } from "react";
import logo from "./assets/images/logo-universal.png";

import { AiFillSetting } from "react-icons/ai";
import { IoIosHelpCircle } from "react-icons/io";
import {
  Drawer,
  Button,
  Typography,
  IconButton,
} from "@material-tailwind/react";

const NavBar = () => {
  const [openRight, setOpenRight] = useState(false);

  const openDrawerRight = () => setOpenRight(true);
  const closeDrawerRight = () => setOpenRight(false);

  return (
    <>
      <nav className="flex items-center justify-between px-10">
        <h1 className="text-4xl font-bold py-5">SISU</h1>
        <ul className="flex gap-4">
          <IconButton
            variant="text"
            className="text-black hover:bg-gray-50 active:bg-gray-100 "
            onClick={openDrawerRight}
          >
            <IoIosHelpCircle className="text-2xl cursor-pointer text-center" />
          </IconButton>
        </ul>
      </nav>
      <main>
        <>
          <Drawer
            placement="right"
            open={openRight}
            onClose={closeDrawerRight}
            className="p-4"
            size={500}
          >
            <div className="mb-6 flex items-center justify-between">
              <Typography variant="h5" color="blue-gray">
                Ajuda
              </Typography>

              <IconButton
                variant="text"
                className="text-black hover:bg-gray-50 active:bg-gray-100"
                onClick={closeDrawerRight}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className="h-5 w-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </IconButton>
            </div>
            <div>
              <Typography variant="h5" color="blue-gray">
                Cotas
              </Typography>
              <Typography variant="h6" color="blue-gray">
                AC:
              </Typography>
              <Typography variant="p" color="blue-gray">
                Ampla concorrência
              </Typography>
              <Typography variant="h6" color="blue-gray">
                C1:
              </Typography>
              <Typography variant="p" color="blue-gray">
                Candidatos Negros ou Indígenas com comprovação de carência.
                socioeconômica
              </Typography>
              <Typography variant="h6" color="blue-gray">
                C2:
              </Typography>
              <Typography variant="p" color="blue-gray">
                Candidatos com deficiência ou filhos de policiais militares,
                bombeiros militares, inspetores de segurança e administração
                penitenciária, mortos ou incapacitados em razão do serviço, com
                comprovação de carência socioeconômica.
              </Typography>
              <Typography variant="h6" color="blue-gray">
                C3:
              </Typography>
              <Typography variant="p" color="blue-gray">
                Candidatos que tenham cursado na rede pública os últimos quatro
                anos do ensino fundamental e todo o ensino médio e com
                comprovação de carência socioeconômica.
              </Typography>
            </div>
          </Drawer>
        </>
      </main>
    </>
  );
};

export default NavBar;
