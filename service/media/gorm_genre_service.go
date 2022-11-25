package media

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/model"
	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

type GORMGenreService struct {
	db *gorm.DB
}

func NewGORMGenreService(db *gorm.DB) *GORMGenreService {
	return &GORMGenreService{
		db: db,
	}
}

func (g *GORMGenreService) GetGenres(ctx context.Context, mediaLibraryID int64) ([]*model.Genre, error) {
	var tags []*gormmodel.Tag
	queryDB := g.db.WithContext(ctx).Distinct("id, tag, tag_type").
		Find(&tags, "tag_type = 1").
		Joins("taggings on taggings.tag_id = tag.id").
		Joins("metadata_items on metadata_items.id = taggings.metadata_item_id and metadata_items.library_section_id = ?", mediaLibraryID)
	if dbErr := queryDB.Error; dbErr != nil {
		return nil, fmt.Errorf("failed to resolve genres for media library %d: %w", mediaLibraryID, dbErr)
	}

	genres := make([]*model.Genre, len(tags))
	for tagIndex, tag := range tags {
		genres[tagIndex] = model.NewGenre(tag.ID, tag.Tag)
	}
	return genres, nil
}
