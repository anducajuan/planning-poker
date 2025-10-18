ALTER TABLE
    votes
ALTER COLUMN
    vote TYPE text USING vote :: text;