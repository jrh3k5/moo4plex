package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// GenreMergeEditor is a combinatory
type GenreMergeEditor struct {
	serviceContainer *services.ServiceContainer
	mergeContainer   *fyne.Container
	containerLabel   *widget.Label
	genreList        *GenreList
	genreMerger      *GenreMerger
}

// NewGenreMergeEditor creates a new instance of GenreMergeEditor
func NewGenreMergeEditor(parentWindow *fyne.Window, serviceContainer *services.ServiceContainer, width int, height int) *GenreMergeEditor {
	genreMergeEditor := &GenreMergeEditor{
		serviceContainer: serviceContainer,
		containerLabel:   widget.NewLabel("Genre:"),
	}

	genreMerger := NewGenreMerger(parentWindow, serviceContainer)
	genreMergeEditor.genreMerger = genreMerger

	genreList := NewGenreList(serviceContainer, width/2, height, func(genre *model.Genre) {
		genreMerger.AddMerge(genre)
	})
	genreMergeEditor.genreList = genreList

	selectorContainer := container.NewGridWithColumns(2, genreList.GetObject(), genreMerger.GetObject())

	genreMergeEditor.mergeContainer = container.NewBorder(genreMergeEditor.containerLabel, nil, nil, nil, selectorContainer)

	return genreMergeEditor
}

func (g *GenreMergeEditor) GetObject() fyne.CanvasObject {
	return g.mergeContainer
}

func (g *GenreMergeEditor) SetGenre(ctx context.Context, genre *model.Genre) {
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
	g.genreMerger.SetMergeTarget(genre)
	g.containerLabel.SetText(fmt.Sprintf("Genre: %s", genre.Name))
}
