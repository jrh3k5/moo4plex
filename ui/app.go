package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/service/media"
	"github.com/jrh3k5/moo4plex/ui/db"
	mediaui "github.com/jrh3k5/moo4plex/ui/media"
)

type App struct {
	libraryService media.LibraryService
}

func NewApp(libraryService media.LibraryService) *App {
	return &App{
		libraryService: libraryService,
	}
}

func (a *App) Run(ctx context.Context) error {
	fyneApp := app.New()
	window := fyneApp.NewWindow("MOO4Plex")

	librarySelector := mediaui.NewLibrarySelector(func(v string) {
		fmt.Printf("chosen is: %v\n", v)
	})

	dbFileSelector := db.NewFileSelector(ctx, &window, *librarySelector)

	window.SetContent(container.NewVBox(
		dbFileSelector.GetObject(),
		librarySelector.GetObject(),
	))

	window.Resize(fyne.NewSize(600, 800))

	window.ShowAndRun()
	return nil
}
