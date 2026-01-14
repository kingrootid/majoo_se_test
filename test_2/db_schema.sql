-- Database Schema for Blog REST API
-- This file describes the database structure for the blog API system

-- Users table
-- Stores user account information
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for user table to improve query performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Posts table
-- Stores blog post information linked to users
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraint
    CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for posts table to improve query performance
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);

-- Comments table
-- Stores comments linked to posts and users
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for comments table to improve query performance
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);

-- Example of how to create a simple migration with sample data
-- Uncomment and run these queries if you want to populate sample data:
/*
INSERT INTO users (username, email, password, created_at, updated_at) VALUES
('johndoe', 'john@example.com', '$2a$10$8K1KqZU1NzFvQhJp9oTj.eVQhBqYq6YxRtOc7nJmG3jL8V2s5Yzv.', NOW(), NOW()),
('janedoe', 'jane@example.com', '$2a$10$8K1KqZU1NzFvQhJp9oTj.eVQhBqYq6YxRtOc7nJmG3jL8V2s5Yzv.', NOW(), NOW());

INSERT INTO posts (title, content, user_id, created_at, updated_at) VALUES
('Welcome to My Blog', 'This is the first post on my new blog. Excited to share my thoughts with the world!', 1, NOW(), NOW()),
('Travel Adventures', 'Visited the beautiful mountains last weekend.', 2, NOW(), NOW());

INSERT INTO comments (content, post_id, user_id, created_at) VALUES
('Great post!', 1, 2, NOW()),
('Amazing experience!', 2, 1, NOW());
*/