package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"todo-app/internal/handlers"
	"todo-app/internal/db"
	"todo-app/internal/middleware"
	"time"
	"todo-app/internal/models"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := db.InitDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := mux.NewRouter()
	
	jwtSecret := []byte("your-secret-key")
	authHandler := handlers.NewAuthHandler(string(jwtSecret))
	taskHandler := handlers.NewTaskHandler()
	notificationHandler := handlers.NewNotificationHandler()

	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")

	taskRouter := r.PathPrefix("/api/tasks").Subrouter()
	taskRouter.Use(middleware.AuthMiddleware(jwtSecret))
	taskRouter.HandleFunc("", taskHandler.Create).Methods("POST", "OPTIONS")
	taskRouter.HandleFunc("", taskHandler.List).Methods("GET", "OPTIONS")
	taskRouter.HandleFunc("/{id}", taskHandler.Update).Methods("PUT", "OPTIONS")
	taskRouter.HandleFunc("/{id}", taskHandler.Delete).Methods("DELETE", "OPTIONS")

	notificationRouter := r.PathPrefix("/api/notifications").Subrouter()
	notificationRouter.Use(middleware.AuthMiddleware(jwtSecret))
	notificationRouter.HandleFunc("", notificationHandler.List).Methods("GET", "OPTIONS")
	notificationRouter.HandleFunc("/{id}/read", notificationHandler.MarkAsRead).Methods("POST", "OPTIONS")
	notificationRouter.HandleFunc("/check", notificationHandler.CheckDueTasks).Methods("POST", "OPTIONS")

	r.Use(corsMiddleware)

	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			if err := models.CheckDueTasks(); err != nil {
				log.Printf("Error checking due tasks: %v", err)
			}
		}
	}()

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
} 