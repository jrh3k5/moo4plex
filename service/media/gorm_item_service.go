package media

import (
	"context"
	"fmt"
	"sort"

	"github.com/jrh3k5/moo4plex/model"

	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

type GORMItemService struct {
	db             *gorm.DB
	gormTagService *GORMTagService
}

func NewGORMItemService(db *gorm.DB, gormTagService *GORMTagService) *GORMItemService {
	return &GORMItemService{
		db:             db,
		gormTagService: gormTagService,
	}
}

func (g *GORMItemService) GetItems(ctx context.Context, mediaLibraryID int64, mediaType model.MediaType) ([]*model.MediaItem, error) {
	metadataType, err := g.toMetadataType(mediaType)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve metadata type for media type '%v': %w", mediaType, err)
	}

	var metadataItems []*gormmodel.MetadataItem
	if dbErr := g.db.Select("id, title, library_section_id").
		Where("library_section_id = ?", mediaLibraryID).
		Where("metadata_type = ?", int(metadataType)).
		Find(&metadataItems).
		Error; dbErr != nil {
		return nil, fmt.Errorf("unable to query for metadata items for library section ID %d: %w", mediaLibraryID, dbErr)
	}

	items := make([]*model.MediaItem, len(metadataItems))
	for metadataItemIndex, metadataItem := range metadataItems {
		items[metadataItemIndex] = model.NewMediaItem(metadataItem.ID, metadataItem.Title, metadataItem.LibrarySectionID)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items, nil
}

func (g *GORMItemService) GetItemsByAttributeSubstring(ctx context.Context, mediaLibraryID int64, textSubstring string) ([]*model.MediaItem, error) {
	metadataItems, err := g.gormTagService.GetMetadataItemsForTagSubstring(ctx, mediaLibraryID, []gormmodel.TagType{gormmodel.Actor, gormmodel.Genre}, textSubstring)
	if err != nil {
		return nil, fmt.Errorf("failed to get items in media library ID %d with text substring '%s': %w", mediaLibraryID, textSubstring, err)
	}
	mediaItems := make([]*model.MediaItem, len(metadataItems))
	for metadataItemIndex, metadataItem := range metadataItems {
		mediaItems[metadataItemIndex] = model.NewMediaItem(metadataItem.ID, metadataItem.Title, mediaLibraryID)
	}
	return mediaItems, nil
}

func (g *GORMItemService) toMetadataType(mediaType model.MediaType) (gormmodel.MetadataType, error) {
	switch mediaType {
	case model.Movie:
		return gormmodel.Movie, nil
	case model.TelevisionSeries:
		return gormmodel.TelevisionSeries, nil
	default:
		return 0, fmt.Errorf("unhandled media type; cannot convert to metadata type: '%v'\n", mediaType)
	}
}
