package media

import (
	"context"

	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// GenreSelector allows for the selection of a genre once a media library has been selected
type GenreSelector struct {
	genreList             *GenreList
	currentMediaLibraryID int64
}

// NewGenreSelector creates a new instance of GenreSelector
func NewGenreSelector(serviceContainer *services.ServiceContainer, onSelect func(*model.Genre)) *GenreSelector {
	genreSelector := &GenreSelector{}

	genreSelector.genreList = NewGenreList(serviceContainer, onSelect)

	return genreSelector
}

// ClearGenres removes all genres from being selectable
func (g *GenreSelector) ClearGenres() {
	g.genreList.ClearGenres()
}

func (g *GenreSelector) GetObject() fyne.CanvasObject {
	return g.genreList.GetObject()
}

// RefreshGenres refreshes the genres shown in this control
func (g *GenreSelector) RefreshGenres(ctx context.Context) error {
	return g.genreList.SetGenres(ctx, g.currentMediaLibraryID)
}

// SetGenres sets the genres to be shown
func (g *GenreSelector) SetGenres(ctx context.Context, mediaLibraryID int64) error {
	g.currentMediaLibraryID = mediaLibraryID
	return g.genreList.SetGenres(ctx, mediaLibraryID)
}
