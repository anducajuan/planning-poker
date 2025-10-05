import { Box, Button, Grid, Typography } from "@mui/material";
import { theme } from "../../theme/theme";
import type { Theme } from "@mui/material/styles";
import { styled } from "@mui/material/styles";
import api from "../../services/api";
import { v4 as uuidv4 } from "uuid";
import { toast } from "react-toastify";
//import { useNavigate } from "react-router-dom";

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
  //const navigate = useNavigate();

  const handleCreateSession = async () => {
    try {
      const response = await api.post("/sessions", {
        name: uuidv4(),
      });

      console.log(response.data);
      //localStorage.setItem("sessionId", "");
      //navigate(`/session/${}`);
    } catch (error) {
      toast.error("Erro ao criar sessão. Tente novamente.");
      console.log(error);
    }
  };

  return (
    <Grid
      container
      xs={12}
      sx={{
        position: "relative",
        overflow: "hidden",
        backgroundColor: theme.palette.background.paper,
        height: 640,
        color: theme.palette.primary.main,
        padding: "48px calc(15% + 24px)",
        clipPath: "polygon(0 0, 100% 0, 100% 100%, 0% 70%)",
      }}
    >
      <Grid item xs={12} md={6}>
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
            onClick={() => {
              handleCreateSession();
            }}
          >
            Criar sessão
          </Button>
        </Grid>
      </Grid>
      <Grid item xs={12} md={6} justifyContent={"flex-end"} display={"flex"}>
        <Box>
          <Image src="./session.png" alt="Session" />
        </Box>
      </Grid>
    </Grid>
  );
}
