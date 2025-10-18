import { Button, Grid, TextField, Typography } from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import Card from "../../components/Cards";
import { mapearCor } from "../../utils/colors";
import { useEffect, useState } from "react";
import { styled } from "@mui/material/styles";
import api from "../../services/api";
import { toast } from "react-toastify";
import { theme } from "../../theme/theme";
import { AxiosError } from "axios";
import { VoteTable } from "./sections/voteTable";
import { SessionData } from "./sections/sessionData";

export interface Player {
  id: number;
  name: string;
  vote: string | number;
  position: number;
}

export interface Story {
  id: number | null;
  name?: string;
  status?: "ACTUAL" | "OLD";
}

export interface Vote {
  id: number | null;
  vote: string | number;
  user_id?: number;
  story_id?: number;
}

export const SessionTextField = styled(TextField)(() => ({
  margin: "0px 15%",
  maxWidth: "360px",
}));

export function Session() {
  const { sessionId: paramsSessionId } = useParams();

  const cards = [1, 2, 3, 5, 8, 13, 21, 34, 55, "∞", "?", "😴"];
  const navigate = useNavigate();

  const [openUserModal, setOpenUserModal] = useState<boolean>(false);
  const [openStoryModal, setOpenStoryModal] = useState<boolean>(false);
  const [player, setPlayer] = useState<Player>();
  const [players, setPlayers] = useState<Player[]>([]);
  const [name, setName] = useState("");
  const [sessionName, setSessionName] = useState("");
  const [vote, setVote] = useState<Vote>({ id: null, vote: "" });
  const [story, setStory] = useState<Story>({
    id: null,
    name: "",
  });
  const [previousSession, setPreviousSession] = useState<string | null>(
    localStorage.getItem("sessionId")
  );
  const [selectedCard, setSelectedCard] = useState<string | number | null>(
    null
  );

  useEffect(() => {
    const loadPlayers = async () => {
      try {
        const { data: playersList } = await api.get(
          `users?session_id=${paramsSessionId}`
        );

        const formattedPlayersList = playersList.map(
          (player: Player, index: number) => ({
            id: player.id,
            name: player.name,
            position: index + 1,
            vote: "",
          })
        );

        setPlayers(formattedPlayersList);

        const storagePlayer = JSON.parse(localStorage.getItem("user") || "{}");

        if (
          storagePlayer?.username &&
          storagePlayer?.session === paramsSessionId
        ) {
          const currentPlayer = formattedPlayersList.find(
            (p: Player) => p.name === storagePlayer.username
          );

          if (currentPlayer) {
            setPlayer(currentPlayer);
            loadStory(currentPlayer);
          }
        } else {
          if (paramsSessionId) {
            localStorage.removeItem("user");
          }
        }
      } catch (error: unknown) {
        if (error instanceof AxiosError) {
          toast.error(error?.response?.data.message);
        } else {
          toast.error("Ocorreu um erro ao criar a sessão.");
        }
      }
    };

    const loadStory = async (player: Player) => {
      try {
        const { data: storyData } = await api.get(
          `stories?session_id=${paramsSessionId}`
        );

        const actualStory = storyData.find(
          (story: Story) => story.status === "ACTUAL"
        );

        if (actualStory) {
          setStory({
            id: actualStory.id,
            name: actualStory.name,
          });
          loadVotes(actualStory.id, player.id);
        }
      } catch (error: unknown) {
        if (error instanceof AxiosError) {
          toast.error(error?.response?.data.message);
        } else {
          toast.error("Ocorreu um erro ao carregar a votação.");
        }
      }
    };

    const loadVotes = async (storyId: number, playerId: number) => {
      try {
        if (storyId) {
          const { data: votesData } = await api.get(
            `votes?story_id=${storyId}`
          );

          const userVote = votesData.find(
            (vote: Vote) => vote.user_id === playerId
          );

          if (userVote) {
            setVote(userVote);
            setSelectedCard(userVote.vote);
          }
        }
      } catch (error: unknown) {
        if (error instanceof AxiosError) {
          toast.error(error?.response?.data.message);
        } else {
          toast.error("Ocorreu um erro ao carregar o voto.");
        }
      }
    };

    loadPlayers();
  }, [paramsSessionId]);

  const handleCardClick = async (card: string | number) => {
    if (card !== selectedCard && player && paramsSessionId && story) {
      if (!vote?.id) {
        const { data: newVote } = await api.post("/votes", {
          vote: card.toString(),
          session_id: paramsSessionId,
          user_id: player.id,
          story_id: story.id,
        });
        setVote(newVote);
      } else {
        const { data: updatedVote } = await api.patch(`/votes/${vote.id}`, {
          vote: card.toString(),
        });
        setVote(updatedVote);
      }
      setSelectedCard(card);
    }
  };

  const handleCreateSession = async () => {
    const trimmedSessionName = sessionName.trim();
    const trimmedName = name.trim();

    if (trimmedSessionName.length === 0) {
      toast.error("Nome de sessão inválido.");
      return;
    }

    if (trimmedName.length === 0) {
      toast.error("Nome de usuário inválido.");
      return;
    }

    try {
      const { data: session } = await api.post("/sessions", {
        name: trimmedSessionName,
      });

      const player = await handleCreateUser(name, session.id);

      if (!player) return;

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

  const handleCreateUser = async (username: string, sessionId?: string) => {
    const trimmedName = username?.trim();
    const session = sessionId || paramsSessionId;

    if (!session) {
      toast.error("Sessão inválida.");
      return;
    }

    if (trimmedName.length === 0) {
      toast.error("Nome de usuário inválido.");
      return;
    }

    try {
      const { data: user } = await api.post("/users", {
        session_id: session,
        name: trimmedName,
      });

      localStorage.setItem(
        "user",
        JSON.stringify({
          username: user.name,
          session: session,
        })
      );

      const newPlayer = {
        id: user.id,
        name: user.name,
        vote: "",
        position: (players.at(-1)?.position || 0) + 1,
      };

      setPlayer(newPlayer);
      setPlayers([...players, newPlayer]);
      setOpenUserModal(false);

      return newPlayer;
    } catch (error: unknown) {
      if (error instanceof AxiosError) {
        toast.error(error?.response?.data.message);
      } else {
        toast.error("Ocorreu um erro ao criar usuário.");
      }
    }
  };

  const handleCreateStory = async (storyName: string) => {
    const trimmedStoryName = storyName.trim();
    const session = paramsSessionId;

    if (!session) {
      toast.error("Sessão inválida.");
      return;
    }

    if (trimmedStoryName.length === 0) {
      toast.error("Nome de votação inválido.");
      return;
    }

    try {
      const { data: storyData } = await api.post("/stories", {
        name: trimmedStoryName,
        session_id: session,
        status: "ACTUAL",
      });

      setStory({
        id: storyData.id,
        name: storyData.name,
      });
      setOpenStoryModal(false);
    } catch (error: unknown) {
      if (error instanceof AxiosError) {
        toast.error(error?.response?.data.message);
      } else {
        toast.error("Ocorreu um erro ao criar a votação.");
      }
    }
  };

  return (
    <Grid container justifyContent="center" alignItems="center">
      {paramsSessionId ? (
        <Grid container justifyContent={"center"} alignItems={"flex-end"}>
          <Grid item xs={12} lg={3}>
            <SessionData />
          </Grid>
          <Grid item xs={12} lg={6}>
            <VoteTable
              playersList={players}
              player={player}
              handleCreateUser={handleCreateUser}
              setOpenUserModal={setOpenUserModal}
              openUserModal={openUserModal}
              story={story}
              setStory={setStory}
              handleCreateStory={handleCreateStory}
              setOpenStoryModal={setOpenStoryModal}
              openStoryModal={openStoryModal}
            />
          </Grid>
          <Grid item xs={12} lg={3}>
            <SessionData />
          </Grid>

          <Grid
            item
            display={"flex"}
            direction={"row"}
            alignItems={"flex-end"}
            style={{ minHeight: "180px", marginTop: "2%" }}
          >
            <Grid container justifyContent={"center"}>
              {cards.map((card, index) => (
                <Grid item display={"flex"} direction={"row"} key={index}>
                  <Card
                    key={card}
                    value={String(card)}
                    selected={card == selectedCard}
                    color={mapearCor({ valor: card })}
                    onClick={() => handleCardClick(card)}
                  />
                </Grid>
              ))}
            </Grid>
          </Grid>
        </Grid>
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
            <SessionTextField
              label="Nome da sessão"
              fullWidth
              variant="standard"
              value={sessionName}
              onChange={(e) => setSessionName(e.target.value)}
            />
          </Grid>
          <Grid item xs={12} justifyContent={"center"} display={"flex"}>
            <SessionTextField
              label="Seu nome"
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
