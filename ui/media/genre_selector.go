package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreSelector struct {
	genreList *GenreList
}

func NewGenreSelector(serviceContainer *services.ServiceContainer, width int, height int, onSelect func(*model.Genre)) *GenreSelector {
	genreSelector := &GenreSelector{}

	genreSelector.genreList = NewGenreList(serviceContainer, width, height, func(g *model.Genre) {
		fmt.Printf("selected: %v\n", g.Name)
	})

	return genreSelector
}

// ClearGenres removes all genres from being selectable
func (g *GenreSelector) ClearGenres() {
	g.genreList.ClearGenres()
}

func (g *GenreSelector) GetObject() fyne.CanvasObject {
	return g.genreList.GetObject()
}

// SetGenres sets the genres to be shown
func (g *GenreSelector) SetGenres(ctx context.Context, mediaLibraryID int64) error {
	return g.genreList.SetGenres(ctx, mediaLibraryID)
}
