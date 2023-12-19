import { Route, Routes } from "react-router-dom";
import { NavBar, SideBar } from "./components";
import {
  ApprovedPage,
  CallPage,
  CallsPage,
  HomePage,
  ReportsPage,
  SubscribePage,
} from "./pages";

function App() {
  return (
    <div className="bg-white flex flex-col h-screen">
      <NavBar />
      <main className="flex h-screen">
        <SideBar />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/subscribe-page" element={<SubscribePage />} />
          <Route path="/approved-page" element={<ApprovedPage />} />
          <Route path="/reports-page" element={<ReportsPage />} />
          <Route path="/calls-page" element={<CallsPage />} />
          <Route path="/call-page/:id" element={<CallPage/>}/>
        </Routes>
      </main>
    </div>
  );
}

export default App;
