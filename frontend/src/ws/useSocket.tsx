// ws/useWebsocket.ts
import { useEffect, useState } from "react";
import { socketManager } from "./socketManager";
import { handlers } from "./handlers";
import type { Player, Story } from "../pages/Session";

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

    return () => socketManager.disconnect();
  }, [sessionId]);

  return { ...state, sendEvent: socketManager.send.bind(socketManager) };
};
