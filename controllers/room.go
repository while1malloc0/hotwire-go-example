package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/pkg/view"
)

type RoomsIndexResponse struct {
	Rooms []*models.Room
}

type EditRoomResponse struct {
	Room *models.Room
}

type GetRoomResponse struct {
	Room *models.Room
}

func RoomsIndex(w http.ResponseWriter, r *http.Request) {
	rooms, err := models.ListRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := view.Render("rooms/index", RoomsIndexResponse{Rooms: rooms})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

func EditRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room, err := models.FindRoom(roomID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := view.Render("rooms/edit", EditRoomResponse{Room: room})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room, err := models.FindRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := view.Render("rooms/show", GetRoomResponse{Room: room})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(200)
}

func UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("room[name]")
	updates := map[string]interface{}{
		"Name": name,
	}
	err = models.UpdateRoom(roomID, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/rooms/%s", roomID), http.StatusFound)
}
