import { Box, Button, CircularProgress, Grid, Typography } from "@mui/material";
import type { Theme } from "@mui/material/styles";
import { alpha, styled } from "@mui/material/styles";
import React, { useEffect, useState } from "react";
import Card from "../../../components/Cards";
import { theme } from "../../../theme/theme";
import { mapearCor } from "../../../utils/colors";
import type { Player, Story } from "..";
import Logo from "../../../components/Logo";
import { CreateModal } from "./sessionModal";
import api from "../../../services/api";
import { AxiosError } from "axios";
import { toast } from "react-toastify";

const tableSides = [
  [1, 5, 9, 11],
  [2, 6],
  [3, 7, 10, 12],
  [4, 8],
];

export const GridTable = styled(Grid)(() => ({
  width: "100%",
  margin: "16px 0px",
  padding: "0px 36px",
}));

export const BoxTable = styled(Box)(({ theme }: { theme: Theme }) => ({
  width: "40%",
  height: "236px",
  backgroundColor: alpha(theme.palette.primary.dark, 0.3),
  opacity: 0.7,
  borderRadius: 24,
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  border: `1.5rem ridge ${theme.palette.primary.dark}`,
  transition: "transform 0.3s ease-in-out, box-shadow 0.3s ease-in-out",
  "&:hover": {
    transform: "scale(1.01)",
    boxShadow: "0 4px 20px rgba(0,0,0,0.2)",
  },
}));

export const VoteTable = ({
  playersList,
  player,
  handleCreateUser,
  openUserModal,
  setOpenUserModal,
  story,
  handleCreateStory,
  openStoryModal,
  setOpenStoryModal,
  setSelectedCard,
  isRevealed,
}: {
  playersList: Player[];
  player: Player | undefined;
  handleCreateUser: (name: string) => void;
  openUserModal: boolean;
  setOpenUserModal: React.Dispatch<React.SetStateAction<boolean>>;
  story: Story | undefined;
  handleCreateStory: (storyName: string) => void;
  openStoryModal: boolean;
  setOpenStoryModal: React.Dispatch<React.SetStateAction<boolean>>;
  setSelectedCard: React.Dispatch<React.SetStateAction<string | null | number>>;
  isRevealed: boolean;
}) => {
  const [players, setPlayers] = useState<Player[]>([]);

  useEffect(() => {
    setPlayers(
      playersList?.map((p) =>
        p.id === player?.id ? { ...p, vote: player?.vote } : p
      )
    );
  }, [playersList, player]);

  const handleReveal = async () => {
    try {
      await api.post(`/stories/${story?.id}/reveal`);
      setSelectedCard("");
    } catch (error: unknown) {
      if (error instanceof AxiosError) {
        toast.error(error?.response?.data.message);
      } else {
        toast.error("Ocorreu um erro ao revelar os votos.");
      }
    }
  };

  return (
    <GridTable container justifyContent={"center"}>
      {!players || players?.length == 0 ? (
        <Grid>
          <CircularProgress />
        </Grid>
      ) : (
        <>
          <Grid
            item
            xs={1}
            justifyContent={"center"}
            alignItems={"center"}
            display={"flex"}
          >
            <Grid item xs={12}>
              <MapPlayers
                idxs={tableSides[3]}
                players={players.filter((player: Player) =>
                  tableSides[3].includes(player.position)
                )}
                isRevealed={isRevealed}
              />
            </Grid>
          </Grid>
          <Grid item xs={8}>
            <Grid container xs={12} gap={2}>
              <Grid
                container
                justifyContent={"space-evenly"}
                alignItems={"center"}
              >
                <MapPlayers
                  idxs={tableSides[0]}
                  players={players.filter((player: Player) =>
                    tableSides[0].includes(player.position)
                  )}
                  isRevealed={isRevealed}
                />
              </Grid>
              <Grid
                item
                xs={12}
                display={"flex"}
                justifyContent={"center"}
                alignItems={"center"}
              >
                <BoxTable>
                  {!player || player.position !== 1 ? (
                    <>
                      {!player ? (
                        <Button
                          variant="contained"
                          onClick={() => setOpenUserModal(true)}
                          style={{
                            minWidth: "40%",
                            height: "44px",
                          }}
                        >
                          Juntar-se à mesa
                        </Button>
                      ) : (
                        <Logo />
                      )}
                    </>
                  ) : (
                    <Button
                      variant="contained"
                      onClick={() =>
                        story?.id ? handleReveal() : setOpenStoryModal(true)
                      }
                      style={{
                        minWidth: "40%",
                        height: "44px",
                      }}
                    >
                      {story?.id ? "Revelar" : "Criar votação"}
                    </Button>
                  )}
                </BoxTable>
              </Grid>
              <Grid
                container
                justifyContent={"space-evenly"}
                alignItems={"center"}
              >
                <MapPlayers
                  idxs={tableSides[2]}
                  players={players.filter((player: Player) =>
                    tableSides[2].includes(player.position)
                  )}
                  isRevealed={isRevealed}
                />
              </Grid>
            </Grid>
          </Grid>
          <Grid
            item
            xs={1}
            justifyContent={"center"}
            alignItems={"center"}
            display={"flex"}
          >
            <Grid item xs={12}>
              <MapPlayers
                idxs={tableSides[1]}
                players={players.filter((player: Player) =>
                  tableSides[1].includes(player.position)
                )}
                isRevealed={isRevealed}
              />
            </Grid>
          </Grid>
        </>
      )}
      <CreateModal
        open={openUserModal || openStoryModal}
        handleClose={() =>
          openStoryModal ? setOpenStoryModal(false) : setOpenUserModal(false)
        }
        handleCreate={openStoryModal ? handleCreateStory : handleCreateUser}
        type={openStoryModal ? "story" : "player"}
      />
    </GridTable>
  );
};

const MapPlayers = ({
  idxs,
  players,
  isRevealed,
}: {
  idxs: Array<number>;
  players: Array<Player>;
  isRevealed: boolean;
}) => {
  return (
    <>
      {players.map(
        (player, index) =>
          idxs.includes(player.position) && (
            <Grid
              item
              key={index}
              alignItems={"center"}
              display={"flex"}
              direction={"column"}
              style={{
                marginTop: "-16px",
              }}
            >
              <Card
                value={isRevealed ? String(player.vote) : ""}
                selected={false}
                fontColor={
                  !["∞", "?", "😴"].includes(String(player.vote))
                    ? mapearCor({ valor: player.vote })
                    : ""
                }
                color={
                  isRevealed
                    ? theme.palette.background.paper
                    : theme.palette.background.paper
                }
                scale={0.75}
              />
              <Typography
                key={index}
                style={{
                  textAlign: "center",
                  color: theme.palette.primary.contrastText,
                  fontSize: 15,
                  fontWeight: "bold",
                  marginTop: "-16px",
                }}
              >
                {player.name}
              </Typography>
            </Grid>
          )
      )}
    </>
  );
};
