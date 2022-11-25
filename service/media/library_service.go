package media

import (
	"context"

	"github.com/jrh3k5/moo4plex/model"
)

// LibraryService defines a service for interacting with library data
type LibraryService interface {
	// GetMediaLibraries gets all known media libraries
	GetMediaLibraries(ctx context.Context) ([]*model.MediaLibrary, error)
}

type InMemoryLibraryService struct {
}

func NewInMemoryLibraryService() *InMemoryLibraryService {
	return &InMemoryLibraryService{}
}

func (i *InMemoryLibraryService) GetMediaLibraries(ctx context.Context) ([]*model.MediaLibrary, error) {
	return []*model.MediaLibrary{
		model.NewMediaLibrary("Movies"),
		model.NewMediaLibrary("Music"),
		model.NewMediaLibrary("Photos"),
	}, nil
}
