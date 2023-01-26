create table if not exists "chats" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "user_id" INT NOT NULL REFERENCES users(id),
    "chat_type" VARCHAR(20) NOT NULL CHECK ("chat_type" IN ('private_chat', 'group_chat')),
    "image_url" TEXT
);

CREATE TABLE IF NOT EXISTS "chat_members" (
    "id" SERIAL PRIMARY KEY,
    "chat_id" INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    "user_id" INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS "chat_messages" (
    "id" SERIAL PRIMARY KEY,
    "message" TEXT NOT NULL,
    "user_id" INT NOT NULL REFERENCES users(id),
    "chat_id" INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);