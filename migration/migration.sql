CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL,
    token VARCHAR(32) NOT NULL UNIQUE,
    expire_time DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    likes INT NOT NULL,
    dislikes INT NOT NULL,
    image TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    category_id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories_posts_association (
    association_id UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    post_id UUID NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (category_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS votes(
    vote_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    post_id UUID,
    comment_id UUID,
    binary INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (comment_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    likes INT NOT NULL,
    dislikes INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_user_id_sessions ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_token_sessions ON sessions (token);

CREATE INDEX IF NOT EXISTS idx_user_id_posts ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_title_posts ON posts (title);

CREATE INDEX IF NOT EXISTS idx_post_id_comments ON comments (post_id);
CREATE INDEX IF NOT EXISTS idx_user_id_comments ON comments (user_id);

CREATE INDEX IF NOT EXISTS idx_category_id_association ON categories_posts_association (category_id);
CREATE INDEX IF NOT EXISTS idx_post_id_association ON categories_posts_association (post_id);

CREATE INDEX IF NOT EXISTS idx_vote_id ON votes (vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_user_id ON votes (user_id);
CREATE INDEX IF NOT EXISTS idx_vote_post_id ON votes (post_id);

CREATE INDEX IF NOT EXISTS idx_comment_id ON comments (comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comment_post_id ON comments (post_id);
CREATE INDEX IF NOT EXISTS idx_comment_comment_id ON comments (comment_id);
