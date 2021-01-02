package routes

import (
	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/controllers"
)

func Register(r chi.Router) {
	r.Get("/", controllers.RoomsIndex)
	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", controllers.RoomsIndex)
		r.Get("/{id}", controllers.GetRoom)
		r.Get("/{id}/edit", controllers.EditRoom)
		r.Post("/{id}", controllers.UpdateRoom)

		r.Route("/{id}/messages", func(r chi.Router) {
			r.Get("/new", controllers.NewMessage)
			r.Post("/", controllers.CreateMessage)
		})
	})

	r.Route("/messages", func(r chi.Router) {
		r.Get("/socket", controllers.MessageSocket)
	})
}
