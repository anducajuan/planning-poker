alter table
    votes
add
    column "status" varchar check (status in ('HIDDEN', 'REVEALED')) DEFAULT 'HIDDEN';