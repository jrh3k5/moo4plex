package media

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreMergeEditor struct {
	serviceContainer *services.ServiceContainer
	mergeContainer   *fyne.Container
	genreList        *GenreList
	genreMerger      *GenreMerger
}

func NewGenreMergeEditor(serviceContainer *services.ServiceContainer, width int, height int) *GenreMergeEditor {
	genreMergeEditor := &GenreMergeEditor{
		serviceContainer: serviceContainer,
	}

	genreMerger := NewGenreMerger(serviceContainer)
	genreMergeEditor.genreMerger = genreMerger

	genreList := NewGenreList(serviceContainer, width/2, height, func(genre *model.Genre) {
		genreMerger.AddMerge(genre)
	})
	genreMergeEditor.genreList = genreList

	genreMergeEditor.mergeContainer = container.NewGridWithColumns(2, genreList.GetObject(), genreMerger.GetObject())

	return genreMergeEditor
}

func (g *GenreMergeEditor) GetObject() fyne.CanvasObject {
	return g.mergeContainer
}

func (g *GenreMergeEditor) SetGenre(ctx context.Context, genre *model.Genre) {
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
	g.genreMerger.SetMergeTarget(genre)
}
