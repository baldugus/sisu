import { Route, Routes } from 'react-router-dom';
import { Toaster } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';
import { Sidebar, ContextBar } from '@/components/layout';
import { lazy, Suspense } from 'react';

const Painel     = lazy(() => import('@/pages/Painel'));
const Candidatos = lazy(() => import('@/pages/Candidatos'));
const Chamadas   = lazy(() => import('@/pages/Chamadas'));
const Chamada    = lazy(() => import('@/pages/Chamada'));
const Relatorios = lazy(() => import('@/pages/Relatorios'));
const Dados      = lazy(() => import('@/pages/Dados'));

function App() {
  return (
    <TooltipProvider delayDuration={400}>
      <div className="flex h-screen bg-background overflow-hidden">
        <Sidebar />

        <div className="flex flex-col flex-1 min-w-0">
          <ContextBar />

          <main className="flex-1 overflow-auto">
            <Suspense fallback={null}>
              <Routes>
                <Route path="/"               element={<Painel />} />
                <Route path="/candidatos"     element={<Candidatos />} />
                <Route path="/chamadas"       element={<Chamadas />} />
                <Route path="/chamadas/:id"   element={<Chamada />} />
                <Route path="/relatorios"     element={<Relatorios />} />
                <Route path="/dados"          element={<Dados />} />
              </Routes>
            </Suspense>
          </main>
        </div>
      </div>

      <Toaster position="top-right" richColors closeButton />
    </TooltipProvider>
  );
}

export default App;
