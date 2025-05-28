package routes

import (
	"log-service/internal/handlers"

	"github.com/gorilla/mux"
)

func UserLogRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/log", handlers.CreateLog).Methods("POST")
	r.HandleFunc("/user/logs", handlers.GetLogs).Methods("GET")

	return r
}
