package controllers

import (
	"bytes"
	"context"
	"net/http"

	"github.com/while1malloc0/hotwire-go-example/models"
	"nhooyr.io/websocket"
)

var messageSocketChan = make(chan []byte)

type MessagesController struct{}

func (*MessagesController) New(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ContextKeyRoomID).(uint64)
	room, err := models.FindRoom(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "messages/new", responseData)
}

func (*MessagesController) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	roomID := r.Context().Value(ContextKeyRoomID).(uint64)
	message := &models.Message{
		Content: r.FormValue("message[content]"),
		RoomID:  int(roomID),
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

	messageSocketChan <- content.Bytes()
}

var started bool
var wss []*websocket.Conn

func (*MessagesController) Socket(w http.ResponseWriter, r *http.Request) {
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
