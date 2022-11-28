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
