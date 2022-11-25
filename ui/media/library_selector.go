package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/service/media"
)

type LibrarySelector struct {
	selector *widget.Select
}

func NewLibrarySelector(onSelect func(string)) *LibrarySelector {
	selector := widget.NewSelect([]string{}, onSelect)
	selector.Disable()

	return &LibrarySelector{
		selector: selector,
	}
}

func (l *LibrarySelector) GetObject() fyne.CanvasObject {
	return l.selector
}

func (l *LibrarySelector) SetLibraries(ctx context.Context, libraryService media.LibraryService) error {
	libraries, err := libraryService.GetMediaLibraries(ctx)
	if err != nil {
		return fmt.Errorf("failed to build library selector: %w", err)
	}

	libraryNames := make([]string, len(libraries))
	for libraryIndex, library := range libraries {
		libraryNames[libraryIndex] = library.Name
	}
	sort.Slice(libraries, func(i, j int) bool {
		return libraryNames[i] < libraryNames[j]
	})

	l.selector.Options = libraryNames
	l.selector.Enable()
	l.selector.Refresh()

	return nil
}
