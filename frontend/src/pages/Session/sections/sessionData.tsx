import { Grid, Typography } from "@mui/material";
import { styled } from "@mui/material/styles";
import { theme } from "../../../theme/theme";
import { useEffect, useState } from "react";
import api from "../../../services/api";
import type { Story } from "..";
import { useMatch } from "react-router-dom";

export const GridTable = styled(Grid)(() => ({
  width: "100%",
}));

export const SessionData = ({ open }: { open: boolean }) => {
  const match = useMatch("/session/:sessionId");
  const sessionId = match?.params?.sessionId;

  const [stories, setStories] = useState<Story[]>([]);

  useEffect(() => {
    const loadStories = async () => {
      const { data: stories } = await api.get(
        `/stories?session_id=${sessionId}`
      );
      setStories(stories);
    };

    if (sessionId && open) {
      loadStories();
    }
  }, [sessionId, open]);

  return (
    <GridTable container>
      {stories.map((story) => (
        <Grid
          item
          xs={12}
          key={story.id}
          style={{
            backgroundColor: theme.palette.background.default,
            padding: "4px 24px",
            height: 56,
            display: "flex",
            alignItems: "center",
            borderRadius: 8,
            marginBottom: 8,
            justifyContent: "space-between",
          }}
        >
          <Typography
            style={{
              color: theme.palette.primary.main,
            }}
          >
            {story.name}
          </Typography>
          <Typography
            style={{
              color: theme.palette.primary.main,
            }}
          >
            4
            {
              // Utilizar a média da story
              1 > 1 ? " hora" : " horas"
            }
          </Typography>
        </Grid>
      ))}
    </GridTable>
  );
};
