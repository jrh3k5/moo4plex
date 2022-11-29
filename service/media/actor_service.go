package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// ActorService defines a means of interacting with actor information
type ActorService interface {
	// GetActorsForItem gets the actors for a given item
	GetActorsForItem(ctx context.Context, mediaItemID int64) ([]*model.Actor, error)
}