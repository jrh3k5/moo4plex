package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ItemSelector allows for the selection of media items
type ItemSelector struct {
	serviceContainer  *services.ServiceContainer
	container         *fyne.Container
	mediaTypeSelector *widget.Select
	mediaItemList     *component.ClickableList[*model.MediaItem]
}

func NewItemSelector(serviceContainer *services.ServiceContainer) *ItemSelector {
	itemSelector := &ItemSelector{
		serviceContainer: serviceContainer,
	}

	mediaTypeSelector := widget.NewSelect([]string{}, func(mediaTypeName string) {
		// TODO; load media items
		fmt.Sprintf("selected '%s'\n", mediaTypeName)
	})

	itemSelector.container = container.NewBorder(mediaTypeSelector, nil, nil, nil)
	itemSelector.mediaTypeSelector = mediaTypeSelector

	return itemSelector
}

func (i *ItemSelector) GetObject() fyne.CanvasObject {
	return i.container
}

// SetMediaLibrary sets the context of the media library whose data is to be shown
func (i *ItemSelector) SetMediaLibrary(ctx context.Context, mediaLibraryID int64) error {
	mediaTypes, err := i.serviceContainer.GetLibraryService().GetAvailableMediaTypes(ctx, mediaLibraryID)
	if err != nil {
		return fmt.Errorf("failed to load items for media library ID %d: %w", mediaLibraryID, err)
	}
	mediaTypeNames := make([]string, len(mediaTypes))
	for mediaTypeIndex, mediaType := range mediaTypes {
		mediaTypeNames[mediaTypeIndex] = string(mediaType)
	}
	sort.Slice(mediaTypeNames, func(i, j int) bool {
		return mediaTypeNames[i] < mediaTypeNames[j]
	})
	i.mediaTypeSelector.Options = mediaTypeNames
	i.mediaTypeSelector.Refresh()
	return nil
}
