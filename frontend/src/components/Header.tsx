import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Grid,
  useMediaQuery,
  Backdrop,
} from "@mui/material";
import { styled } from "@mui/material/styles";
import type { Theme } from "@mui/material/styles";
import { theme } from "../theme/theme";
import CottageIcon from "@mui/icons-material/Cottage";
import GridViewIcon from "@mui/icons-material/GridView";
import LibraryBooksIcon from "@mui/icons-material/LibraryBooks";
import Logo from "./Logo";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { SessionData } from "../pages/Session/sections/sessionData";

export const StyledToolbar = styled(Toolbar)(({ theme }: { theme: Theme }) => ({
  height: 96,
  backgroundColor: theme.palette.background.paper,
  display: "flex",
  justifyContent: "space-between",
  margin: "0px 15%",
  [theme.breakpoints.down("md")]: {
    height: 64,
    margin: 0,
  },
}));

export const NavButton = styled(Button)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.primary.contrastText,
  textTransform: "none",
  backgroundColor: "transparent",
  boxShadow: "none",
  "&:hover": {
    color: theme.palette.secondary.contrastText,
    backgroundColor: "transparent",
    boxShadow: "none",
  },
}));

export const LogoText = styled(Typography)(({ theme }: { theme: Theme }) => ({
  cursor: "pointer",
  "&:hover": {
    transition: "all 0.5s ease 0s",
    color: theme.palette.primary.main,
  },
}));

export const History = styled(Grid, {
  shouldForwardProp: (prop) => prop !== "open",
})<{ open: boolean }>(({ open, theme }) => ({
  position: "fixed",
  top: 0,
  right: 0,
  width: "400px",
  height: "100%",
  backgroundColor: theme.palette.background.paper,
  padding: "16px",
  boxShadow: "-4px 0 24px rgba(0, 0, 0, 0.1)",
  transform: open ? "translateX(0)" : "translateX(100%)",
  opacity: open ? 1 : 0,
  transition: "transform 0.2s ease-out, opacity 0.2s ease-out",
  zIndex: 1301,
}));

export default function Header() {
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const navigate = useNavigate();

  const [openHistory, setOpenHistory] = useState(false);

  return (
    <AppBar
      position="static"
      color="default"
      elevation={0}
      sx={{ backgroundColor: theme.palette.background.paper, width: "100%" }}
    >
      <StyledToolbar>
        <Grid container alignItems="center" gap={1} xs={isMobile ? 6 : 4}>
          <Logo
            sx={{ fontSize: 48, cursor: "pointer" }}
            onClick={() => {
              navigate("/");
            }}
          />
          <LogoText
            variant="h6"
            sx={{ color: theme.palette.primary.contrastText }}
            onClick={() => {
              navigate("/");
            }}
          >
            Planning Poker
          </LogoText>
        </Grid>
        <Grid
          container
          alignItems="center"
          justifyContent={"flex-end"}
          gap={isMobile ? 1 : 2}
          xs={isMobile ? 6 : 8}
        >
          <NavButton
            onClick={() => setOpenHistory(true)}
            startIcon={<LibraryBooksIcon />}
          >
            Histórico
          </NavButton>
          <NavButton href="/" startIcon={<CottageIcon />}>
            Home
          </NavButton>
          <NavButton href="/session" startIcon={<GridViewIcon />}>
            Sessão
          </NavButton>
        </Grid>
      </StyledToolbar>
      <Backdrop
        open={openHistory}
        onClick={() => setOpenHistory(false)}
        sx={{
          zIndex: 1300,
          backgroundColor: "rgba(0,0,0,0.35)",
          transition: "opacity 0.2s ease-out",
        }}
      />
      <History container open={openHistory}>
        <Grid
          item
          xs={12}
          display={"flex"}
          gap={2}
          style={{
            height: 48,
            boxShadow: `0px 0.5px 0 0 ${theme.palette.primary.contrastText}`,
            alignItems: "center",
            paddingBottom: 8,
          }}
        >
          <LibraryBooksIcon style={{ color: theme.palette.primary.main }} />
          <Typography
            variant="h6"
            sx={{ lineHeight: 1.35, color: theme.palette.primary.main }}
          >
            Votações Anteriores
          </Typography>
        </Grid>
        <Grid
          item
          xs={12}
          style={{
            marginTop: 16,
            height: "100%",
            overflowY: "auto",
          }}
        >
          <SessionData open={openHistory} />
        </Grid>
      </History>
    </AppBar>
  );
}
