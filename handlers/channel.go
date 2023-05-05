package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dg/acordia/middleware"
	"github.com/dg/acordia/models"
	"github.com/dg/acordia/repository"
	"github.com/dg/acordia/responses"
	"github.com/dg/acordia/server"
	"github.com/gorilla/mux"
)

type InsertChannelRequest struct {
	Color               string `bson:"color" json:"color"`
	Background          string `bson:"background" json:"background"`
	DesertRefBackground string `bson:"desert_ref_background" json:"desert_ref_background"`
	Image               string `bson:"image" json:"image"`
	DesertRefImage      string `bson:"desert_ref_image" json:"desert_ref_image"`
	Description         string `bson:"description" json:"description"`
	Name                string `bson:"name" json:"name"`
}

type InsertMessageRequest struct {
	Description string `bson:"description" json:"description"`
	Image       string `bson:"image" json:"image"`
	DesertRef   string `bson:"desert_ref" json:"desert_ref"`
}

func CreateChannelHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		profile, _ := middleware.ValidateToken(s, w, r)
		// Handle request
		w.Header().Set("Content-Type", "application/json")
		var req = InsertChannelRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			responses.BadRequest(w, "Invalid request")
			return
		}
		users := []models.Profile{*profile}
		channel := models.InsertChannel{
			Users:               users,
			Color:               req.Color,
			Background:          req.Background,
			DesertRefBackground: req.Background,
			Image:               req.Image,
			DesertRefImage:      req.DesertRefImage,
			Description:         req.Description,
			Name:                req.Name,
			Messages:            []models.ChannelMessage{},
		}
		insertChannel, err := repository.CreateChannel(r.Context(), channel)
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertChannel)
	}
}

func ListOfChannelsHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		profile, _ := middleware.ValidateToken(s, w, r)
		listChannels, err := repository.ListOfChannels(r.Context(), profile.Id)
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(listChannels)
	}
}

func UpdateChannelHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		_, _ = middleware.ValidateToken(s, w, r)
		// Handle request
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		var req = models.UpdateChannel{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			responses.BadRequest(w, "Invalid request")
			return
		}
		updateChannel, err := repository.UpdateChannel(r.Context(), params["id"], req)
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(updateChannel)
	}
}

func DeleteChannelHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		_, _ = middleware.ValidateToken(s, w, r)
		// Handle request
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		err := repository.DeleteChannel(r.Context(), params["id"])
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		responses.DeleteResponse(w, "Channel deleted")
	}
}

func AddUserToChannelHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		_, _ = middleware.ValidateToken(s, w, r)
		// Handle request
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		channel, err := repository.AddUserToChannel(r.Context(), params["user"], params["id"])
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(channel)
	}
}

func AddMessagesToChannelHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		profile, _ := middleware.ValidateToken(s, w, r)
		// Handle request
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		var req = InsertMessageRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			responses.BadRequest(w, "Invalid request")
			return
		}
		loc, error := time.LoadLocation("America/Bogota")
		if error != nil {
			responses.InternalServerError(w, "Error loading location")
			return
		}
		message := models.ChannelMessage{
			User:        *profile,
			Date:        time.Now().In(loc).Format("2006-01-02 15:04:05"),
			Description: req.Description,
			Image:       req.Image,
			DesertRef:   req.DesertRef,
		}
		insertMessage, err := repository.AddMessagesToChannel(r.Context(), &message, params["id"])
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertMessage)
	}
}

func RemoveUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Token validation
		_, _ = middleware.ValidateToken(s, w, r)
		// Handle request
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		removeUser, err := repository.RemoveUser(r.Context(), params["id"], params["user"])
		if err != nil {
			responses.InternalServerError(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(removeUser)
	}
}
