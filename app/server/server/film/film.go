package film

import (
	"encoding/json"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/J3olchara/VKIntern/app/server/support"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
}

type IDHandler struct {
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	if r.Method == http.MethodPost && user.Staff {
		h.post(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h Handler) post(w http.ResponseWriter, r *http.Request) {
	var film models.Film
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !film.Create() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filmJson, err := json.Marshal(film)
	support.PanicErr(err)
	models.CreateRelations(&film, film.Actors)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(filmJson)
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	var search models.Search
	search.Search = r.URL.Query().Get("search")
	search.Ordering = r.URL.Query().Get("ordering")
	search.Field = r.URL.Query().Get("field")
	search.ActorName = r.URL.Query().Get("actor")
	films := models.SearchFilms(search)
	filmsJson, err := json.Marshal(films)
	support.PanicErr(err)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsJson)
}

func (h IDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	if r.Method == http.MethodPut && user.Staff {
		h.put(w, r)
		return
	}
	if r.Method == http.MethodDelete && user.Staff {
		h.delete(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h IDHandler) put(w http.ResponseWriter, r *http.Request) {
	var film *models.Film
	intID, err := strconv.Atoi(strings.Trim(r.URL.Path, "/actor/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !film.Save() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	film.ID = uint(intID)
	models.UpdateRelations(film, film.Actors)
	filmJson, err := json.Marshal(film)
	support.PanicErr(err)
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(filmJson)
}

func (h IDHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Trim(r.URL.Path, "/film/"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	film := &models.Film{ID: uint(id)}
	if !film.Delete() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
