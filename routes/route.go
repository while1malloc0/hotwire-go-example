// Package routes wires up request paths to their Controllers
package routes

import (
	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/controllers"
	"github.com/while1malloc0/hotwire-go-example/pkg/notice"
)

// Register wires up request paths to controllers for the given router
func Register(r chi.Router) {
	roomsController := &controllers.RoomsController{}
	messagesController := &controllers.MessagesController{}

	// parse notices out of cookies if they exist. Useful for one-off storage
	// between requests.
	r.Use(notice.Context)

	// Root route, i.e. /
	r.Get("/", roomsController.Index)

	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", roomsController.Index)   // GET /rooms
		r.Get("/new", roomsController.New)  // GET /rooms/new
		r.Post("/", roomsController.Create) // POST /rooms

		r.Route("/{id}", func(r chi.Router) {
			// parse a room out of the {id} param
			r.Use(roomsController.Context)

			r.Get("/", roomsController.Get)     // GET /rooms/{id}
			r.Post("/", roomsController.Update) // POST /rooms/{id}

			r.Get("/edit", roomsController.Edit)       // GET /rooms/{id}/edit
			r.Get("/destroy", roomsController.Destroy) // GET /rooms/{id}/destroy

			r.Route("/messages", func(r chi.Router) {
				r.Post("/", messagesController.Create) // POST /rooms/{id}/messages

				r.Get("/new", messagesController.New)       // GET /rooms/{id}/messages/new
				r.Get("/socket", messagesController.Socket) // GET /rooms/{id}/messages/socket
			})
		})
	})
}
