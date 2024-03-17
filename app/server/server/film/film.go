package film

import (
	"encoding/json"
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
}

type IDHandler struct {
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
	log.Println(film.Actors)
	res := db.Conn.Omit("id").Create(&film)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filmJson, err := json.Marshal(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	for _, actor := range film.Actors {
		db.Conn.Create(&models.FilmActor{
			FilmID:  film.ID,
			ActorID: actor.ID,
		})
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
		Joins("LEFT JOIN film_actors ON films.id = film_actors.film_id").
		Joins("LEFT JOIN actors ON film_actors.actor_id = actors.id").
		Where("LOWER(films.name) LIKE LOWER(?)", "%"+search.Search+"%")
	if search.ActorName != "" {
		res = res.Where("LOWER(actors.name) LIKE LOWER(?)", "%"+search.ActorName+"%")
	}
	res = res.
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

func (h *IDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		h.put(w, r)
	} else if r.Method == "DELETE" {
		h.delete(w, r)
	}
}

func (h *IDHandler) put(w http.ResponseWriter, r *http.Request) {
	intID, err := strconv.Atoi(strings.Trim(r.URL.Path, "/actor/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := uint(intID)
	var film models.Film
	if err = json.NewDecoder(r.Body).Decode(&film); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	film.ID = id
	res := db.Conn.Save(&film)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(res.Error)
	}
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	filmJson, err := json.Marshal(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusCreated)

	for _, actor := range film.Actors {
		db.Conn.
			Set("gorm:insert_option", "ON CONFLICT DO NOTHING").
			Create(&models.FilmActor{
				FilmID:  film.ID,
				ActorID: actor.ID,
			})
	}
	_, _ = w.Write(filmJson)
}

func (h *IDHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Trim(r.URL.Path, "/film/"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res := db.Conn.Delete(&models.Film{}, uint(id))
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(res.Error)
	}
	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
