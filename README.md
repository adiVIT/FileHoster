# FileHoster

Here's a detailed README for the file manager project, including information about the environment variables used:

# Go File Manager

Go File Manager is a robust and scalable file management system built with Go. It provides secure file storage, retrieval, and sharing capabilities, along with user authentication and caching for improved performance.

## Features

- User authentication using JWT
- File upload and management
- File search functionality
- File sharing between users
- Redis caching for improved performance
- Integration with AWS S3 for file storage
- PostgreSQL database for user and file metadata storage

## Project Structure

```
go_file_manager/
│
├── cmd/
│   └── server/
│       └── main.go          # Entry point of the application
│
├── internal/
│   ├── auth/
│   │   └── auth.go          # Authentication logic (JWT/OAuth2)
│   │
│   ├── db/
│   │   ├── db.go            # Database connection and queries
│   │   └── models.go        # Database models for users and files
│   │
│   ├── files/
│   │   ├── upload.go        # File upload handling
│   │   ├── manage.go        # File management logic
│   │   └── search.go        # File search functionality
│   │
│   ├── cache/
│   │   └── redis.go         # Redis caching logic
│   │
│   ├── handlers/
│   │   ├── auth_handler.go  # HTTP handlers for authentication
│   │   ├── file_handler.go  # HTTP handlers for file operations
│   │   └── share_handler.go # HTTP handlers for file sharing
│   │
│   └── utils/
│       └── utils.go         # Utility functions (e.g., error handling)
│
├── scripts/
│   └── migrate.sql          # SQL scripts for database schema setup
│
├── config/
│   └── config.go            # Configuration settings (e.g., DB, Redis, S3)
│
├── test/
│   ├── auth_test.go         # Tests for authentication
│   ├── file_test.go         # Tests for file operations
│   └── cache_test.go        # Tests for caching logic
│
├── Dockerfile               # Docker configuration for containerization
├── docker-compose.yml       # Docker Compose file for multi-container setup
├── go.mod                   # Go module file
└── go.sum                   # Go module dependencies
```

## Environment Variables

The application uses the following environment variables for configuration:

```
PORT=8080                    # Application port
DB_HOST=localhost            # PostgreSQL host
DB_PORT=5433                 # PostgreSQL port
DB_USER=postgres             # PostgreSQL user
DB_PASSWORD=1234             # PostgreSQL password
DB_NAME=postgres             # PostgreSQL database name
AWS_REGION=eu-north-1        # AWS region for S3
S3_BUCKET=gobucketfile       # S3 bucket name
REDIS_ADDR=172.21.22.102:6379 # Redis address
REDIS_PASSWORD=1234          # Redis password
JWT_SECRET=magic_pawns       # Secret key for JWT token generation
```


## Setup and Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go_file_manager.git
   cd go_file_manager
   ```

2. Set up the environment variables:
   - Create a `.env` file in the root directory
   - Add the environment variables listed above to the `.env` file

3. Install dependencies:
   ```
   go mod download
   ```

4. Set up the database:
   - Ensure PostgreSQL is running
   - Run the migration script:
     ```
     psql -U postgres -d postgres -f scripts/migrate.sql
     ```

5. Start the Redis server

6. Build and run the application:
   ```
   go build ./cmd/server
   ./server
   ```

DOCKER

1. Pull the Docker image:
   ```
   docker pull adityabajaj22/filehostproject:latest
   ```

2. Run the Docker container:
   ```
   docker run -d --name filehostproject adityabajaj22/filehostproject:latest
   ```
Using Docker Compose
If you prefer using Docker Compose, make sure you have a docker-compose.yml file configured for this image. Then run:
docker-compose up
