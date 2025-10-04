import { ThemeProvider, CssBaseline } from "@mui/material";
import { theme } from "./theme/theme";
import Header from "./components/Header";
import { AppRoutes } from "./routes";

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Header />
      <AppRoutes />
    </ThemeProvider>
  );
}
