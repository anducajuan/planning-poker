import { ThemeProvider, CssBaseline } from "@mui/material";
import { theme } from "./theme/theme";
import { AppRoutes } from "./routes";
import { ToastContainer, Bounce } from "react-toastify";

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <ToastContainer
        position="bottom-right"
        autoClose={3000}
        newestOnTop={true}
        hideProgressBar={true}
        theme="colored"
        transition={Bounce}
      />
      <AppRoutes />
    </ThemeProvider>
  );
}
