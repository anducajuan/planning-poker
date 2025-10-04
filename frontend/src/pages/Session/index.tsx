import { Grid } from "@mui/material";
//import { useParams } from "react-router-dom";
import Card from "../../components/Cards";
import { mapearCor } from "../../utils/colors";
import { useState } from "react";

export function Session() {
  //const { sessionId } = useParams();
  const cards = [1, 2, 3, 5, 8, 13, 21, 34, 55, "∞", "?", "😴"];
  const [selectedCard, setSelectedCard] = useState<string | number | null>(
    null
  );

  const handleCardClick = (card: string | number) => {
    setSelectedCard(card);
  };

  return (
    <Grid
      container
      justifyContent="center"
      alignItems="center"
      style={{ minHeight: "200px" }}
    >
      {cards.map((card) => (
        <Card
          key={card}
          value={String(card)}
          selected={card == selectedCard}
          color={mapearCor({ valor: card })}
          onClick={() => handleCardClick(card)}
        />
      ))}
    </Grid>
  );
}
