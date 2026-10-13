package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/thedeepak12/argus/internal/auth"
	"github.com/thedeepak12/argus/internal/db"
	"github.com/thedeepak12/argus/internal/handlers"
)

func main() {
	ctx := context.Background()
	pool := db.NewPool(ctx)
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	queries := db.New(sqlDB)

	authHandler := handlers.NewAuthHandler(queries)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)
	})

	r.Route("/api/protected", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			userID := auth.GetUserID(r)
			w.Write([]byte("Hello, user " + userID))
		})
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
