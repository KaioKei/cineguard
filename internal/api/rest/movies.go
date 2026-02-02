package rest

import (
	"cineguard/internal/utils"
	"crypto/sha256"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

// DATA MODELS

// postCreateMovie defines the request body to create a movie
type postCreateMovie struct {
	Title    string `json:"title"`
	Year     string `year:"year"`
	Director string `json:"director"`
	Synopsis string `json:"synopsis"`
	// IMDB rating is between 0 and 10
	ImdbRating uint8 `json:"imdbRating"`
	// User Rating is between 0 and 5
	UserRating uint8 `json:"userRating"`
}

// UTILS

func getMovieId(title string, director string, year string) string {
	movieIDClear := fmt.Sprintf("%s%s%s", title, director, year)
	h := sha256.New()
	h.Write([]byte(movieIDClear))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// REST METHODS

func ListMovies(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"movies": "OK"})
}

func CreateMovie(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Call BindJSON to bind the received JSON request body
		var body postCreateMovie
		if err := c.BindJSON(&body); err != nil {
			logrus.Errorf("Error parsing request body: %s", err.Error())
		}

		// verify database availability
		err := db.Ping()
		if err != nil {
			logrus.Errorf("could not connect to database: %v", err)
		}

		// first verifies that the movie does not exist yet in the table by computing its ID
		movieID := getMovieId(body.Title, body.Director, body.Year)
		movie, err := utils.SelectMovieByID(db, movieID)
		if err != nil {
			logrus.Errorf("Error selecting movie: %s", err.Error())

		}
		logrus.Infof("movie found: %s", movie.Title)

		// Create movie entry
		//movie := &models.Movie{
		//	BaseModel:  bun.BaseModel{},
		//	ID:         0,
		//	Title:      "",
		//	Year:       0,
		//	Synopsis:   "",
		//	ImdbRating: 0,
		//	UserRating: 0,
		//	CreatedAt:  time.Time{},
		//	Cast:       nil,
		//	Crew:       nil,
		//	Genres:     nil,
		//	Themes:     nil,
		//}
		//
		//// insert movie in table
		//
		//_, err = db.NewInsert().Model(user).Exec(context.Background())
		//// user.ID is now populated
		//
		//if err != nil {
		//	logrus.Errorf("could not create user: %v", err)
		//}
	}
}
