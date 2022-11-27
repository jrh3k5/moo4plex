package ui

import (
	"context"

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
	height := 800
	width := 700

	fyneApp := app.New()
	window := fyneApp.NewWindow("MOO4Plex")
	window.Resize(fyne.NewSize(float32(width), float32(height)))

	serviceContainer := services.NewServiceContainer()

	genreMerger := mediaui.NewGenreMergeEditor(&window, serviceContainer, width, 300)

	genreSelector := mediaui.NewGenreSelector(serviceContainer, width, 300, func(genre *model.Genre) {
		genreMerger.SetGenre(ctx, genre)
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
