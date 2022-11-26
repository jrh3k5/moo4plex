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

type GenreList struct {
	serviceContainer  *services.ServiceContainer
	selectorContainer *fyne.Container
	genresList        *widget.List
	genreFilter       *widget.Entry
	allGenres         []*model.Genre
	currentGenres     []*model.Genre
}

func NewGenreList(serviceContainer *services.ServiceContainer, width int, height int, onSelected func(*model.Genre)) *GenreList {
	component := &GenreList{
		serviceContainer: serviceContainer,
	}

	genreList := widget.NewList(func() int {
		return len(component.currentGenres)
	}, func() fyne.CanvasObject {
		button := widget.NewButton("", func() {})
		button.Alignment = widget.ButtonAlignLeading
		return button
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		button := o.(*widget.Button)
		genre := component.currentGenres[i]
		button.SetText(genre.Name)
		button.OnTapped = func() {
			onSelected(genre)
		}
	})
	genreList.Resize(fyne.NewSize(float32(width), float32(height)))

	genreFilter := widget.NewEntry()
	genreFilter.Disable()
	genreFilter.SetPlaceHolder("Filter genres")
	genreFilter.OnChanged = func(v string) {
		component.applyFilter(v)
		component.genresList.Refresh()
	}

	listContainer := fyne.NewContainer(genreList)
	listContainer.Resize(fyne.NewSize(float32(width), float32(height)))

	component.selectorContainer = container.NewVBox(genreFilter, listContainer)
	component.genreFilter = genreFilter
	component.genresList = genreList

	return component
}

// ClearGenres removes all genres from being selectable
func (g *GenreList) ClearGenres() {
	g.allGenres = nil
	g.currentGenres = nil
	g.genresList.Refresh()
}

func (g *GenreList) GetObject() fyne.CanvasObject {
	return g.selectorContainer
}

// SetGenres sets the genres to be shown
func (g *GenreList) SetGenres(ctx context.Context, mediaLibraryID int64) error {
	genres, err := g.serviceContainer.GetGenreService().GetGenres(ctx, mediaLibraryID)
	if err != nil {
		return fmt.Errorf("unable to load genres into selector for media library ID %d: %w", mediaLibraryID, err)
	}
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})
	g.allGenres = genres
	g.applyFilter(g.genreFilter.Text)
	g.genresList.Refresh()

	g.genreFilter.Enable()
	return nil
}

func (g *GenreList) applyFilter(textFilter string) {
	var currentGenres []*model.Genre
	for _, genre := range g.allGenres {
		if strings.Contains(strings.ToLower(genre.Name), strings.ToLower(textFilter)) {
			currentGenres = append(currentGenres, genre)
		}
	}
	g.currentGenres = currentGenres
}
