import { Box, Button, Grid, Typography, useMediaQuery } from "@mui/material";
import { theme } from "../../theme/theme";
import type { Theme } from "@mui/material/styles";
import { styled } from "@mui/material/styles";

export const Image = styled("img")(({ theme }: { theme: Theme }) => ({
  width: "100%",
  maxWidth: 456,
  borderRadius: 24,
  border: `0.35rem ridge ${theme.palette.primary.dark}`,
  transition: "transform 0.3s ease-in-out, box-shadow 0.3s ease-in-out",
  "&:hover": {
    transform: "scale(1.02)",
    boxShadow: "0 4px 20px rgba(0,0,0,0.3)",
  },
}));

export function Home() {
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  return (
    <Grid
      container
      xs={12}
      sx={{
        position: "relative",
        overflow: "hidden",
        backgroundColor: theme.palette.background.paper,
        height: isMobile ? 864 : 640,
        color: theme.palette.primary.main,
        padding: "48px calc(15% + 24px)",
        clipPath: isMobile
          ? "polygon(0 0, 100% 0, 100% 100%, 0% 85%)"
          : "polygon(0 0, 100% 0, 100% 100%, 0% 70%)",
      }}
      justifyContent={"space-between"}
      gap={2}
    >
      <Grid item xs={12} md={5}>
        <Typography
          sx={{
            fontSize: 48,
            fontWeight: "bold",
          }}
        >
          Planning poker
        </Typography>
        <Typography
          sx={{
            marginTop: 2,
            fontSize: 22,
            maxWidth: 600,
            lineHeight: 1.4,
          }}
        >
          Ferramenta simples e eficaz para estimar tarefas em equipe utilizando
          a técnica de Planning Poker.
        </Typography>
        <Grid style={{ marginTop: 70 }}>
          <Button
            variant="contained"
            style={{ width: 240, height: 40 }}
            href="/session"
          >
            Criar uma sessão
          </Button>
        </Grid>
      </Grid>
      <Grid item xs={12} md={5} justifyContent={"flex-end"} display={"flex"}>
        <Box>
          <Image src="./session.png" alt="Session" />
        </Box>
      </Grid>
    </Grid>
  );
}
