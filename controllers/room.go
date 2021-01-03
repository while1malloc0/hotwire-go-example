package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/while1malloc0/hotwire-go-example/models"
	"github.com/while1malloc0/hotwire-go-example/pkg/notice"
)

type contextKey struct{}

var (
	// ContextKeyRoom is a type safe representation of the key "room" inside of a context.Context
	ContextKeyRoom = contextKey{}
)

// RoomsController implements Controller functionality for the Room model
type RoomsController struct{}

// Context is a middleware that parses the Room ID from a request, loads the
// corresponding Room model, and makes it available as part of the request's
// context
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

// Index shows a list of all Rooms
func (*RoomsController) Index(w http.ResponseWriter, r *http.Request) {
	rooms, err := models.ListRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	notice := r.Context().Value(notice.ContextKey)
	responseData := map[string]interface{}{"Rooms": rooms, "Notice": notice}
	render.HTML(w, http.StatusOK, "rooms/index", responseData)
}

// New renders a form for creating a new Room
func (*RoomsController) New(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, http.StatusOK, "rooms/new", nil)
}

// Edit renders a form for making changes to an existing Room
func (*RoomsController) Edit(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/edit", responseData)
}

// Get retrieves a single Room by its ID
func (*RoomsController) Get(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	responseData := map[string]interface{}{"Room": room}
	render.HTML(w, http.StatusOK, "rooms/show", responseData)
}

// Update makes changes to a Room given its ID
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
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/rooms/%d", room.ID), http.StatusFound)
}

// Create makes a new Room
func (*RoomsController) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("room[name]")
	err = models.CreateRoom(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notice.Set(w, "Room created successfully")
	http.Redirect(w, r, "/rooms", http.StatusFound)
}

// Destroy deletes an existing room
func (*RoomsController) Destroy(w http.ResponseWriter, r *http.Request) {
	room := r.Context().Value(ContextKeyRoom).(*models.Room)
	err := models.DeleteRoom(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notice.Set(w, "Room deleted successfully")
	http.Redirect(w, r, "/rooms", http.StatusFound)
}
