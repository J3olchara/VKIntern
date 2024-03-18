package actor

import (
	"encoding/json"
	"github.com/J3olchara/VKIntern/app/server/db"
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
	var actor *models.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	actor.Create()
	if !actor.Create() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	actorJson, err := json.Marshal(actor)
	support.PanicErr(err)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(actorJson)
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	var actors []models.Actor
	res := db.Conn.Preload("Films").Find(&actors)
	support.PanicErr(res.Error)
	actorsJson, err := json.Marshal(actors)
	support.PanicErr(err)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorsJson)
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
	var actor *models.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(strings.Trim(r.URL.Path, "/actor/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	actor.ID = uint(id)
	if !actor.Save() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	actorJson, err := json.Marshal(actor)
	support.PanicErr(err)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(actorJson)
}

func (h IDHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Trim(r.URL.Path, "/actor/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	actor := &models.Actor{ID: uint(id)}
	if !actor.Delete() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
