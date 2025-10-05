import { Grid, Typography } from "@mui/material";
import { styled } from "@mui/material/styles";
import type { Theme } from "@mui/material/styles";
import { theme } from "../theme/theme";
import Logo from "./Logo";

interface CardProps {
  value?: string;
  selected?: boolean;
  color: string;
  fontColor?: string;
  onClick?: () => void;
  scale?: number;
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
    height: selected ? 136 : 124,
    width: selected ? 96 : 84,
    borderRadius: 12,
    boxShadow: "0px 4px 12px rgba(0, 0, 0, 0.1)",
    color: theme.palette.primary.main,
    margin: "8px",
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
      height: selected ? 136 : 124,
      width: selected ? 96 : 84,
      margin: 4,
    },
  })
);

export default function Card({
  value,
  selected,
  color,
  fontColor,
  onClick,
  scale,
}: CardProps) {
  return (
    <PokerCard
      container
      selected={selected || false}
      color={color}
      theme={theme}
      onClick={onClick}
      style={{
        scale: scale || 1,
        color: fontColor || theme.palette.primary.main,
      }}
    >
      <Grid xs={12} style={{ height: "15%" }}>
        <Typography style={{ fontSize: 12 }}>{value || ""}</Typography>
      </Grid>
      <Grid
        xs={12}
        style={{ height: "70%" }}
        alignItems={"center"}
        display="flex"
        justifyContent="center"
      >
        <Typography style={{ fontSize: value ? 40 : 56 }}>
          {value || <Logo fontSize={"large"} alignmentBaseline={"baseline"} />}
        </Typography>
      </Grid>
      <Grid
        xs={12}
        style={{ height: "15%" }}
        justifyContent={"flex-end"}
        display="flex"
      >
        <Typography style={{ fontSize: 12 }}>{value || ""}</Typography>
      </Grid>
    </PokerCard>
  );
}
