package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/while1malloc0/hotwire-go-example/controllers"
	"github.com/while1malloc0/hotwire-go-example/models"
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

	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/messages/socket", controllers.MessageSocket)

	r.Get("/", controllers.RoomsIndex)
	r.Get("/rooms", controllers.RoomsIndex)
	r.Get("/rooms/{id}", controllers.GetRoom)
	r.Post("/rooms/{id}", controllers.UpdateRoom)
	r.Get("/rooms/{id}/edit", controllers.EditRoom)
	r.Get("/rooms/{id}/messages/new", controllers.NewMessage)
	r.Post("/rooms/{id}/messages", controllers.CreateMessage)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	log.Print("Serving on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
