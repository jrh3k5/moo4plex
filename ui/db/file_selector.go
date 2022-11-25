package db

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/service/media"
	mediaui "github.com/jrh3k5/moo4plex/ui/media"
)

type FileSelector struct {
	selectorContainer *fyne.Container
}

func NewFileSelector(ctx context.Context, parentWindow *fyne.Window, librarySelector mediaui.LibrarySelector) *FileSelector {
	filePathTextbox := widget.NewEntry()
	filePathTextbox.Disable()

	onOpen := func(readCloser fyne.URIReadCloser, openErr error) {
		if readCloser == nil {
			// user cancelled the operation; exit out
			return
		}

		uriString := readCloser.URI().String()
		filePathTextbox.SetText(uriString)

		// TODO: read file, initialize database
		libraryService := media.NewInMemoryLibraryService()
		if err := librarySelector.SetLibraries(ctx, libraryService); err != nil {
			// TODO: report that the libraries failed to load
		}
	}

	openButton := widget.NewButton("Open Plex DB File", func() {
		dialog.NewFileOpen(onOpen, *parentWindow).Show()
	})

	selectorContainer := container.NewGridWithColumns(2, filePathTextbox, openButton)
	return &FileSelector{
		selectorContainer: selectorContainer,
	}
}

func (f *FileSelector) GetObject() fyne.CanvasObject {
	return f.selectorContainer
}
