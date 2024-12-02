-- migrations/0003_extend_songs.down.sql
ALTER TABLE
    IF EXISTS songs DROP COLUMN release_date;

--
ALTER TABLE
    IF EXISTS songs DROP COLUMN "text";

--
ALTER TABLE
    IF EXISTS songs DROP COLUMN link;