package media

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreSelector struct {
	serviceContainer  *services.ServiceContainer
	selectorContainer *fyne.Container
	genreList         *widget.List
	genreFilter       *widget.Entry
	allGenres         []*model.Genre
	currentGenres     []*model.Genre
}

func NewGenreSelector(width int, height int, serviceContainer *services.ServiceContainer, onSelect func(*model.Genre)) *GenreSelector {
	genreSelector := &GenreSelector{
		serviceContainer: serviceContainer,
	}

	genreList := widget.NewList(func() int {
		return len(genreSelector.currentGenres)
	}, func() fyne.CanvasObject {
		button := widget.NewButton("", func() {})
		button.Alignment = widget.ButtonAlignLeading
		return button
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		button := o.(*widget.Button)
		genre := genreSelector.currentGenres[i]
		button.SetText(genre.Name)
		button.OnTapped = func() {
			fmt.Printf("chose genre: %s\n", genre.Name)
		}
	})
	genreList.Resize(fyne.NewSize(float32(width), float32(height)))

	genreFilter := widget.NewEntry()
	genreFilter.Disable()
	genreFilter.SetPlaceHolder("Filter genres")
	genreFilter.OnChanged = func(v string) {
		genreSelector.applyFilter(v)
		genreSelector.genreList.Refresh()
	}

	listContainer := fyne.NewContainer(genreList)
	listContainer.Resize(fyne.NewSize(float32(width), float32(height)))

	genreSelector.selectorContainer = container.NewVBox(genreFilter, listContainer)
	genreSelector.genreFilter = genreFilter
	genreSelector.genreList = genreList

	return genreSelector
}

// ClearGenres removes all genres from being selectable
func (g *GenreSelector) ClearGenres() {
	g.allGenres = nil
	g.currentGenres = nil
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
	g.allGenres = genres
	g.applyFilter(g.genreFilter.Text)
	g.genreList.Refresh()

	g.genreFilter.Enable()
	return nil
}

func (g *GenreSelector) applyFilter(textFilter string) {
	var currentGenres []*model.Genre
	for _, genre := range g.allGenres {
		if strings.Contains(strings.ToLower(genre.Name), strings.ToLower(textFilter)) {
			currentGenres = append(currentGenres, genre)
		}
	}
	g.currentGenres = currentGenres
}
