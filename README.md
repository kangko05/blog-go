# blog-go

A modern personal blog system built with Go and React, featuring a powerful CLI for content management and a clean web interface.

## Features

- **Clean Architecture**: Well-structured Go backend with repository pattern
- **CLI Management**: Full-featured command-line interface for post management
- **Markdown Support**: Write posts in Markdown with HTML rendering
- **Category System**: Organize posts as Notes or Projects
- **REST API**: Clean API endpoints for frontend integration
- **Authentication**: JWT-based auth system (optional)
- **Database**: SQLite for simplicity and portability

## Tech Stack

### Backend

- **Go 1.25** - Modern Go with latest features
- **Gin** - Fast HTTP web framework
- **SQLite** - Lightweight, embedded database
- **Cobra** - Powerful CLI framework
- **gomarkdown** - Markdown to HTML conversion
- **JWT** - Token-based authentication

## Quick Start

### Prerequisites

- Go 1.25 or higher

### Backend Setup

1. **Clone and build**

   ```bash
   git clone <repository-url>
   cd blog-go
   go mod tidy
   ```

2. **Environment Configuration**
   Create a `.env` file:

   ```env
   DBPATH=./dev.db
   SERVERPORT=:8000
   JWTSECRET=your-secret-key-here
   ```

3. **Build the applications**

   ```bash
   # Build CLI tool
   go build -o cli ./cmd/cli

   # Build server
   go build -o server ./cmd/server
   ```

4. **Start the server**
   ```bash
   ./server
   ```

## CLI Usage

The CLI tool provides complete post management functionality:

### Creating Posts

```bash
# Create a new note
./cli create --notes

# Create a new project post
./cli create --proj
```

This opens your default editor (set `EDITOR` environment variable). The first `# Heading` becomes the post title.

### Managing Posts

```bash
# List all posts
./cli list

# List only notes
./cli list --notes

# List only projects
./cli list --proj

# Show a specific post
./cli show <id>

# Update a post
./cli update <id>

# Delete a post (with confirmation)
./cli delete <id>
```

### Import/Export

```bash
# Import markdown file as a post
./cli import <filepath> --notes
./cli import <filepath> --proj

# Export post to markdown file
./cli export <id> <filepath>
```

## API Endpoints

### Public Endpoints

- `GET /notes` - List all notes
- `GET /projects` - List all projects
- `GET /posts/:id` - Get specific post
- `GET /checkhealth` - Health check

### Development Endpoints (commented out)

- `POST /posts` - Create post (requires auth)
- `PUT /posts/:id` - Update post (requires auth)
- `DELETE /posts/:id` - Delete post (requires auth)
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout

## Project Structure

```
.
├── cmd/
│   ├── cli/main.go          # CLI entry point
│   └── server/main.go       # Server entry point
├── internal/
│   ├── auth/                # Authentication logic
│   ├── cli/                 # CLI commands
│   ├── config/              # Configuration management
│   ├── post/                # Post domain logic
│   ├── repo/                # Database repositories
│   └── server/              # HTTP handlers & middleware
├── test/                    # Test markdown files
├── docs/                    # Documentation
└── logs/                    # Access logs
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/post
```

### Database Schema

The application automatically creates SQLite tables:

**Posts Table:**

- `id` - Primary key
- `title` - Post title
- `content` - Markdown content
- `category` - "notes" or "proj"
- `created_at` - Creation timestamp
- `updated_at` - Last modification timestamp

**Users Table (for auth):**

- `id` - Primary key
- `username` - Unique username
- `hashed_password` - Bcrypt hashed password

### Middleware Features

- **CORS**: Configured for development origins
- **Rate Limiting**: 10 requests per second, burst of 30
- **Security Headers**: XSS protection, HSTS, etc.
- **Access Logging**: JSON-formatted request logs
- **JWT Authentication**: Token-based auth for protected routes

## Configuration

### Environment Variables

| Variable     | Description               | Default  |
| ------------ | ------------------------- | -------- |
| `DBPATH`     | SQLite database file path | Required |
| `SERVERPORT` | Server listen address     | Required |
| `JWTSECRET`  | JWT signing secret        | Optional |
| `EDITOR`     | Preferred editor for CLI  | `nvim`   |

## Deployment

### Building for Production

```bash
# Build optimized binaries
CGO_ENABLED=1 go build -ldflags="-w -s" -o blog-server ./cmd/server
CGO_ENABLED=1 go build -ldflags="-w -s" -o blog-cli ./cmd/cli
```

### Production Environment

1. Set production environment variables
2. Use a process manager (systemd, PM2, etc.)
3. Configure reverse proxy (nginx, caddy)
4. Set up HTTPS and proper CORS origins
5. Enable authentication endpoints if needed

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

Future features planned as the blog grows and more content is added:

- [ ] Search functionality
- [ ] Tag system
- [ ] Image upload support
- [ ] Pagination for post lists
- [ ] Static site generation
