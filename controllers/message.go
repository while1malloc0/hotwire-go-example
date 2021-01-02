package controllers

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
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
	render.HTML(w, http.StatusOK, "messages/new", NewMessageResponseData{Room: room})
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
	var content bytes.Buffer
	err = render.HTML(&content, http.StatusCreated, "messages/create.turbostream", CreateMessageResponseData{Message: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	messageSocketChan <- content.Bytes()
}

var started bool
var wss []*websocket.Conn

func MessageSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wss = append(wss, ws)

	if !started {
		go listenWs()
		started = true
	}
}

func listenWs() {
	for {
		select {
		case content := <-messageSocketChan:
			for i := range wss {
				wss[i].Write(context.TODO(), websocket.MessageText, content)
			}
		}
	}
}
