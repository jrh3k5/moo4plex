package services

import "github.com/jrh3k5/moo4plex/service/media"

// ServiceContainer is a convenient means of sharing services between UI components that cannot be initialized at startup
type ServiceContainer struct {
	genreService      media.GenreService
	hasGenreService   bool
	hasLibraryService bool
	libraryService    media.LibraryService
}

// NewServiceContainer builds a new instance of ServiceContainer
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{}
}

// GetGenreService gets the genre service, if set.
// This will panic if there is no genre service set yet.
func (sc *ServiceContainer) GetGenreService() media.GenreService {
	if !sc.hasGenreService {
		panic("No genre service set")
	}
	return sc.genreService
}

// GetLibraryService gets the library service, if set.
// This will panic if there is no library service set yet.
func (sc *ServiceContainer) GetLibraryService() media.LibraryService {
	if !sc.hasLibraryService {
		panic("No library service set")
	}
	return sc.libraryService
}

// SetGenreService sets the media.GenreService to be used by UI components
func (sc *ServiceContainer) SetGenreService(genreService media.GenreService) {
	sc.genreService = genreService
	sc.hasGenreService = true
}

// SetLibraryService sets the media.LibraryService to be used by UI components
func (sc *ServiceContainer) SetLibraryService(libraryService media.LibraryService) {
	sc.libraryService = libraryService
	sc.hasLibraryService = true
}
