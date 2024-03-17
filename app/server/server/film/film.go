package film

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
	var film models.Film
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := db.Conn.Create(&film)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filmJson, err := json.Marshal(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(filmJson)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var films []models.Film
	var search models.Search
	search.Search = r.URL.Query().Get("search")
	search.Ordering = r.URL.Query().Get("ordering")
	search.Field = r.URL.Query().Get("field")
	search.ActorName = r.URL.Query().Get("actor")
	if search.Ordering == "" {
		search.Ordering = "desc"
	}
	if search.Field == "" {
		search.Field = "rating"
	}
	res := db.Conn.
		Preload("Actors").
		Where("LOWER(name) LIKE LOWER(?)", "%"+search.Search+"%").
		Order(search.Field + " " + search.Ordering).
		Find(&films)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filmsJson, err := json.Marshal(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsJson)
}
