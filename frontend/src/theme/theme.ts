import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: "#EBEEF0",
      light: "#FFE46A",
      dark: "#EA4700",
      contrastText: "#C3C3C3",
    },
    secondary: {
      main: "#8CC0B7",
      light: "#2995A3",
      dark: "#0B1926",
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
