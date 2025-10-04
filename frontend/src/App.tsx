import { ThemeProvider, CssBaseline } from "@mui/material";
import { StylesProvider } from "@mui/styles";
import { theme } from "./theme/theme";
import Header from "./components/Header";
import { AppRoutes } from "./routes";

export default function App() {
  return (
    <StylesProvider injectFirst>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Header />
        <AppRoutes />
      </ThemeProvider>
    </StylesProvider>
  );
}
