package rest

import (
	"cineguard/internal/data/models"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

func CreateUser(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := db.Ping()
		if err != nil {
			logrus.Error(err)
		}

		// Insert single user
		user := &models.User{Name: "John Doe", Email: "john@example.com"}
		_, err = db.NewInsert().Model(user).Exec(context.Background())
		// user.ID is now populated

		if err != nil {
			logrus.Errorf("could not create user: %v", err)
		}
	}

}
