package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
)

type contextKey struct{}

var ContextKeyRoom = contextKey{}

type RoomsController struct{}

func (*RoomsController) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var room *models.Room
		var err error
		if idStr := chi.URLParam(r, "id"); idStr != "" {
			room, err = roomFromIDString(idStr)
		}
		if models.IsRecordNotFound(err) {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Fatal error: %v", err), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyRoom, room)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func roomFromIDString(idStr string) (*models.Room, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, err
	}
	room, err := models.FindRoom(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (*RoomsController) Index(w http.ResponseWriter, r *http.Request) {
	rooms, err := models.ListRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseData := map[string]interface{}{"Rooms": rooms}
	render.HTML(w, http.StatusOK, "rooms/index", responseData)
}

func (*RoomsController) Edit(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/edit", responseData)
}

func (*RoomsController) Get(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/show", responseData)
}

func (*RoomsController) Update(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("room[name]")
	updates := map[string]interface{}{
		"Name": name,
	}
	err = models.UpdateRoom(room, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/rooms/%s", room.ID), http.StatusFound)
}
