package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// LibraryService defines a service for interacting with library data
type LibraryService interface {
	// GetAvailableMediaTypes gets the media types that are supported by the given library
	GetAvailableMediaTypes(ctx context.Context, mediaLibraryID int64) ([]model.MediaType, error)

	// GetMediaLibraries gets all known media libraries
	GetMediaLibraries(ctx context.Context) ([]*model.MediaLibrary, error)
}
