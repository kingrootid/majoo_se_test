# Blog REST API

A complete RESTful API for a blog system with user authentication, posts, and comments.

## Features

- User authentication and authorization (JWT-based)
- CRUD operations for posts and comments
- Input validation and error responses
- Database integration with PostgreSQL
- Transaction support for database operations

## Prerequisites

- Go 1.19+
- PostgreSQL database
- Environment variables configured

## Environment Variables

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=majoo_test
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
```

## Setup

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Configure environment variables
4. Run the seeding script: `go run cmd/seed/main.go`
5. Start the server: `go run cmd/api/main.go`

## Development with Live Reload (using Air)

For faster development workflow, you can use Air for hot reloading:

1. Install Air:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Or use the binary installer:
   ```bash
   # On macOS
   brew install air

   # On Linux
   curl -fLo air https://git.io/air_for_linux && chmod +x air

   # On Windows (PowerShell)
   Invoke-WebRequest -Uri https://git.io/air_for_windows_amd64.zip -OutFile air.zip
   Expand-Archive -Path air.zip -DestinationPath . -Force
   ```

3. Run the application with live reload:
   ```bash
   air
   ```

Air will automatically rebuild and restart the server when you make changes to the code. The configuration is already set up in `.air.toml`.

## API Endpoints

### Authentication

| Method | Endpoint       | Description              | Auth Required |
|--------|----------------|--------------------------|---------------|
| POST   | `/api/v1/register` | Register a new user     | No            |
| POST   | `/api/v1/login`    | Authenticate user      | No            |

### User Profile

| Method | Endpoint        | Description                 | Auth Required |
|--------|-----------------|-----------------------------|---------------|
| GET    | `/api/v1/profile` | Get authenticated user profile | Yes           |
| PUT    | `/api/v1/profile` | Update authenticated user profile | Yes          |
| DELETE | `/api/v1/profile` | Delete authenticated user account | Yes         |

### Posts

| Method | Endpoint               | Description                | Auth Required |
|--------|------------------------|----------------------------|---------------|
| POST   | `/api/v1/posts`        | Create a new post          | Yes           |
| GET    | `/api/v1/posts`        | Get all posts              | No            |
| GET    | `/api/v1/posts/{id}`   | Get a specific post        | No            |
| PUT    | `/api/v1/posts/{id}`   | Update a specific post     | Yes           |
| DELETE | `/api/v1/posts/{id}`   | Delete a specific post     | Yes           |

### Comments

| Method | Endpoint                          | Description                   | Auth Required |
|--------|----------------------------------|-------------------------------|---------------|
| POST   | `/api/v1/posts/{id}/comments`    | Create a comment on a post   | Yes           |
| GET    | `/api/v1/posts/{id}/comments`    | Get all comments for a post  | No            |
| GET    | `/api/v1/users/{user_id}/comments` | Get all comments by a user   | No            |
| GET    | `/api/v1/comments/{id}`           | Get a specific comment       | No            |
| PUT    | `/api/v1/comments/{id}`           | Update a specific comment    | Yes           |
| DELETE | `/api/v1/comments/{id}`           | Delete a specific comment    | Yes           |

## Example Requests

### Register a new user

```bash
curl -X POST http://localhost:8090/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8090/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Create a post (requires authentication)

```bash
curl -X POST http://localhost:8090/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "title": "My New Post",
    "content": "This is the content of my post."
  }'
```

### Create a comment (requires authentication)

```bash
curl -X POST http://localhost:8090/api/v1/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "content": "This is a great post!"
  }'
```

Note: The route structure is now nested under the post resource, so `/api/v1/posts/{id}/comments` instead of `/api/v1/posts/{post_id}/comments`.

## Response Format

Successful responses follow this format:
```json
{
  "message": "Operation successful",
  "data": { ... },
  "meta": { ... }  // Optional metadata (pagination info, etc.)
}
```

Error responses follow this format:
```json
{
  "error": "Error description",
  "code": 400
}
```

## Database Schema

The application uses three main tables:

- `users`: Stores user information (ID, username, email, password hash)
- `posts`: Stores blog posts (ID, title, content, user_id)
- `comments`: Stores comments (ID, content, post_id, user_id)

## Security Features

- JWT-based authentication
- Password hashing using bcrypt
- Input validation and sanitization
- Authorization checks for protected resources
- SQL injection prevention through ORM

## Testing the API

After starting the server, you can test the API using tools like:

- curl
- Postman
- Insomnia
- A browser-based REST client

The server runs on `http://localhost:8090` by default.