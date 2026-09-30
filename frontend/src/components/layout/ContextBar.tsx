import { useEffect, useState } from 'react';
import { CalendarDays, HelpCircle } from 'lucide-react';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { FetchApprovedSelection } from '@/lib/backend';

interface SelectionInfo {
  year: number;
}

export function ContextBar() {
  const [selection, setSelection] = useState<SelectionInfo | null>(null);

  useEffect(() => {
    FetchApprovedSelection()
      .then((sel) => {
        if (sel) {
          setSelection({ year: sel.Year });
        }
      })
      .catch(() => {});
  }, []);

  const QUOTAS = [
    { code: 'AC', label: 'Ampla Concorrência' },
    { code: 'C1', label: 'Escola Pública — Renda ≤ 1,5 SM' },
    { code: 'C2', label: 'Escola Pública — Renda ≤ 1,5 SM + PPI' },
    { code: 'C3', label: 'Escola Pública — PPI' },
  ];

  return (
    <header className="h-[54px] border-b border-border flex items-center justify-between px-6 shrink-0 bg-background">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        {selection ? (
          <>
            <CalendarDays className="size-4" />
            <span className="font-medium text-foreground font-mono">
              {selection.year}
            </span>
          </>
        ) : (
          <span className="italic">Nenhum ciclo importado</span>
        )}
      </div>

      <Tooltip>
        <TooltipTrigger asChild>
          <button
            aria-label="Referência de cotas"
            className="p-1.5 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <HelpCircle className="size-4" />
          </button>
        </TooltipTrigger>
        <TooltipContent side="left" className="max-w-xs">
          <p className="font-semibold mb-1 text-xs">Cotas SISU</p>
          <ul className="text-xs space-y-0.5">
            {QUOTAS.map(({ code, label }) => (
              <li key={code}>
                <span className="font-mono font-bold">{code}</span>
                {' — '}
                {label}
              </li>
            ))}
          </ul>
        </TooltipContent>
      </Tooltip>
    </header>
  );
}
