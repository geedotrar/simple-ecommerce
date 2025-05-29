package middleware

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"product-service/config"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		accessURL := os.Getenv("USER_AUTH_ACCESS_URL")
		if accessURL == "" {
			log.Println("[AuthMiddleware] USER_AUTH_ACCESS_URL not configured")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "USER_AUTH_ACCESS_URL not configured"})
			return
		}

		req, err := http.NewRequestWithContext(context.Background(), "GET", accessURL, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth request"})
			return
		}
		req.Header.Set("Authorization", authHeader)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[AuthMiddleware] Auth request error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("[AuthMiddleware] Auth service responded with status: %d\n", resp.StatusCode)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[AuthMiddleware] Failed to read auth response: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read auth response"})
			return
		}

		var respData struct {
			User UserClaims `json:"user"`
		}

		if err := json.Unmarshal(body, &respData); err != nil {
			log.Printf("[AuthMiddleware] Invalid auth response: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth response"})
			return
		}

		ctx := context.Background()
		redisKey := "laravel_database_role:" + respData.User.Role

		log.Printf("[AuthMiddleware] Fetching permissions from Redis key: %s\n", redisKey)

		permJson, err := config.Client.Get(ctx, redisKey).Result()
		if err != nil {
			log.Printf("[AuthMiddleware] Redis get error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found or Redis error"})
			return
		}

		var permissions []string
		if err := json.Unmarshal([]byte(permJson), &permissions); err != nil {
			log.Printf("[AuthMiddleware] Failed to parse permissions JSON: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse permissions from Redis"})
			return
		}

		claims := struct {
			UserClaims
			Permissions []string
		}{
			UserClaims:  respData.User,
			Permissions: permissions,
		}

		c.Set("claims", claims)
		c.Next()
	}
}
