import { Box, Button, Grid, Modal, TextField, Typography } from "@mui/material";
import { theme } from "../../../theme/theme";
import { useState } from "react";

export const CreateModal = ({
  open,
  handleClose,
  handleCreate,
  type,
}: {
  open: boolean;
  handleClose: () => void;
  handleCreate: (name: string) => void;
  type: string;
}) => {
  const [name, setName] = useState<string>("");

  return (
    <Modal
      open={open}
      onClose={handleClose}
      aria-labelledby="modal-title"
      aria-describedby="modal-description"
    >
      <Box
        sx={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          minWidth: 400,
          bgcolor: "background.paper",
          borderRadius: 2,
          boxShadow: 8,
          p: 4,
        }}
      >
        <Typography
          style={{
            color: theme.palette.primary.main,
            fontWeight: "bold",
            fontSize: 24,
          }}
        >
          {type === "story" ? "Nova votação" : "Juntar-se à mesa"}
        </Typography>
        <Grid
          item
          xs={12}
          justifyContent={"center"}
          display={"flex"}
          style={{ margin: "16px 0px 16px" }}
        >
          <TextField
            label={type === "story" ? "Votação" : "Seu nome"}
            fullWidth
            variant="standard"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </Grid>
        <Grid container gap={2} justifyContent={"flex-end"}>
          <Button variant="outlined" onClick={handleClose} sx={{ mt: 2 }}>
            Fechar
          </Button>
          <Button
            variant="contained"
            onClick={() => handleCreate(name)}
            sx={{ mt: 2 }}
          >
            Salvar
          </Button>
        </Grid>
      </Box>
    </Modal>
  );
};
