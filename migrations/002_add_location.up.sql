ALTER TABLE chatUser
ADD COLUMN location TEXT;

UPDATE chatUser
SET location = ''
WHERE location IS NULL;

ALTER TABLE chatUser
ALTER COLUMN location SET NOT NULL,
ALTER COLUMN location SET DEFAULT '';
