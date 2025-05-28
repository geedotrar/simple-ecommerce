# Product Service (Go + Gin + GORM)

This is a microservice built with Go that handles product management. It supports full CRUD operations and image uploads to AWS S3. Built using the Gin web framework and GORM ORM with PostgreSQL as the database.

## Features

- Product CRUD (Create, Read, Update, Delete)
- Upload product images to AWS S3
    - Fallback to local file storage if S3 upload fails or is not configured
- Input validation using go-playground/validator
- CORS support for frontend-backend communication
- Clean modular project structure

## Tech Stack
- Language : Go (v1.23.1)
- Framework : Gin
- ORM : GORM
- Database : PostgreSQL
- Cloud Storage AWS S3
- Env Config : godotenv
- Cache : Redis

## 🚀 Setup
1. **Navigate to the Product Service directory**
   ```
   cd ags-microservices/backend/microservices/product-service
   ```
2. **Copy the example environment file**
   ```
   cp .env.example .env
   ```
3. **Fill in your .env file**

   Example: 
   ```env
    # Server Configuration
    SERVER_PORT=8080

    # PostgreSQL Database Configuration
    DATABASE_URL=host=localhost user=postgres password=password dbname=ags_products_db port=5432 sslmode=disable

    # Local File Upload Directory
    UPLOAD_DIR=uploads

    # Auth Service Access URL
    AUTH_ACCESS_URL=http://localhost:8000/api/access

    # Redis Configuration
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=

    # AWS S3 Configuration
    AWS_ACCESS_KEY_ID=your-access-key-id
    AWS_SECRET_ACCESS_KEY=your-secret-access-key
    AWS_REGION=us-east-1
    AWS_S3_BUCKET=ags-products-bucket
    AWS_S3_FOLDER=ags-products-image
   ```

3. **Install Go dependencies**
   ```
   go mod tidy
   ```
4. **Run database migrations and seeders**
   ```
   make migrate-up
   ```
5. **Start the service**
   ```
   go run cmd/main.go
   ```

