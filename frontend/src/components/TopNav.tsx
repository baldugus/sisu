import { Link, useLocation } from "react-router-dom";
import { Home, Users, UserCheck, Database, PhoneCall, BarChart3, HelpCircle } from "lucide-react";

const TopNav = () => {
  const location = useLocation();

  const getLinkClass = (path: string) => {
    const base = "flex items-center gap-2 text-sm font-medium transition-colors hover:text-primary";
    return location.pathname === path ? `${base} text-primary` : `${base} text-muted-foreground`;
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center px-4 max-w-full">
        <div className="mr-4 flex">
          <Link to="/" className="mr-6 flex items-center space-x-2 text-foreground">
            <span className="font-bold sm:inline-block">SISU Dashboard</span>
          </Link>
          <nav className="flex items-center space-x-6 text-sm font-medium">
            <Link to="/" className={getLinkClass("/")}>
              <Home className="h-4 w-4" /> Início
            </Link>
            <Link to="/data-management" className={getLinkClass("/data-management")}>
              <Database className="h-4 w-4" /> Dados
            </Link>
            <Link to="/approved-page" className={getLinkClass("/approved-page")}>
              <UserCheck className="h-4 w-4" /> Aprovados
            </Link>
            <Link to="/subscribe-page" className={getLinkClass("/subscribe-page")}>
              <Users className="h-4 w-4" /> Em espera
            </Link>
            <Link to="/calls-page" className={getLinkClass("/calls-page")}>
              <PhoneCall className="h-4 w-4" /> Chamadas
            </Link>
            <Link to="/reports-page" className={getLinkClass("/reports-page")}>
              <BarChart3 className="h-4 w-4" /> Relatórios
            </Link>
          </nav>
        </div>
        <div className="flex flex-1 items-center justify-end space-x-2">
          <button className="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9">
            <HelpCircle className="h-5 w-5" />
            <span className="sr-only">Ajuda</span>
          </button>
        </div>
      </div>
    </header>
  );
};

export default TopNav;
