package routes

import (
	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/controllers"
	"github.com/while1malloc0/hotwire-go-example/pkg/notice"
)

func Register(r chi.Router) {
	roomsController := &controllers.RoomsController{}
	messagesController := &controllers.MessagesController{}

	r.Use(notice.Context)
	r.Get("/", roomsController.Index)

	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", roomsController.Index)
		r.Get("/new", roomsController.New)
		r.Post("/", roomsController.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(roomsController.Context)

			r.Get("/", roomsController.Get)
			r.Post("/", roomsController.Update)

			r.Get("/edit", roomsController.Edit)
			r.Get("/destroy", roomsController.Destroy)

			r.Route("/messages", func(r chi.Router) {
				r.Post("/", messagesController.Create)

				r.Get("/new", messagesController.New)
				r.Get("/socket", messagesController.Socket)
			})
		})
	})
}
