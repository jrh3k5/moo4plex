package services

import "github.com/jrh3k5/moo4plex/service/media"

// ServiceContainer is a convenient means of sharing services between UI components that cannot be initialized at startup
type ServiceContainer struct {
	actorService      media.ActorService
	genreService      media.GenreService
	hasActorService   bool
	hasGenreService   bool
	hasItemService    bool
	hasLibraryService bool
	libraryService    media.LibraryService
	itemService       media.ItemService
}

// NewServiceContainer builds a new instance of ServiceContainer
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{}
}

// GetGenreService gets the actor service, if set.
// This will panic if there is no actor service set yet.
func (sc *ServiceContainer) GetActorService() media.ActorService {
	if !sc.hasActorService {
		panic("No actor service set")
	}
	return sc.actorService
}

// GetGenreService gets the genre service, if set.
// This will panic if there is no genre service set yet.
func (sc *ServiceContainer) GetGenreService() media.GenreService {
	if !sc.hasGenreService {
		panic("No genre service set")
	}
	return sc.genreService
}

// GetItemService gets the item service, if set.
// This will panic if there is no item service set yet.
func (sc *ServiceContainer) GetItemService() media.ItemService {
	if !sc.hasItemService {
		panic("No item service set")
	}
	return sc.itemService
}

// GetLibraryService gets the library service, if set.
// This will panic if there is no library service set yet.
func (sc *ServiceContainer) GetLibraryService() media.LibraryService {
	if !sc.hasLibraryService {
		panic("No library service set")
	}
	return sc.libraryService
}

// SetActorService sets the media.ActorService to be used by the UI components
func (sc *ServiceContainer) SetActorService(actorService media.ActorService) {
	sc.actorService = actorService
	sc.hasActorService = true
}

// SetGenreService sets the media.GenreService to be used by UI components
func (sc *ServiceContainer) SetGenreService(genreService media.GenreService) {
	sc.genreService = genreService
	sc.hasGenreService = true
}

// SetItemService sets the media.ItemService to be used by UI components
func (sc *ServiceContainer) SetItemService(itemService media.ItemService) {
	sc.itemService = itemService
	sc.hasItemService = true
}

// SetLibraryService sets the media.LibraryService to be used by UI components
func (sc *ServiceContainer) SetLibraryService(libraryService media.LibraryService) {
	sc.libraryService = libraryService
	sc.hasLibraryService = true
}
