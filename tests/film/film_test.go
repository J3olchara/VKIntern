package film

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

var actor1, actor2 *models.Actor
var film1, film2 *models.Film
var admin, user models.User

var host = fmt.Sprintf("http://web:%s", os.Getenv("PORT"))
var filmHost = fmt.Sprintf("%s/film", host)

func setUp() {
	var date1, date2, date3 time.Time
	db.NewConnection()
	date1 = time.Date(1990, 10, 3, 0, 0, 0, 0, time.UTC)
	date2 = time.Date(1998, 7, 11, 0, 0, 0, 0, time.UTC)
	date3 = time.Date(2005, 1, 3, 0, 0, 0, 0, time.UTC)

	admin = models.User{Username: "senya", Password: "password", Staff: true}
	user = models.User{Username: "hehehe", Password: "password", Staff: false}

	actor1 = models.NewActor("John doe", true, date1, []models.Film{})
	actor2 = models.NewActor("Arseniy Borodulin", true, date3, []models.Film{})

	film1 = models.NewFilm("John Wick", "crime film",
		date1, 10, []models.Actor{*actor1})
	film2 = models.NewFilm("VK Intern", "too crime film",
		date2, 10, []models.Actor{*actor2})

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

func TestCreateGood(t *testing.T) {
	tests := []struct {
		film models.Film
	}{
		{film: *film1},
		{film: *models.NewFilm("test", "test", time.Now(), 1, []models.Actor{})},
		{film: *models.NewFilm("test2", strings.Repeat("a", 1000), time.Now(), 9, []models.Actor{})},
	}
	for _, tc := range tests {
		t.Run("TestFilmPostGood", func(t *testing.T) {
			data, err := json.Marshal(tc.film)
			support.FatalErr(err)
			req, err := http.NewRequest(http.MethodPost, filmHost, bytes.NewBuffer(data))
			support.FatalErr(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", admin.Username, admin.Password))
			res, err := http.DefaultClient.Do(req)
			support.FatalErr(err)
			var got models.Film
			err = json.NewDecoder(res.Body).Decode(&got)
			support.FatalErr(err)
			require.Equal(t, res.StatusCode, http.StatusCreated, "/film POST is not working")
			assert.NotEqual(t, tc.film.ID, 0, "/film POST is not working correctly")
			assert.Equal(t, tc.film.Rating, got.Rating, "/film POST is not working correctly")
			assert.Equal(t, tc.film.Description, got.Description, "/film POST is not working correctly")
			assert.Equal(t, tc.film.Name, got.Name, "/film POST is not working correctly")
			got.Delete()
		})
	}
}

func TestCreateBad(t *testing.T) {
	tests := []struct {
		film *models.Film
	}{
		{film: models.NewFilm("badtest", "test", time.Now(), 14, []models.Actor{})},
		{film: models.NewFilm("badtest2", strings.Repeat("a", 1001), time.Now(), 9, []models.Actor{})},
	}
	for _, tc := range tests {
		t.Run("TestFilmPostBad", func(t *testing.T) {
			data, err := json.Marshal(tc.film)
			support.FatalErr(err)
			req, err := http.NewRequest(http.MethodPost, filmHost, bytes.NewBuffer(data))
			support.FatalErr(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", admin.Username, admin.Password))
			res, err := http.DefaultClient.Do(req)
			support.FatalErr(err)
			require.Equal(t, http.StatusBadRequest, res.StatusCode, "/film POST creates incorrect data")

		})
	}
}

func TestGet(t *testing.T) {
	film1.Create()
	film2.Create()
	expected := []models.Film{*film1, *film2}
	req, err := http.NewRequest(http.MethodGet, filmHost+"?field=id&ordering=asc", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", admin.Username, admin.Password))
	support.FatalErr(err)
	res, err := http.DefaultClient.Do(req)
	support.FatalErr(err)

	var films []models.Film
	err = json.NewDecoder(res.Body).Decode(&films)
	support.FatalErr(err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 2, len(films))
	assert.Equal(t, films[0].ID, expected[0].ID)
	assert.Equal(t, films[1].ID, expected[1].ID)
	film1.Delete()
	film2.Delete()
}
