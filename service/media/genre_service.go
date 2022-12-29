package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// GenreService defines means of interacting with genres
type GenreService interface {
	// GetGenres gets genres listed within the given media library
	GetGenres(ctx context.Context, mediaLibraryID int64) ([]*model.Genre, error)

	// MergeGenres merges the given genres into the target genre
	MergeGenres(ctx context.Context, mergeTarget *model.Genre, toMerge []*model.Genre) error

	// SplitGenres splits the given genre into the given target genres
	SplitGenres(ctx context.Context, splitSource *model.Genre, splitTargets []*model.Genre) error
}
