package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"log-service/config"
	"log-service/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.InitMongo()

	r := routes.UserLogRouter()

	port := os.Getenv("PORT")

	fmt.Printf("🚀 Log Service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
