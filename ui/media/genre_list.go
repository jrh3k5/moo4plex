package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreList struct {
	serviceContainer *services.ServiceContainer
	clickableList    *component.ClickableList[*model.Genre]
}

func NewGenreList(serviceContainer *services.ServiceContainer, onSelected func(*model.Genre)) *GenreList {
	genreList := &GenreList{
		serviceContainer: serviceContainer,
	}

	clickableList := component.NewClickableList[*model.Genre](func(g *model.Genre) string {
		return g.Name
	}, func(g *model.Genre) {
		onSelected(g)
	})
	genreList.clickableList = clickableList
	return genreList
}

// ClearGenres removes all genres from being selectable
func (g *GenreList) ClearGenres() {
	g.clickableList.ClearData()
}

func (g *GenreList) GetObject() fyne.CanvasObject {
	return g.clickableList.GetObject()
}

// SetGenres sets the genres to be shown
func (g *GenreList) SetGenres(ctx context.Context, mediaLibraryID int64) error {
	genres, err := g.serviceContainer.GetGenreService().GetGenres(ctx, mediaLibraryID)
	if err != nil {
		return fmt.Errorf("unable to load genres into selector for media library ID %d: %w", mediaLibraryID, err)
	}
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})
	g.clickableList.SetData(genres)
	return nil
}
