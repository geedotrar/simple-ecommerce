# Auth Service (Laravel)

This is a microservice built with Laravel that handles user authentication and authorization using JWT (JSON Web Tokens). It supports role-based access to control permissions for different user types.

## Features

- RBAC (Role-Based Access Control):
- Login (Admin | User)
- Register

## Tech Stack
- Framework: Laravel (version 10.10 up to < 11.0)
- Language: PHP 8.1, 8.2, or higher
- Database: PostgreSQL
- Authentication: JWT (using tymon/jwt-auth package)
- Cache : Redis

## 🚀 Setup
1. **Navigate to the Auth Service directory**
   ```
   cd ags-microservices/backend/microservices/auth-service  
   ```
2. **Copy the environment configuration file bash**
   ```
   cp .env.example .env
   ```
3. **Install Dependencies**
   ```
   composer install
   ```
4. **Generate JWT Secret**
   ```
   php artisan jwt:secret
   ```

5. **Run database migrations**
   ```
   php artisan migrate
   ```
6. **Seed database with dummy users**
   ```
   php artisan db:seed
   ```

### 👥 Dummy Users

**ADMIN:**
```json
{
    "email": "admin@ags.com",
    "password": "password"
}
```

**USER:**
```json
{
    "email": "user@ags.com",
    "password": "password"
}
```

7. **Run the Service**
   ```
   php artisan serve
   ```

