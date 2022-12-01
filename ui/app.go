package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/shutdown"
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
	shutdownHookRegistrar := shutdown.NewSliceHookRegistrar()
	defer shutdownHookRegistrar.ExecuteHooks(ctx)

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

	itemEditor := mediaui.NewItemEditor(ctx, serviceContainer, &window)

	librarySelector := mediaui.NewLibrarySelector(serviceContainer, func(m *model.MediaLibrary) {
		if m == nil {
			genreSelector.ClearGenres()
			return
		}

		if setErr := genreSelector.SetGenres(ctx, m.ID); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set genres in genre seletor after media library selection: %w", setErr), window)
		}

		if setErr := itemEditor.SetMediaLibrary(ctx, m.ID); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set media library for item selector: %w", setErr), window)
		}
	})

	dbFileSelector := db.NewFileSelector(ctx, serviceContainer, &window, librarySelector, shutdownHookRegistrar)

	dbMediaContainer := container.NewVBox(dbFileSelector.GetObject(), librarySelector.GetObject())
	genreDataContainer := container.NewGridWithRows(2,
		genreSelector.GetObject(),
		genreMerger.GetObject(),
	)

	tabbedContainer := container.NewAppTabs(
		container.NewTabItem("Genres", genreDataContainer),
		container.NewTabItem("Items", itemEditor.GetObject()),
	)

	window.SetContent(container.NewBorder(
		dbMediaContainer,
		nil,
		nil,
		nil,
		tabbedContainer,
	))

	window.ShowAndRun()
	return nil
}
