package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// GenreDeleter is a component used to delete a genre
type GenreDeleter struct {
	deleteButton *widget.Button
	toDelete     *model.Genre
}

func NewGenreDeleter(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, onSave func()) *GenreDeleter {
	deleter := &GenreDeleter{}

	deleteButton := widget.NewButton("Delete Genre", func() {
		dialog.ShowConfirm("Delete Genre", fmt.Sprintf("You are about to delete the genre '%s'. Are you sure you wish to delete it?", deleter.toDelete.Name), func(confirmed bool) {
			if !confirmed {
				return
			}

			if deleteErr := serviceContainer.GetGenreService().DeleteGenre(ctx, deleter.toDelete); deleteErr != nil {
				dialog.ShowError(fmt.Errorf("failed to delete tag '%s': %w", deleter.toDelete.Name, deleteErr), *parentWindow)
				return
			}

			onSave()
		}, *parentWindow)
	})
	deleteButton.Disable()
	deleter.deleteButton = deleteButton

	return deleter
}

func (g *GenreDeleter) ClearGenre() {
	g.toDelete = nil
	g.deleteButton.Disable()
}

func (g *GenreDeleter) GetObject() fyne.CanvasObject {
	return g.deleteButton
}

// SetGenre sets the genre to be deleted
func (g *GenreDeleter) SetGenre(toDelete *model.Genre) {
	g.toDelete = toDelete
	g.deleteButton.Enable()
}
