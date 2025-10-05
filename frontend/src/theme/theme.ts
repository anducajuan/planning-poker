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
      main: "#a8db93",
      light: "#2995A3",
      dark: "#1c6c7e",
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
    MuiTextField: {
      variants: [
        {
          props: { variant: "standard" },
          style: ({ theme }) => ({
            "& .MuiInputBase-input": {
              color: theme.palette.primary.main,
            },
            "& .MuiInputLabel-root": {
              color: theme.palette.primary.contrastText,
            },
            "& .MuiInput-underline:": {
              borderBottomColor: theme.palette.primary.dark,
            },
            "& .MuiInput-underline:before": {
              borderBottomColor: theme.palette.primary.dark,
              color: theme.palette.primary.main,
            },
            "& .MuiInput-underline:hover:before": {
              borderBottomColor: theme.palette.primary.dark,
            },
            "& .MuiInput-underline:after": {
              borderBottomColor: theme.palette.primary.dark,
            },
            "& .MuiInput-underline:hover:after": {
              borderBottomColor: theme.palette.primary.dark,
            },
          }),
        },
      ],
    },
  },
});
