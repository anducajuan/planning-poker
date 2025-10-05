
CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL
);


CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE
);


CREATE TABLE stories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('ACTUAL', 'OLD')),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE
);


CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    vote INTEGER NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE
);
