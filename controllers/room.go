package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
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

	render.HTML(w, http.StatusOK, "rooms/index", RoomsIndexResponse{Rooms: rooms})
}

func EditRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room, err := models.FindRoom(roomID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.HTML(w, http.StatusOK, "rooms/edit", EditRoomResponse{Room: room})
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room, err := models.FindRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.HTML(w, http.StatusOK, "rooms/show", GetRoomResponse{Room: room})
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
