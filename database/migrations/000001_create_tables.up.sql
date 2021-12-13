CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_name TEXT UNIQUE NOT NULL,
    pwdhash TEXT NOT NULL,
    protected BOOL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS photos (
    id SERIAL PRIMARY KEY,
    file_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS gallery_items (
    id TEXT UNIQUE NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    caption TEXT NOT NULL,
    project_info TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    embed_url TEXT
);
CREATE TABLE IF NOT EXISTS project_images (
    id SERIAL PRIMARY KEY,
    gallery_id TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    CONSTRAINT fk_gallery FOREIGN KEY(gallery_id) REFERENCES gallery_items(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT UNIQUE NOT NULL PRIMARY KEY,
    user_name TEXT UNIQUE NOT NULL,
    user_id SERIAL UNIQUE NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    max_age BIGINT NOT NULL,
    CONSTRAINT fk_user_name FOREIGN KEY(user_name) REFERENCES users(user_name) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);