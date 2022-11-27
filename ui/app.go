package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	var genreSelector *mediaui.GenreSelector

	genreMerger := mediaui.NewGenreMergeEditor(ctx, &window, serviceContainer, func() {
		if refreshErr := genreSelector.RefreshGenres(ctx); refreshErr != nil {
			dialog.ShowError(fmt.Errorf("failed to refresh genre selector after save: %w", refreshErr), window)
		}
	})

	genreSelector = mediaui.NewGenreSelector(serviceContainer, func(genre *model.Genre) {
		genreMerger.SetGenre(ctx, genre)
	})

	librarySelector := mediaui.NewLibrarySelector(serviceContainer, func(m *model.MediaLibrary) {
		if m == nil {
			genreSelector.ClearGenres()
			return
		}

		if setErr := genreSelector.SetGenres(ctx, m.ID); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set genres in genre seletor after media library selection: %w", setErr), window)
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
