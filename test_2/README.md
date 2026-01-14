# Blog REST API

A complete RESTful API for a blog system with user authentication, posts, and comments.

## Features

- User authentication and authorization (JWT-based)
- CRUD operations for posts and comments
- Input validation and error responses
- Database integration with PostgreSQL
- Transaction support for database operations
- Comprehensive error handling
- Docker deployment support
- OpenAPI/Swagger documentation
- Health checking capabilities

## Architecture

The application follows a clean architecture pattern with separation of concerns:

```
cmd/api/main.go                    # Entry point with routing configuration
├── internal/
│   ├── database/                 # Database connection and initialization
│   ├── middleware/              # Authentication middleware (JWT)
│   ├── models/                  # Data models (GORM structs)
│   ├── responses/               # Standardized response formats
│   ├── routes/                  # Route definitions
│   ├── users/                   # User management (Controller, Service, Repository)
│   ├── posts/                   # Post management (Controller, Service, Repository)
│   └── comments/                # Comment management (Controller, Service, Repository)
```

### Layered Architecture

- **Controller Layer**: HTTP request/response handling
- **Service Layer**: Business logic processing
- **Repository Layer**: Database operations
- **Models**: Data structures and relations
- **Middleware**: Cross-cutting concerns (auth, logging, etc.)

## Prerequisites

- Go 1.19+
- PostgreSQL database
- Environment variables configured
- Docker and Docker Compose (for containerized deployment)

## Setup

### Local Development Setup

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Configure environment variables (see below)
4. Run the seeding script: `go run cmd/seed/main.go`
5. Start the server: `go run cmd/api/main.go`

### Using Docker

1. Build and run with docker-compose:
   ```bash
   docker-compose up --build
   ```

2. The API will be available at `http://localhost:8090`
3. The database will be available at `localhost:5432`

### Environment Variables

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=majoo_test
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
```

For Docker deployment, these are automatically loaded from the docker-compose.yml file.

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

## API Documentation

The API includes OpenAPI/Swagger documentation available at:
- JSON format: `http://localhost:8090/swagger/doc.json`
- Interactive UI: `http://localhost:8090/swagger/index.html`

Documentation is generated using [Swag](https://github.com/swaggo/swag) and available in the `docs/` directory.

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

Validation error responses follow this format:
```json
{
  "error": "Validation failed",
  "code": 422,
  "errors": [
    {
      "field": "username",
      "message": "This field is required",
      "value": ""
    }
  ]
}
```

## Database Schema

Detailed database schema documentation is available in [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md).

## Security Features

- JWT-based authentication with configurable expiration
- Password hashing using bcrypt with salt
- Input validation and sanitization using Go validators
- Authorization checks for protected resources
- SQL injection prevention through GORM ORM
- Token signature validation with HMAC signing method
- Secure token storage (tokens are not stored server-side)

## Transaction Support

The application implements transaction support for complex operations:
- Repository layer supports transaction contexts
- Service layer supports WithTransaction functionality
- Database operations can be wrapped in ACID transactions
- Example: Creating a user with additional profile data in a single transaction

## Testing the API

After starting the server, you can test the API using tools like:

- curl
- Postman
- Insomnia
- A browser-based REST client
- Automated tests

The server runs on `http://localhost:8090` by default.

## Docker Deployment

The application can be deployed using Docker:

### Building the image
```bash
docker build -t blog-api .
```

### Running with Docker Compose (recommended)
```bash
docker-compose up --build
```

### Running the container
```bash
docker run -p 8090:8090 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e DB_NAME=your-db-name \
  -e JWT_SECRET_KEY=your-secret-key \
  blog-api
```

## Known Limitations

- **Performance**: For very high-load scenarios, consider caching with Redis
- **File Uploads**: The current API doesn't support file uploads (images, documents)
- **Real-time Features**: No WebSocket support for real-time notifications or chat
- **Advanced Queries**: Limited search and filtering capabilities
- **Rate Limiting**: No built-in rate limiting (should be implemented at proxy level)
- **Analytics**: No built-in analytics or monitoring
- **Soft Deletes**: Only users and posts support soft deletes; comments are hard deleted

## Future Improvements

- Add Redis caching for better performance
- Implement rate limiting middleware
- Add comprehensive logging with log levels
- Add unit and integration tests
- Implement file upload functionality
- Add email notifications for comments
- Implement advanced search and filtering
- Add real-time features with WebSockets
- Add metrics and monitoring (Prometheus integration)

## Docker Files

- `Dockerfile`: Production-ready container image
- `docker-compose.yml`: Multi-service orchestration (API + Database)
- `generate_swagger.sh`: Script to generate API documentation

## Git Best Practices

- Follow semantic commit messages
- Use feature branches for new functionality
- Ensure all tests pass before merging
- Update documentation when adding new features
- Keep PRs small and focused

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes and commit (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request