import { Grid, Typography } from "@mui/material";
import { styled } from "@mui/material/styles";
import type { Theme } from "@mui/material/styles";
import { theme } from "../theme/theme";

interface CardProps {
  value?: string;
  selected?: boolean;
  color: string;
  onClick: () => void;
}

export const PokerCard = styled(Grid)(
  ({
    theme,
    selected,
    color,
  }: {
    theme: Theme;
    selected: boolean;
    color: string;
  }) => ({
    backgroundColor: color,
    height: selected ? 162 : 144,
    width: selected ? 124 : 108,
    borderRadius: 12,
    boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
    color: theme.palette.primary.main,
    margin: "16px",
    padding: "16px",
    display: "flex",
    justifyContent: "space-between",
    cursor: "pointer",
    transition: "all 0.2s ease 0s",
    "&:hover": {
      boxShadow: "0px 6px 16px rgba(0, 0, 0, 0.2)",
      transform: "translateY(-4px)",
    },
    [theme.breakpoints.down("sm")]: {
      height: 90,
      margin: 0,
    },
  })
);

export default function Card({ value, selected, color, onClick }: CardProps) {
  return (
    <PokerCard
      container
      selected={selected || false}
      color={color}
      theme={theme}
      onClick={onClick}
    >
      <Grid xs={12} style={{ height: "15%" }}>
        <Typography style={{ fontSize: 12 }}>{value || "?"}</Typography>
      </Grid>
      <Grid
        xs={12}
        style={{ height: "70%" }}
        alignItems={"center"}
        display="flex"
        justifyContent="center"
      >
        <Typography style={{ fontSize: 56 }}>{value || "?"}</Typography>
      </Grid>
      <Grid
        xs={12}
        style={{ height: "15%" }}
        justifyContent={"flex-end"}
        display="flex"
      >
        <Typography style={{ fontSize: 12 }}>{value || "?"}</Typography>
      </Grid>
    </PokerCard>
  );
}
