import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { Session } from "./pages/Session";
import Header from "./components/Header";

export function AppRoutes() {
  return (
    <BrowserRouter>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/session" element={<Session />} />
        <Route path="/session/:sessionId" element={<Session />} />
      </Routes>
    </BrowserRouter>
  );
}
