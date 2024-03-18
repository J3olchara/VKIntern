package middleware

import (
	"context"
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/J3olchara/VKIntern/app/server/support"
	"log"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(sw, r)
		reqTime := time.Now().Sub(start)
		log.Printf("%s %s %d %f mcs", r.Method, r.URL.Path, sw.status, float64(reqTime.Nanoseconds())/1000)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			switch err.(type) {
			case error:
				w.WriteHeader(http.StatusInternalServerError)
				support.WarningErr(err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if auth[0] == "Basic" {
			data := strings.Split(auth[1], ":")
			if models.Authenticate(data[0], data[1]) {
				var user models.User
				db.Conn.Where("username = ?", data[0]).First(&user)
				ctx := context.WithValue(r.Context(), "user", user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}
