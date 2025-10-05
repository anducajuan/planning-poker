import { Button, Grid, TextField, Typography } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import Card from "../../components/Cards";
import { mapearCor } from "../../utils/colors";
import { useState } from "react";
import { styled } from "@mui/material/styles";
import api from "../../services/api";
import { toast } from "react-toastify";
import { theme } from "../../theme/theme";
import { AxiosError } from "axios";

export const SessionNameTextField = styled(TextField)(() => ({
  margin: "0px 15%",
  maxWidth: "360px",
}));

export function Session() {
  const { sessionId } = useParams();

  const cards = [1, 2, 3, 5, 8, 13, 21, 34, 55, "∞", "?", "😴"];
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [previousSession, setPreviousSession] = useState<string | null>(
    localStorage.getItem("sessionId")
  );
  const [selectedCard, setSelectedCard] = useState<string | number | null>(
    null
  );

  const handleCardClick = (card: string | number) => {
    setSelectedCard(card);
  };

  const handleCreateSession = async () => {
    const sessionName = name.trim();

    if (sessionName.length === 0) {
      toast.error("Nome de sessão inválido.");
      return;
    }

    try {
      const response = await api.post("/sessions", {
        name: name,
      });

      const { data: session } = response.data;

      localStorage.setItem("sessionId", session.id);
      setPreviousSession(session.id);
      navigate(`/session/${session.id}`);
    } catch (error: unknown) {
      if (error instanceof AxiosError) {
        toast.error(error?.response?.data.message);
      } else {
        toast.error("Ocorreu um erro ao criar a sessão.");
      }
    }
  };

  return (
    <Grid
      container
      justifyContent="center"
      alignItems="center"
      style={{ minHeight: "200px" }}
    >
      {sessionId ? (
        <>
          {cards.map((card) => (
            <Card
              key={card}
              value={String(card)}
              selected={card == selectedCard}
              color={mapearCor({ valor: card })}
              onClick={() => handleCardClick(card)}
            />
          ))}
        </>
      ) : (
        <Grid
          container
          justifyContent={"center"}
          alignItems={"center"}
          display={"flex"}
          gap={2}
        >
          <Grid item xs={12}>
            <Typography
              style={{
                color: theme.palette.primary.main,
                fontWeight: "bold",
                fontSize: 24,
                margin: "36px 15%",
                textAlign: "center",
              }}
            >
              Crie uma sessão para começar!
            </Typography>
          </Grid>
          <Grid item xs={12} justifyContent={"center"} display={"flex"}>
            <SessionNameTextField
              label="Nome da sessão"
              fullWidth
              variant="standard"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>
          <Grid
            item
            xs={12}
            alignItems={"center"}
            display={"flex"}
            direction={"column"}
            style={{
              margin: "36px 15%",
              maxWidth: "360px",
            }}
            gap={3}
          >
            <Button
              variant="contained"
              fullWidth
              style={{ height: 40 }}
              onClick={() => handleCreateSession()}
            >
              Criar sessão
            </Button>
            {previousSession && (
              <Button
                variant="outlined"
                fullWidth
                style={{ height: 40 }}
                onClick={() => navigate(`/session/${previousSession}`)}
              >
                Restaurar sessão anterior
              </Button>
            )}
          </Grid>
        </Grid>
      )}
    </Grid>
  );
}
