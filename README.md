# Test Task Kami

This project is designed to manage room reservations. The project includes an API for creating and retrieving reservations, built with Go.

## Prerequisites

Before starting, make sure you have the following installed:

- Docker
- Docker Compose
- Make
- Goose (for migrations https://github.com/pressly/goose)

## Getting Started

To set up and run the project, follow these steps:

### 1. Build the Project

First, build the Docker images:

```bash
make docker-build
```

Second, run docker compose
```bash
make compose-up
```

Third, run migration to create the schema
```bash
make migration-up
```

App will started on `localhost:8080`

### 2. Run Tests

```bash
make test
```

### 3. Was used in project

- godotenv - to read config file - https://github.com/joho/godotenv
- zap - to write logs - https://go.uber.org/zap
- pgx - to interact with database - https://github.com/jackc/pgx
- squirrel - to construct queries - https://github.com/Masterminds/squirrel
- validator - to validate input - http://github.com/go-playground/validator/v10
- chi - http router - https://github.com/go-chi/chi/v5
- testify - for tests - https://github.com/stretchr/testify
- swaggo - for swagger - https://github.com/swaggo/swag


### 4. Database schema 

