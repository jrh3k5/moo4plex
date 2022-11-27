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

// mediaLibrarySelectionCallback describes a callback to be invoked when a media library has been selected
// This can be given nil if there is no matching selection
type mediaLibrarySelectionCallback func(*model.MediaLibrary)

type LibrarySelector struct {
	serviceContainer *services.ServiceContainer
	selector         *widget.Select
	mediaLibraries   []*model.MediaLibrary
}

// NewLibrarySelector creates a new library selector
func NewLibrarySelector(serviceContainer *services.ServiceContainer, onSelect mediaLibrarySelectionCallback) *LibrarySelector {
	librarySelector := &LibrarySelector{
		serviceContainer: serviceContainer,
	}

	selector := widget.NewSelect([]string{}, func(mediaLibraryName string) {
		for _, mediaLibrary := range librarySelector.mediaLibraries {
			if mediaLibrary.Name == mediaLibraryName {
				onSelect(mediaLibrary)
				return
			}
		}
		onSelect(nil)
	})
	selector.Disable()

	librarySelector.selector = selector
	return librarySelector
}

func (l *LibrarySelector) GetObject() fyne.CanvasObject {
	return l.selector
}

func (l *LibrarySelector) SetLibraries(ctx context.Context) error {
	libraries, err := l.serviceContainer.GetLibraryService().GetMediaLibraries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get media libraries: %w", err)
	}

	libraryNames := make([]string, len(libraries))
	for libraryIndex, library := range libraries {
		libraryNames[libraryIndex] = library.Name
	}
	sort.Slice(libraries, func(i, j int) bool {
		return libraryNames[i] < libraryNames[j]
	})

	l.selector.Options = libraryNames
	l.mediaLibraries = libraries
	l.selector.Enable()
	l.selector.Refresh()

	return nil
}
