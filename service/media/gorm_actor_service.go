package media

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/model"
	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
)

type GORMActorService struct {
	gormTagService *GORMTagService
}

func NewGORMActorService(gormTagService *GORMTagService) *GORMActorService {
	return &GORMActorService{
		gormTagService: gormTagService,
	}
}

func (g *GORMActorService) GetActorsForItem(ctx context.Context, mediaItemID int64) ([]*model.Actor, error) {
	tags, err := g.gormTagService.GetTagsForMetadataItem(ctx, gormmodel.Actor, mediaItemID)
	if err != nil {
		return nil, fmt.Errorf("unable to get actor tags for media item ID %d: %w", mediaItemID, err)
	}

	actors := make([]*model.Actor, len(tags))
	for tagIndex, tag := range tags {
		actors[tagIndex] = model.NewActor(tag.ID, tag.Tag, tag.UserThumbURL)
	}
	return actors, nil
}

func (g *GORMActorService) GetMediaItemsForActor(ctx context.Context, actorID int64, mediaType model.MediaType) ([]*model.MediaItem, error) {
	metadataItems, err := g.gormTagService.GetMetadataItemsForTags(ctx, []int64{actorID})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve metadata items for actor ID %d: %w", actorID, err)
	}

	mediaItems := make([]*model.MediaItem, len(metadataItems))
	for metadataItemIndex, metadataItem := range metadataItems {
		mediaItems[metadataItemIndex] = model.NewMediaItem(metadataItem.ID, metadataItem.Title)
	}
	return mediaItems, nil
}

func (g *GORMActorService) RemoveActorFromItem(ctx context.Context, mediaItemID int64, actorID int64) error {
	return g.gormTagService.RemoveTagsFromItem(ctx, mediaItemID, gormmodel.Actor, []int64{actorID})
}
