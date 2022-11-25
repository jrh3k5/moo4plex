package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// GenreService defines means of interacting with genres
type GenreService interface {
	// GetGenres gets genres listed within the given media library
	GetGenres(ctx context.Context, mediaLibraryID int64) ([]*model.Genre, error)
}
