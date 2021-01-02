package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/pkg/view"
	"nhooyr.io/websocket"
)

var messageSocketChan = make(chan []byte)

type NewMessageResponseData struct {
	Room *models.Room
}

type CreateMessageResponseData struct {
	Message *models.Message
}

func NewMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	room, err := models.FindRoom(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	content, err := view.Render("messages/new", NewMessageResponseData{Room: room})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message := &models.Message{
		Content: r.FormValue("message[content]"),
		RoomID:  id,
	}

	err = models.CreateMessage(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html; turbo-stream; charset=utf-8")
	content, err := view.Render("messages/create.turbostream", CreateMessageResponseData{Message: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	messageSocketChan <- content
	w.WriteHeader(201)
}

func MessageSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close(websocket.StatusInternalError, "")

	for {
		select {
		case content := <-messageSocketChan:
			ws.Write(context.TODO(), websocket.MessageText, content)
		}
	}
}
