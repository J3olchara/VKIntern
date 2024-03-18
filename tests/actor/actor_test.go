package actor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/J3olchara/VKIntern/app/server/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var actor1, actor2, actor4 *models.Actor

var host = fmt.Sprintf("http://web:%s", os.Getenv("PORT"))
var actorHost = fmt.Sprintf("%s/actor", host)
var admin, user models.User

func setUp() {
	var date1, date2, date3 time.Time
	db.NewConnection()
	date1 = time.Date(1990, 10, 3, 0, 0, 0, 0, time.UTC)
	date2 = time.Date(1998, 7, 11, 0, 0, 0, 0, time.UTC)
	date3 = time.Date(2005, 1, 3, 0, 0, 0, 0, time.UTC)

	admin = models.User{Username: "senya", Password: "password", Staff: true}
	user = models.User{Username: "hehehe", Password: "password", Staff: false}
	admin.Create()
	user.Create()

	actor1 = models.NewActor("John doe", true, date1, []models.Film{})
	actor2 = models.NewActor("Arseniy Borodulin", true, date3, []models.Film{})
	actor4 = models.NewActor("Mizullina Ekaterina", false, date2, []models.Film{})
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

func TestCreateGood(t *testing.T) {
	tests := []struct {
		actor models.Actor
	}{
		{actor: *actor1},
		{actor: *models.NewActor("test", false, time.Now(), []models.Film{})},
		{actor: *models.NewActor("test2", false, time.Now(), []models.Film{})},
	}
	for _, tc := range tests {
		t.Run("TestFilmPostGood", func(t *testing.T) {
			data, err := json.Marshal(tc.actor)
			support.FatalErr(err)
			req, err := http.NewRequest(http.MethodPost, actorHost, bytes.NewBuffer(data))
			support.FatalErr(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", admin.Username, admin.Password))
			res, err := http.DefaultClient.Do(req)
			support.FatalErr(err)
			var got models.Actor
			err = json.NewDecoder(res.Body).Decode(&got)
			support.FatalErr(err)
			require.Equal(t, res.StatusCode, http.StatusCreated, "/film POST is not working")
			assert.NotEqual(t, tc.actor.ID, 0, "/film POST is not working correctly")
			assert.Equal(t, tc.actor.Gender, got.Gender, "/film POST is not working correctly")
			assert.Equal(t, tc.actor.Name, got.Name, "/film POST is not working correctly")
			got.Delete()
		})
	}
}

func TestCreateBad(t *testing.T) {
	tests := []struct {
		actor *models.Actor
	}{
		{actor: models.NewActor(strings.Repeat("a", 151), false, time.Now(), []models.Film{})},
	}
	for _, tc := range tests {
		t.Run("TestFilmPostBad", func(t *testing.T) {
			data, err := json.Marshal(tc.actor)
			support.FatalErr(err)
			req, err := http.NewRequest(http.MethodPost, actorHost, bytes.NewBuffer(data))
			support.FatalErr(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", admin.Username, admin.Password))
			res, err := http.DefaultClient.Do(req)
			support.FatalErr(err)
			require.Equal(t, res.StatusCode, http.StatusBadRequest, "/film POST creates incorrect data")
		})
	}
}
