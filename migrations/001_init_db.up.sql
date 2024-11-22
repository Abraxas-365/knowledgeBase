-- First, create the timestamp update function
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create the tables
CREATE TABLE "user" (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    provider TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_whitelist (
    email TEXT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    user_email TEXT NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Modified chatUser table with timestamp columns
CREATE TABLE chatUser (
    id TEXT PRIMARY KEY,
    age INTEGER,
    gender TEXT,
    occupation TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for chatUser table
CREATE TRIGGER update_chatuser_timestamp
    BEFORE UPDATE ON chatUser
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

CREATE TABLE interactions (
    id SERIAL PRIMARY KEY,
    user_chat_id TEXT NOT NULL REFERENCES chatUser(id) ON DELETE CASCADE,
    context_interaction TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create triggers for each table that has updated_at column

-- For "user" table
CREATE TRIGGER update_user_timestamp
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- For files table
CREATE TRIGGER update_files_timestamp
    BEFORE UPDATE ON files
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();

-- For interactions table
CREATE TRIGGER update_interactions_timestamp
    BEFORE UPDATE ON interactions
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
