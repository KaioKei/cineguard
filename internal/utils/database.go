package utils

import (
	"cineguard/internal/data/models"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SelectMovieByID(db *bun.DB, id string) (*models.Movie, error) {
	ctx := context.Background()

	var selectedMovie models.Movie
	if err := db.NewSelect().Model(&selectedMovie).Where("ID = ?", id).Scan(ctx); err != nil {
		return nil, fmt.Errorf("cannot select movie by ID '%s': %w", id, err)
	}

	return &selectedMovie, nil
}
