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
