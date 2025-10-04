import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: "#EBEEF0",
      light: "#25516B",
      dark: "#1A3443",
      contrastText: "#C3C3C3",
    },
    secondary: {
      main: "#8D9CB4",
      dark: "#748497",
      contrastText: "#FFFFFF",
    },
    background: {
      default: "#191919",
      paper: "#111",
    },
  },
  typography: {
    fontFamily: "Inter, Roboto, Arial, sans-serif",
  },
});
