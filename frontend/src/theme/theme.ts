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
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: "none",
          borderRadius: 8,
          fontWeight: 600,
        },
      },
      variants: [
        {
          props: { variant: "contained" },
          style: ({ theme }) => ({
            backgroundColor: theme.palette.primary.dark,
            color: theme.palette.primary.contrastText,
            boxShadow: "0 4px 20px rgba(0,0,0,0.3)",
            "&:hover": {
              backgroundColor: theme.palette.primary.dark,
              color: theme.palette.primary.main,
              boxShadow: "0 4px 20px rgba(0,0,0,0.3)",
            },
          }),
        },
        {
          props: { variant: "outlined" },
          style: ({ theme }) => ({
            border: `2px solid ${theme.palette.primary.dark}`,
            color: theme.palette.primary.dark,
            backgroundColor: "transparent",
            "&:hover": {
              border: `2px solid ${theme.palette.primary.dark}`,
              backgroundColor: "rgba(234, 71, 0, 0.05)",
            },
          }),
        },
      ],
    },
  },
});
