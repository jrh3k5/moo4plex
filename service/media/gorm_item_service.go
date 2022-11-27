package media

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/model"

	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

type GORMItemService struct {
	db *gorm.DB
}

func NewGORMItemService(db *gorm.DB) *GORMItemService {
	return &GORMItemService{
		db: db,
	}
}

func (g *GORMItemService) GetItems(ctx context.Context, mediaLibraryID int64) ([]*model.MediaItem, error) {
	var metadataItems []*gormmodel.MetadataItem
	if dbErr := g.db.Select("SELECT id, title FROM metadata_items WHERE library_subsection_id = ?", mediaLibraryID).Scan(&metadataItems).Error; dbErr != nil {
		return nil, fmt.Errorf("unable to query for metadata items for library section ID %d: %w", mediaLibraryID, dbErr)
	}

	items := make([]*model.MediaItem, len(metadataItems))
	for metadataItemIndex, metadataItem := range metadataItems {
		items[metadataItemIndex] = model.NewMediaItem(metadataItem.ID, metadataItem.Title)
	}
	return items, nil
}
