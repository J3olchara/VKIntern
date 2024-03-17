package server

import (
	"fmt"
	"github.com/J3olchara/VKIntern/app/server/server/actor"
	"github.com/J3olchara/VKIntern/app/server/server/film"
	"github.com/J3olchara/VKIntern/app/server/server/middleware"
	"log"
	"net/http"
	"os"
)

func ApplyRoutes(mux *http.ServeMux) {
	mux.Handle("/film", middleware.RequestLogger(&film.Handler{}))
	mux.Handle("/actor", middleware.RequestLogger(&actor.Handler{}))
}

func StartServer() {
	mux := http.NewServeMux()
	ApplyRoutes(mux)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
