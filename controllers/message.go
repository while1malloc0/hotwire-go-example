package controllers

import (
	"bytes"
	"net/http"

	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/pkg/pubsub"
	"nhooyr.io/websocket"
)

type MessagesController struct{}

func (*MessagesController) New(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "messages/new", responseData)
}

func (*MessagesController) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	message := &models.Message{
		Content: r.FormValue("message[content]"),
		Room:    *room,
	}

	err = models.CreateMessage(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html; turbo-stream; charset=utf-8")
	var content bytes.Buffer
	responseData := map[string]interface{}{"Message": message}
	err = render.HTML(&content, http.StatusCreated, "messages/create.turbostream", responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	pubsub.Broadcast(room.ID, content.Bytes())
}

func (*MessagesController) Socket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	pubsub.Subscribe(room.ID, ws)
}
