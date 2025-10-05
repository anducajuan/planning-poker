import { Box, Button, Grid, Typography } from "@mui/material";
import type { Theme } from "@mui/material/styles";
import { alpha, styled } from "@mui/material/styles";
import { useEffect, useState } from "react";
import Card from "../../../components/Cards";
import { theme } from "../../../theme/theme";
import { mapearCor } from "../../../utils/colors";

interface Player {
  id: number;
  name: string;
  vote: string | number;
  position: number;
}

const tableSides = [
  [1, 5, 9, 11],
  [2, 6],
  [3, 7, 10, 12],
  [4, 8],
];

const playersLit = [
  {
    id: 1,
    name: "Gustavo Akyama 1",
    vote: 1,
    position: 1,
  },
  {
    id: 4,
    name: "Erick 1",
    vote: 21,
    position: 2,
  },
  {
    id: 5,
    name: "João 1",
    vote: 55,
    position: 3,
  },
  {
    id: 2,
    name: "Juan 1",
    vote: 8,
    position: 4,
  },
  {
    id: 4,
    name: "Erick 2",
    vote: 21,
    position: 5,
  },
  {
    id: 3,
    name: "Breno 1",
    vote: 55,
    position: 6,
  },
  {
    id: 2,
    name: "Juan 2",
    vote: 8,
    position: 7,
  },
  {
    id: 2,
    name: "Maia 1",
    vote: "?",
    position: 8,
  },
  {
    id: 1,
    name: "Gustavo Akyama 2",
    vote: 1,
    position: 9,
  },
  {
    id: 2,
    name: "Maia 2",
    vote: "?",
    position: 10,
  },
  {
    id: 5,
    name: "João 2",
    vote: 55,
    position: 11,
  },
  {
    id: 3,
    name: "Breno 2",
    vote: 55,
    position: 12,
  },
];

export const GridTable = styled(Grid)(() => ({
  width: "100%",
  margin: "16px 0px",
  padding: "0px 36px",
}));

export const BoxTable = styled(Box)(({ theme }: { theme: Theme }) => ({
  width: "80%",
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

export const VoteTable = () => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [isRevealed, setIsRevealed] = useState<boolean>(true);

  useEffect(() => {
    setPlayers(playersLit);
  }, []);

  return (
    <GridTable container justifyContent={"center"}>
      {players.length === 0 ? (
        <Grid>
          <Typography>Carregando...</Typography>
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
                  <Button
                    variant="contained"
                    style={{
                      minWidth: "40%",
                      height: "44px",
                    }}
                  >
                    Revelar
                  </Button>
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
