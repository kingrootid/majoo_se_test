# Database Schema Documentation

## Overview
This document describes the database schema for the Blog REST API system. The system uses PostgreSQL as the primary database with GORM as the ORM layer.

## Tables

### 1. Users Table (`users`)
Stores user account information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Unique identifier for the user |
| username | VARCHAR(255) | UNIQUE, NOT NULL | User's display name |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User's email address |
| password | VARCHAR(255) | NOT NULL | Hashed password |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | Record update timestamp |
| deleted_at | TIMESTAMP WITH TIME ZONE | - | Soft delete timestamp |

**Indexes:**
- `idx_users_username` on `username`
- `idx_users_email` on `email`
- `idx_users_deleted_at` on `deleted_at`

### 2. Posts Table (`posts`)
Stores blog post information linked to users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Unique identifier for the post |
| title | VARCHAR(255) | NOT NULL | Post title |
| content | TEXT | NOT NULL | Post content |
| user_id | INTEGER | NOT NULL, FOREIGN KEY | Reference to the author user |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | Record update timestamp |
| deleted_at | TIMESTAMP WITH TIME ZONE | - | Soft delete timestamp |

**Foreign Keys:**
- `fk_posts_user` references `users.id` with CASCADE delete

**Indexes:**
- `idx_posts_user_id` on `user_id`
- `idx_posts_created_at` on `created_at`
- `idx_posts_deleted_at` on `deleted_at`

### 3. Comments Table (`comments`)
Stores comments linked to posts and users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Unique identifier for the comment |
| content | TEXT | NOT NULL | Comment content |
| post_id | INTEGER | NOT NULL, FOREIGN KEY | Reference to the post |
| user_id | INTEGER | NOT NULL, FOREIGN KEY | Reference to the comment author |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT CURRENT_TIMESTAMP | Record creation timestamp |
| deleted_at | TIMESTAMP WITH TIME ZONE | - | Soft delete timestamp |

**Foreign Keys:**
- `fk_comments_post` references `posts.id` with CASCADE delete
- `fk_comments_user` references `users.id` with CASCADE delete

**Indexes:**
- `idx_comments_post_id` on `post_id`
- `idx_comments_user_id` on `user_id`
- `idx_comments_created_at` on `created_at`
- `idx_comments_deleted_at` on `deleted_at`

## Relationships

### User → Posts (One-to-Many)
- One user can create many posts
- Cascade delete: If a user is deleted, all their posts are also deleted

### Post → Comments (One-to-Many)
- One post can have many comments
- Cascade delete: If a post is deleted, all its comments are also deleted

### User → Comments (One-to-Many)
- One user can make many comments
- Cascade delete: If a user is deleted, all their comments are also deleted

## Data Integrity

### Constraints
- Unique constraints on usernames and emails
- Foreign key constraints with cascade deletion
- Not-null constraints on required fields

### Soft Delete Strategy
- All tables except comments support soft deletes using the `deleted_at` field
- This preserves referential integrity while hiding deleted records
- GORM handles soft delete filtering automatically

## Security Considerations

### Password Storage
- Passwords are stored as bcrypt hashes
- Plain text passwords are never stored in the database
- Passwords are filtered out from JSON responses using `-` tag

### Access Control
- Each entity has proper ownership relationships
- Authorization checks ensure users can only modify their own content
- Foreign key constraints enforce referential integrity

## Performance Optimizations

### Indexing Strategy
- Primary keys are automatically indexed
- Foreign key columns are indexed for join performance
- Searchable fields (username, email) are indexed
- Timestamp columns are indexed for sorting and filtering
- Soft delete columns are indexed for query optimization

## Migration Notes

### Initial Migration
When deploying the application for the first time, the auto-migration feature will create all necessary tables, indexes, and constraints as defined by the GORM models.

### Future Changes
- Schema changes should be handled through proper migration scripts
- Always backup the database before running migrations
- Test migrations on a copy of production data before applying