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

var ContextKeyRoomID = contextKey{}

type RoomsController struct{}

func (*RoomsController) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")
		if idStr != "" {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			ctx = context.WithValue(r.Context(), ContextKeyRoomID, id)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	roomID := r.Context().Value(ContextKeyRoomID).(uint64)
	room, err := models.FindRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/edit", responseData)
}

func (*RoomsController) Get(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.Context().Value(ContextKeyRoomID).(uint64)
	if !ok {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	room, err := models.FindRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/show", responseData)
}

func (*RoomsController) Update(w http.ResponseWriter, r *http.Request) {
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
