CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(15) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    member_1 INTEGER NOT NULL,
    member_2 INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_member1 FOREIGN KEY (member_1) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_member2 FOREIGN KEY (member_2) REFERENCES users(id) ON DELETE CASCADE,
    
    CONSTRAINT force_member_order CHECK (member_1 < member_2),
    CONSTRAINT unique_chat_pair UNIQUE (member_1, member_2)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_sender_id FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_chat_id FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);

