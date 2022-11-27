package media

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/model"
	gormmodel "github.com/jrh3k5/moo4plex/model/gorm"
	"gorm.io/gorm"
)

type GORMLibraryService struct {
	db *gorm.DB
}

func NewGORMLibraryService(db *gorm.DB) *GORMLibraryService {
	return &GORMLibraryService{
		db: db,
	}
}

func (g *GORMLibraryService) GetAvailableMediaTypes(ctx context.Context, mediaLibraryID int64) ([]model.MediaType, error) {
	var metadataTypes []int
	if dbErr := g.db.WithContext(ctx).Raw("SELECT DISTINCT metadata_type FROM metadata_items WHERE library_section_id = ?", mediaLibraryID).Find(&metadataTypes).Error; dbErr != nil {
		return nil, fmt.Errorf("unable to get metadata types for library section ID %d: %w", mediaLibraryID, dbErr)
	}

	var mediaTypes []model.MediaType
	for _, metadataType := range metadataTypes {
		switch metadataType {
		case int(gormmodel.Movie):
			mediaTypes = append(mediaTypes, model.Movie)
		case int(gormmodel.TelevisionSeries):
			mediaTypes = append(mediaTypes, model.TelevisionSeries)
		}
	}
	return mediaTypes, nil
}

func (g *GORMLibraryService) GetMediaLibraries(ctx context.Context) ([]*model.MediaLibrary, error) {
	var librarySections []*gormmodel.LibrarySection
	if dbErr := g.db.WithContext(ctx).Find(&librarySections).Error; dbErr != nil {
		return nil, fmt.Errorf("failed to look up all libraries: %w", dbErr)
	}
	mediaLibraries := make([]*model.MediaLibrary, len(librarySections))
	for sectionIndex, librarySection := range librarySections {
		mediaLibraries[sectionIndex] = model.NewMediaLibrary(librarySection.ID, librarySection.Name)
	}
	return mediaLibraries, nil
}
