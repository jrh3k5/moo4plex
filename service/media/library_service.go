package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// LibraryService defines a service for interacting with library data
type LibraryService interface {
	// GetMediaLibraries gets all known media libraries
	GetMediaLibraries(ctx context.Context) ([]*model.MediaLibrary, error)
}
