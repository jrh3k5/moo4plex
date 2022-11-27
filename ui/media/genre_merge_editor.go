package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreMergeEditor struct {
	serviceContainer *services.ServiceContainer
	mergeContainer   *fyne.Container
	genreList        *GenreList
}

func NewGenreMergeEditor(serviceContainer *services.ServiceContainer, width int, height int) *GenreMergeEditor {
	genreMerger := &GenreMergeEditor{
		serviceContainer: serviceContainer,
	}

	genreList := NewGenreList(serviceContainer, width/2, height, func(g *model.Genre) {
		fmt.Printf("merge %s\n", g.Name)
	})
	genreMerger.genreList = genreList

	genreMerger.mergeContainer = container.NewGridWithColumns(2, genreList.GetObject())

	return genreMerger
}

func (g *GenreMergeEditor) GetObject() fyne.CanvasObject {
	return g.mergeContainer
}

func (g *GenreMergeEditor) SetGenre(ctx context.Context, genre *model.Genre) {
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
}
