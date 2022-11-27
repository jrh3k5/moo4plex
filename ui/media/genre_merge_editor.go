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
func NewGenreMergeEditor(ctx context.Context, parentWindow *fyne.Window, serviceContainer *services.ServiceContainer, onSaveCallback func()) *GenreMergeEditor {
	genreMergeEditor := &GenreMergeEditor{
		serviceContainer: serviceContainer,
		containerLabel:   widget.NewLabel("Genre:"),
	}

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	genreMerger := NewGenreMerger(ctx, parentWindow, serviceContainer, progressBar, func() {
		genreMergeEditor.ClearGenre()
	})
	genreMergeEditor.genreMerger = genreMerger

	genreList := NewGenreList(serviceContainer, func(genre *model.Genre) {
		genreMerger.AddMerge(genre)
	})
	genreMergeEditor.genreList = genreList

	selectorContainer := container.NewGridWithColumns(2, genreList.GetObject(), genreMerger.GetObject())

	genreMergeEditor.mergeContainer = container.NewBorder(genreMergeEditor.containerLabel, nil, nil, progressBar, selectorContainer)

	return genreMergeEditor
}

// ClearGenre clears the genre information in this control
func (g *GenreMergeEditor) ClearGenre() {
	g.genreList.ClearGenres()
	g.genreMerger.ClearMergeTarget()
	g.genreMerger.ClearMerges()
	g.containerLabel.SetText("Genre:")
}

func (g *GenreMergeEditor) GetObject() fyne.CanvasObject {
	return g.mergeContainer
}

// SetGenre sets the genre for which changes are to be made
func (g *GenreMergeEditor) SetGenre(ctx context.Context, genre *model.Genre) {
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
	g.genreMerger.SetMergeTarget(genre)
	g.genreMerger.ClearMerges()
	g.containerLabel.SetText(fmt.Sprintf("Genre: %s", genre.Name))
}
