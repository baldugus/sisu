import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  Users,
  ListOrdered,
  FileText,
  Database,
} from 'lucide-react';
import { cn } from '@/lib/utils';

const NAV = [
  { to: '/',           label: 'Painel',      Icon: LayoutDashboard },
  { to: '/candidatos', label: 'Candidatos',  Icon: Users },
  { to: '/chamadas',   label: 'Chamadas',    Icon: ListOrdered },
  { to: '/relatorios', label: 'Relatórios',  Icon: FileText },
  { to: '/dados',      label: 'Dados',       Icon: Database },
] as const;

export function Sidebar() {
  return (
    <aside className="flex flex-col w-16 lg:w-56 h-full bg-sidebar text-sidebar-foreground shrink-0 rounded-se-3xl">
      {/* Wordmark */}
      <div className="px-4 py-6 hidden lg:block">
        <span className="font-heading font-black text-xl tracking-tight text-white uppercase">
          SISU
        </span>
      </div>
      <div className="px-4 py-6 lg:hidden flex justify-center">
        <span className="font-heading font-black text-lg text-white">S</span>
      </div>

      {/* Nav */}
      <nav className="flex flex-col gap-1 px-2 flex-1">
        {NAV.map(({ to, label, Icon }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-colors',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sidebar-ring',
                isActive
                  ? 'bg-sidebar-primary text-white'
                  : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-foreground'
              )
            }
          >
            <Icon className="size-5 shrink-0" />
            <span className="hidden lg:block">{label}</span>
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="px-4 py-4 hidden lg:block">
        <p className="text-xs text-sidebar-foreground/40 font-mono">
          SISU 2025
        </p>
      </div>
    </aside>
  );
}
