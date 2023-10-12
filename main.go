package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tmbrody/blogAggregatorGo/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type contextKey string

const dbContextKey contextKey = "db"

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("CONN")
	if dbURL == "" {
		log.Fatal("CONN environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open a connection to the postgres database")
	}

	dbQueries := database.New(db)

	var apiCfg apiConfig
	apiCfg.DB = dbQueries

	r := chi.NewRouter()
	r_handlers := chi.NewRouter()

	r_handlers.Get("/readiness", readinessHandler)
	r_handlers.Get("/err", errorHandler)

	r_handlers.Post("/users", withDB(createUserHandler, apiCfg.DB))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/v1", r_handlers)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Serving files on port: %s", port)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func withDB(next http.HandlerFunc, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dbContextKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
