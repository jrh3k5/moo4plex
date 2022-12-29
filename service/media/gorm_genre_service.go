package media

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/model"
	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
)

type GORMGenreService struct {
	gormTagService *GORMTagService
}

func NewGORMGenreService(gormTagService *GORMTagService) *GORMGenreService {
	return &GORMGenreService{
		gormTagService: gormTagService,
	}
}

func (g *GORMGenreService) DeleteGenre(ctx context.Context, toDelete *model.Genre) error {
	return g.gormTagService.DeleteTagAssociations(ctx, toDelete.ID, gormmodel.Genre)
}

func (g *GORMGenreService) GetGenres(ctx context.Context, mediaLibraryID int64) ([]*model.Genre, error) {
	tags, err := g.gormTagService.GetTagsForLibrarySection(ctx, gormmodel.Genre, mediaLibraryID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve genres for media library %d: %w", mediaLibraryID, err)
	}

	genres := make([]*model.Genre, len(tags))
	for tagIndex, tag := range tags {
		genres[tagIndex] = model.NewGenre(tag.ID, tag.Tag, tag.LibrarySectionID)
	}
	return genres, nil
}

func (g *GORMGenreService) MergeGenres(ctx context.Context, mergeTarget *model.Genre, toMerge []*model.Genre) error {
	toMergeIDs := make([]int64, len(toMerge))
	for mergeIndex, mergeable := range toMerge {
		toMergeIDs[mergeIndex] = mergeable.ID
	}

	return g.gormTagService.ReplaceTags(ctx, mergeTarget.MediaLibraryID, gormmodel.Genre, toMergeIDs, []int64{mergeTarget.ID})
}

func (g *GORMGenreService) SplitGenres(ctx context.Context, splitSource *model.Genre, splitTargets []*model.Genre) error {
	splitTargetIDs := make([]int64, len(splitTargets))
	for splitTargetIndex, splitTarget := range splitTargets {
		splitTargetIDs[splitTargetIndex] = splitTarget.ID
	}

	return g.gormTagService.ReplaceTags(ctx, splitSource.MediaLibraryID, gormmodel.Genre, []int64{splitSource.ID}, splitTargetIDs)
}
