import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Grid,
  useMediaQuery,
} from "@mui/material";
import { styled } from "@mui/material/styles";
import type { Theme } from "@mui/material/styles";
import { theme } from "../theme/theme";
import CottageIcon from "@mui/icons-material/Cottage";
import GridViewIcon from "@mui/icons-material/GridView";
import Logo from "./Logo";
import { useNavigate } from "react-router-dom";

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

export default function Header() {
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const navigate = useNavigate();

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
          <NavButton href="/" startIcon={<CottageIcon />}>
            Home
          </NavButton>
          <NavButton href="/session" startIcon={<GridViewIcon />}>
            Sessão
          </NavButton>
        </Grid>
      </StyledToolbar>
    </AppBar>
  );
}
