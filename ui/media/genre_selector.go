package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreSelector struct {
	serviceContainer  *services.ServiceContainer
	selectorContainer *fyne.Container
	genreList         *widget.List
	genres            []*model.Genre
}

func NewGenreSelector(width int, height int, serviceContainer *services.ServiceContainer, onSelect func(*model.Genre)) *GenreSelector {
	genreSelector := &GenreSelector{
		serviceContainer: serviceContainer,
	}

	genreList := widget.NewList(func() int {
		return len(genreSelector.genres)
	}, func() fyne.CanvasObject {
		button := widget.NewButton("", func() {})
		button.Alignment = widget.ButtonAlignLeading
		return button
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		button := o.(*widget.Button)
		genre := genreSelector.genres[i]
		button.SetText(genre.Name)
		button.OnTapped = func() {
			fmt.Printf("chose genre: %s\n", genre.Name)
		}
	})
	genreList.Resize(fyne.NewSize(float32(width), float32(height)))

	genreSelector.selectorContainer = fyne.NewContainer(genreList)
	genreSelector.genreList = genreList

	return genreSelector
}

// ClearGenres removes all genres from being selectable
func (g *GenreSelector) ClearGenres() {
	g.genres = nil
	g.genreList.Refresh()
}

func (g *GenreSelector) GetObject() fyne.CanvasObject {
	return g.selectorContainer
}

// SetGenres sets the genres to be shown
func (g *GenreSelector) SetGenres(ctx context.Context, mediaLibraryID int64) error {
	genres, err := g.serviceContainer.GetGenreService().GetGenres(ctx, mediaLibraryID)
	if err != nil {
		return fmt.Errorf("unable to load genres into selector for media library ID %d: %w", mediaLibraryID, err)
	}
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})
	g.genres = genres
	g.genreList.Refresh()
	return nil
}
