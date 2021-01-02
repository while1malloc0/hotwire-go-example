package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Print("Starting app")

	log.Print("Running migrations")
	db, err := gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
	db.AutoMigrate(&models.Room{}, &models.Message{})

	result := db.Raw("TRUNCATE TABLE rooms;")
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	db.Create(&models.Room{Name: "Test Room"})

	log.Print("Registering routes")
	router := registerRoutes()

	log.Print("Serving on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes() http.Handler {
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
