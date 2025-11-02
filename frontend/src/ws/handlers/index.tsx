import { playerJoinedHandler, playerLeftHandler } from "./playerHandler";
import { storyCreatedHandler, storyRevealedHandler } from "./storyHandler";
import { voteCreatedHandler } from "./voteHandler";

export const handlers = {
  STORY_REVEALED: storyRevealedHandler,
  STORY_CREATED: storyCreatedHandler,
  VOTE_CREATED: voteCreatedHandler,
  USER_JOINED: playerJoinedHandler,
  USER_LEFT: playerLeftHandler,
};
