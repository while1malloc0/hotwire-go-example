package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/routes"
)

func main() {
	log.Println("Starting app")

	setupDB()

	router := registerRoutes()

	port := getPort()
	log.Printf("Serving on port %d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes() http.Handler {
	log.Println("Registering routes")
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	routes.Register(r)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	logRoutes(r)

	return r
}

func logRoutes(r chi.Router) {
	log.Println("Serving with routes")
	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Println(method, route)
		return nil
	})
}

func setupDB() {
	log.Println("Running migrations")
	err := models.Migrate()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	log.Println("Seeding database")
	err = models.Seed()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func getPort() int {
	port := 8080
	if p, ok := os.LookupEnv("HOTWIRE_CHAT_PORT"); ok {
		parsed, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Fatal error: %v", err)
		}
		port = parsed
	}
	return port
}
