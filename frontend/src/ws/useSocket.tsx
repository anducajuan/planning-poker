// ws/useWebsocket.ts
import { useEffect, useState } from "react";
import { socketManager } from "./socketManager";
import { handlers } from "./handlers";
import type { Player, Story } from "../pages/Session";
import api from "../services/api";

export const useWebsocket = (sessionId: string) => {
  const [state, setState] = useState<{
    players: Player[];
    revealed: boolean;
    story: Story;
  }>({
    players: [],
    revealed: false,
    story: { id: null, name: "", status: "ACTUAL" },
  });

  useEffect(() => {
    if (!sessionId) return;

    socketManager.connect(sessionId, (event, payload) => {
      const handler = handlers[event as keyof typeof handlers];
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      if (handler) handler(payload as any, setState as any);
      else console.log("[WS] Evento não tratado:", event, payload);
    });

    const fetchInitial = async () => {
      try {
        const { data: playersList } = await api.get(
          `users?session_id=${sessionId}`
        );

        const { data: storyData } = await api.get(
          `stories?session_id=${sessionId}`
        );

        const actualStory = storyData.find(
          (story: Story) => story.status === "ACTUAL"
        );

        const players = (playersList as Player[]).map((p, i) => ({
          id: p.id,
          name: p.name,
          vote: "",
          position: i + 1,
        }));
        const story = actualStory as Story;

        setState((prev) => ({
          ...prev,
          players: players ?? [],
          story: story ?? prev.story,
        }));
      } catch (err) {
        console.error("[WS] Erro ao buscar dados iniciais:", err);
      }
    };

    fetchInitial();

    return () => {
      socketManager.disconnect();
    };
  }, [sessionId]);

  return { ...state, sendEvent: socketManager.send.bind(socketManager) };
};
