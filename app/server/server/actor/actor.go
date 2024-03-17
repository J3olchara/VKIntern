package actor

import (
	"encoding/json"
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"log"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.post(w, r)
	} else if r.Method == "GET" {
		h.get(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var actor *models.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := db.Conn.Create(actor)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actorJson, err := json.Marshal(actor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(actorJson)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var actors []models.Actor
	res := db.Conn.Find(&actors)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(res.Error)
	}
	actorsJson, err := json.Marshal(actors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorsJson)
}
