import type { Player } from "../../pages/Session";

export const voteCreatedHandler = (
  payload: { user_id: string | number; vote: string | number },
  setState: React.Dispatch<React.SetStateAction<{ players: Player[] }>>
) => {
  setState((prev: { players: Player[] }) => ({
    ...prev,
    players: prev.players.map((p) =>
      p.id === payload.user_id ? { ...p, vote: payload.vote } : p
    ),
  }));
};
