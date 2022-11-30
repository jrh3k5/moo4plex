package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// ActorService defines a means of interacting with actor information
type ActorService interface {
	// GetActorsForItem gets the actors for a given item
	GetActorsForItem(ctx context.Context, mediaItemID int64) ([]*model.Actor, error)

	// GetMediaItemsForActor gets all media for the given actor of the given media type
	GetMediaItemsForActor(ctx context.Context, actorID int64, mediaType model.MediaType) ([]*model.MediaItem, error)

	// RemoveActorFromItem disassociates the given actor from the given media item
	RemoveActorFromItem(ctx context.Context, mediaItemID int64, actorID int64) error
}
