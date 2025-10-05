import { Grid } from "@mui/material";
import type { Theme } from "@mui/material/styles";
import { styled } from "@mui/material/styles";

export const DataTable = styled(Grid)(({ theme }: { theme: Theme }) => ({
  width: "100%",
  margin: "24px 0px",
  padding: "36px",
}));

export const SessionData = () => {
  return <DataTable container></DataTable>;
};
