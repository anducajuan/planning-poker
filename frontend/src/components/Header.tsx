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

export const StyledToolbar = styled(Toolbar)(({ theme }: { theme: Theme }) => ({
  height: 96,
  backgroundColor: theme.palette.background.paper,
  display: "flex",
  justifyContent: "space-between",
  margin: "0px 15%",
  [theme.breakpoints.down("sm")]: {
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

export default function Header() {
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <AppBar
      position="static"
      color="default"
      elevation={0}
      sx={{ backgroundColor: theme.palette.background.paper }}
    >
      <StyledToolbar>
        <Typography
          variant="h6"
          component="div"
          sx={{ flexGrow: 1, color: theme.palette.primary.contrastText }}
        >
          Planning Poker
        </Typography>
        <Grid container spacing={isMobile ? 2 : 4}>
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
