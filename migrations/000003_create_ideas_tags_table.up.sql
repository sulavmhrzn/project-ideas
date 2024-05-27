CREATE TABLE IF NOT EXISTS ideas(
    id serial PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL ,
    user_id int REFERENCES users ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT NOW() 
);

CREATE TABLE IF NOT EXISTS tags(
    id serial PRIMARY KEY,
    title text NOT NULL
);

CREATE TABLE IF NOT EXISTS ideas_tags(
    idea_id int REFERENCES ideas ON DELETE CASCADE,
    tag_id int REFERENCES tags ON DELETE CASCADE
);