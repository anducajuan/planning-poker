import type { Dispatch, SetStateAction } from "react";
import type { Player, Story } from "../../pages/Session";

interface WSState {
  players: Player[];
  revealed: boolean;
  story: Story;
}

export const playerJoinedHandler = (
  payload: { user: Player },
  setState: Dispatch<SetStateAction<WSState>>
) => {
  setState((prev) => ({
    ...prev,
    players: prev.players.some((p) => p.id === payload.user.id)
      ? prev.players
      : [
          ...prev.players,
          { ...payload.user, position: prev.players.length + 1 },
        ],
  }));
};

export const playerLeftHandler = (
  payload: Player,
  setState: Dispatch<SetStateAction<{ players: Player[] }>>
) => {
  setState((prev: { players: Player[] }) => ({
    ...prev,
    players: prev.players.filter((p) => p.id !== payload.id),
  }));
};
