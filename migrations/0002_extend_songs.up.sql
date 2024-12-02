-- migrations/0002_extend_songs.up.sql

-- Проверяем, существует ли колонка release_date, и если нет, добавляем ее
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='songs' AND column_name='release_date') THEN
        ALTER TABLE songs ADD COLUMN release_date DATE;
    END IF;
END $$;

-- Проверяем, существует ли колонка "text", и если нет, добавляем ее
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='songs' AND column_name='text') THEN
        ALTER TABLE songs ADD COLUMN "text" TEXT;
    END IF;
END $$;

-- Проверяем, существует ли колонка link, и если нет, добавляем ее
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='songs' AND column_name='link') THEN
        ALTER TABLE songs ADD COLUMN link TEXT;
    END IF;
END $$;