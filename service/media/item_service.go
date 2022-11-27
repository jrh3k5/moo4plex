package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// ItemService describes a service used to interact with media items
type ItemService interface {
	// GetItems gets all media items for the given media library ID
	GetItems(ctx context.Context, mediaLibraryID int64, mediaType model.MediaType) ([]*model.MediaItem, error)
}
