package db

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/service/media"
	"github.com/jrh3k5/moo4plex/shutdown"
	mediaui "github.com/jrh3k5/moo4plex/ui/media"
	"github.com/jrh3k5/moo4plex/ui/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// FileSelector is used to select the Plex database file to be read
type FileSelector struct {
	serviceContainer  *services.ServiceContainer
	selectorContainer *fyne.Container
}

// NewFileSelector creates a new instance of FileSelector
func NewFileSelector(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, librarySelector *mediaui.LibrarySelector, hookRegistrar shutdown.HookRegistrar) *FileSelector {
	filePathTextbox := widget.NewEntry()
	filePathTextbox.Disable()

	onOpen := func(readCloser fyne.URIReadCloser, openErr error) {
		if readCloser == nil {
			// user cancelled the operation; exit out
			return
		}

		uriString := readCloser.URI().Path()
		filePathTextbox.SetText(uriString)

		db, err := gorm.Open(sqlite.Open(uriString), &gorm.Config{})
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to open database at '%s': %w", uriString, err), *parentWindow)
			return
		}
		hookRegistrar.AddHook(func(_ context.Context) error {
			if closableDB, dbErr := db.DB(); dbErr != nil {
				return fmt.Errorf("failed to get closable database reference during closure: %w", dbErr)
			} else if closeErr := closableDB.Close(); closeErr != nil {
				return fmt.Errorf("failed to close database on shutdown: %w", closeErr)
			}
			return nil
		})

		gormTagService := media.NewGORMTagService(db)
		serviceContainer.SetLibraryService(media.NewGORMLibraryService(db))
		serviceContainer.SetGenreService(media.NewGORMGenreService(gormTagService))
		serviceContainer.SetItemService(media.NewGORMItemService(db, gormTagService))
		serviceContainer.SetActorService(media.NewGORMActorService(gormTagService))

		if err := librarySelector.SetLibraries(ctx); err != nil {
			dialog.ShowError(fmt.Errorf("failed to set libraries: %w", err), *parentWindow)
			return
		}
	}

	openButton := widget.NewButton("Open Plex DB File", func() {
		dialog.NewFileOpen(onOpen, *parentWindow).Show()
	})

	selectorContainer := container.NewGridWithColumns(2, filePathTextbox, openButton)
	return &FileSelector{
		serviceContainer:  serviceContainer,
		selectorContainer: selectorContainer,
	}
}

func (f *FileSelector) GetObject() fyne.CanvasObject {
	return f.selectorContainer
}
