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

// GenreEditor is a combinatory
type GenreEditor struct {
	serviceContainer *services.ServiceContainer
	editorContainer  *fyne.Container
	containerLabel   *widget.Label
	genreList        *GenreList
	genreMerger      *GenreMerger
}

// NewGenreEditor creates a new instance of GenreMergeEditor
func NewGenreEditor(ctx context.Context, parentWindow *fyne.Window, serviceContainer *services.ServiceContainer, onSaveCallback func()) *GenreEditor {
	genreEditor := &GenreEditor{
		serviceContainer: serviceContainer,
		containerLabel:   widget.NewLabel("Genre:"),
	}

	genreMerger := NewGenreMerger(ctx, parentWindow, serviceContainer, func() {
		genreEditor.ClearGenre()
		onSaveCallback()
	})
	genreEditor.genreMerger = genreMerger

	genreList := NewGenreList(serviceContainer, func(genre *model.Genre) {
		genreMerger.AddMerge(genre)
	})
	genreEditor.genreList = genreList

	genreMergeContainer := container.NewGridWithColumns(2, genreList.GetObject(), genreMerger.GetObject())

	genreEditorTabs := container.NewAppTabs(
		container.NewTabItem("Merge Genres", genreMergeContainer),
	)
	genreEditorTabs.SetTabLocation(container.TabLocationBottom)

	genreEditor.editorContainer = container.NewBorder(genreEditor.containerLabel, nil, nil, nil, genreEditorTabs)

	return genreEditor
}

// ClearGenre clears the genre information in this control
func (g *GenreEditor) ClearGenre() {
	g.genreList.ClearGenres()
	g.genreMerger.ClearMergeTarget()
	g.genreMerger.ClearMerges()
	g.containerLabel.SetText("Genre:")
}

func (g *GenreEditor) GetObject() fyne.CanvasObject {
	return g.editorContainer
}

// SetGenre sets the genre for which changes are to be made
func (g *GenreEditor) SetGenre(ctx context.Context, genre *model.Genre) {
	g.genreList.SetGenres(ctx, genre.MediaLibraryID)
	g.genreMerger.SetMergeTarget(genre)
	g.genreMerger.ClearMerges()
	g.containerLabel.SetText(fmt.Sprintf("Genre: %s", genre.Name))
}
