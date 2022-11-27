package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/db"
	mediaui "github.com/jrh3k5/moo4plex/ui/media"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	fyneApp := app.New()
	window := fyneApp.NewWindow("MOO4Plex")
	window.Resize(fyne.NewSize(600, 800))

	serviceContainer := services.NewServiceContainer()

	genreMerger := mediaui.NewGenreMerger(serviceContainer, 600, 300)

	genreSelector := mediaui.NewGenreSelector(serviceContainer, 600, 300, func(g *model.Genre) {
		fmt.Printf("selected genre %v\n", g.Name)
	})

	librarySelector := mediaui.NewLibrarySelector(serviceContainer, func(m *model.MediaLibrary) {
		if m == nil {
			// TODO: reset genres
			return
		}

		if setErr := genreSelector.SetGenres(ctx, m.ID); setErr != nil {
			// TODO: handle error
		}
	})

	dbFileSelector := db.NewFileSelector(ctx, serviceContainer, &window, librarySelector)

	dbMediaContainer := container.NewVBox(dbFileSelector.GetObject(), librarySelector.GetObject())
	genreDataContainer := container.NewGridWithRows(2,
		genreSelector.GetObject(),
		genreMerger.GetObject(),
	)

	window.SetContent(container.NewBorder(
		dbMediaContainer,
		nil,
		nil,
		nil,
		genreDataContainer,
	))

	window.ShowAndRun()
	return nil
}
