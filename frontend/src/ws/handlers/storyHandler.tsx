import type { Story } from "../../pages/Session";

export const storyRevealedHandler = (
  _payload: unknown,
  setState: React.Dispatch<React.SetStateAction<{ revealed: boolean }>>
) => {
  setState((prev) => ({
    ...prev,
    revealed: true,
  }));
};

export const storyCreatedHandler = (
  payload: { story: Story },
  setState: React.Dispatch<
    React.SetStateAction<{
      story: Story;
    }>
  >
) => {
  setState((prev) => ({
    ...prev,
    story: {
      id: Number(payload.story.id),
      name: payload.story.name || "",
      status: payload.story.status as unknown as "ACTUAL" | "OLD",
    },
  }));
};
