package tests

import (
	"fmt"
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/J3olchara/VKIntern/app/server/support"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var host = fmt.Sprintf("http://web:%s", os.Getenv("PORT"))
var actorHost = fmt.Sprintf("%s/actor", host)
var filmHost = fmt.Sprintf("%s/film", host)
var admin, user models.User

func setUp() {
	db.NewConnection()

	admin = models.User{Username: "senya", Password: "password", Staff: true}
	user = models.User{Username: "hehehe", Password: "password", Staff: false}
	admin.Create()
	user.Create()
}

func tearDown() {
	admin.Delete()
	user.Delete()
	err := db.Conn.Close()
	support.FatalErr(err)
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestAdminEndpoints(t *testing.T) {
	tests := []struct {
		sender models.User
		code   int
		url    string
		method string
	}{
		{sender: admin, code: http.StatusBadRequest, url: filmHost, method: http.MethodPost},
		{sender: admin, code: http.StatusOK, url: actorHost, method: http.MethodGet},
		{sender: admin, code: http.StatusBadRequest, url: filmHost + "/1", method: http.MethodPut},
		{sender: user, code: http.StatusNotFound, url: filmHost, method: http.MethodPost},
		{sender: user, code: http.StatusOK, url: actorHost, method: http.MethodGet},
		{sender: user, code: http.StatusNotFound, url: filmHost + "/1", method: http.MethodPut},
	}
	for _, tc := range tests {
		t.Run("TestFilmPostGood", func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			support.FatalErr(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", tc.sender.Username, tc.sender.Password))
			res, err := http.DefaultClient.Do(req)
			support.FatalErr(err)
			assert.Equal(t, tc.code, res.StatusCode)
		})
	}
}
