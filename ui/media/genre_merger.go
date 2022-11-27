package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreMerger struct {
	serviceContainer *services.ServiceContainer
	mergeContainer   *fyne.Container
	genreList        *GenreList
}

func NewGenreMerger(serviceContainer *services.ServiceContainer, width int, height int) *GenreMerger {
	genreMerger := &GenreMerger{
		serviceContainer: serviceContainer,
	}

	genreList := NewGenreList(serviceContainer, width/2, height, func(g *model.Genre) {
		fmt.Printf("merge %s\n", g.Name)
	})
	genreMerger.genreList = genreList

	genreMerger.mergeContainer = container.NewGridWithColumns(2, genreList.GetObject())

	return genreMerger
}

func (g *GenreMerger) GetObject() fyne.CanvasObject {
	return g.mergeContainer
}

func (g *GenreMerger) SetGenre(ctx context.Context, genre *model.Genre) {
	fmt.Printf("loading genres for library ID %d\n", genre.MediaLibraryID)
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
}
