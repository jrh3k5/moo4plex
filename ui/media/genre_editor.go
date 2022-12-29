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
	genreMergerList  *GenreList
	genreMerger      *GenreMerger
	genreSplitter    *GenreSplitter
}

// NewGenreEditor creates a new instance of GenreMergeEditor
func NewGenreEditor(ctx context.Context, parentWindow *fyne.Window, serviceContainer *services.ServiceContainer, onSaveCallback func()) *GenreEditor {
	genreEditor := &GenreEditor{
		serviceContainer: serviceContainer,
		containerLabel:   widget.NewLabel("Genre:"),
	}

	onSave := func() {
		genreEditor.ClearGenre()
		onSaveCallback()
	}

	genreMerger := NewGenreMerger(ctx, parentWindow, serviceContainer, onSave)
	genreEditor.genreMerger = genreMerger

	genreList := NewGenreList(serviceContainer, func(genre *model.Genre) {
		genreMerger.AddMerge(genre)
	})
	genreEditor.genreMergerList = genreList

	// TODO: refactor this so that these components are all contained within genre_merger
	genreMergeContainer := container.NewGridWithColumns(2, genreList.GetObject(), genreMerger.GetObject())

	genreSplitter := NewGenreSplitter(ctx, serviceContainer, parentWindow, onSave)
	genreEditor.genreSplitter = genreSplitter
	genreSplitContainer := genreSplitter.GetObject()

	genreEditorTabs := container.NewAppTabs(
		container.NewTabItem("Merge Genres", genreMergeContainer),
		container.NewTabItem("Split Genres", genreSplitContainer),
	)
	genreEditorTabs.SetTabLocation(container.TabLocationBottom)

	genreEditor.editorContainer = container.NewBorder(genreEditor.containerLabel, nil, nil, nil, genreEditorTabs)

	return genreEditor
}

// ClearGenre clears the genre information in this control
func (g *GenreEditor) ClearGenre() {
	// clear genre merger
	g.genreMergerList.ClearGenres()
	g.genreMerger.ClearMergeTarget()
	g.genreMerger.ClearMerges()

	// clear genre splitter
	g.genreSplitter.ClearSplits()

	// reset the name to indicate that there's nothing currently selected
	g.containerLabel.SetText("Genre:")
}

func (g *GenreEditor) GetObject() fyne.CanvasObject {
	return g.editorContainer
}

// SetGenre sets the genre for which changes are to be made
func (g *GenreEditor) SetGenre(ctx context.Context, genre *model.Genre) error {
	// set the contents of the genre merger
	if mergeGenresSetErr := g.genreMergerList.SetGenres(ctx, genre.MediaLibraryID); mergeGenresSetErr != nil {
		return fmt.Errorf("failed to set genres in genre merger: %w", mergeGenresSetErr)
	}
	g.genreMerger.SetMergeTarget(genre)
	g.genreMerger.ClearMerges()

	// set the genres for splitting
	if splitGenresSetErr := g.genreSplitter.SetTargettableGenres(ctx, genre.MediaLibraryID); splitGenresSetErr != nil {
		return fmt.Errorf("failed to set targettable genres in genre splitter: %w", splitGenresSetErr)
	}
	g.genreSplitter.SetSplitSource(genre)

	// indicate the current, in-context genre
	g.containerLabel.SetText(fmt.Sprintf("Genre: %s", genre.Name))

	return nil
}
