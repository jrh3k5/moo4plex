package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ItemSelector allows for the selection of media items
type ItemSelector struct {
	serviceContainer      *services.ServiceContainer
	container             *fyne.Container
	mediaTypeSelector     *widget.Select
	mediaItemList         *component.ClickableList[*model.MediaItem]
	currentMediaLibraryID int64
}

// NewItemSelector creates a new instance of ItemSelector
func NewItemSelector(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, onSelect func(*model.MediaItem)) *ItemSelector {
	itemSelector := &ItemSelector{
		serviceContainer: serviceContainer,
	}

	filterItems := func(filterText string) {
		mediaItems, err := serviceContainer.GetItemService().GetItemsByAttributeSubstring(ctx, itemSelector.currentMediaLibraryID, filterText)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to filter media items: %w", err), *parentWindow)
			return
		}
		sort.Slice(mediaItems, func(i, j int) bool {
			return mediaItems[i].Name < mediaItems[j].Name
		})
		itemSelector.mediaItemList.SetData(mediaItems)
	}

	filterEntry := widget.NewEntry()
	filterEntry.Disable()
	filterEntry.OnSubmitted = filterItems
	filterButton := widget.NewButton("Filter by genre or actor", func() {
		filterItems(filterEntry.Text)
	})
	filterButton.Disable()
	filterContainer := container.NewBorder(nil, nil, nil, filterButton, filterEntry)

	mediaTypeSelector := widget.NewSelect([]string{}, func(mediaTypeName string) {
		mediaType := model.MediaType(mediaTypeName)
		mediaItems, err := serviceContainer.GetItemService().GetItems(ctx, itemSelector.currentMediaLibraryID, mediaType)
		if err != nil {
			dialog.ShowError(fmt.Errorf("unable to load items for media type '%s': %w", mediaTypeName, err), *parentWindow)
		}
		itemSelector.mediaItemList.SetData(mediaItems)
		filterButton.Enable()
		filterEntry.Enable()
	})

	mediaItemList := component.NewClickableList(func(m *model.MediaItem) string {
		return m.Name
	}, onSelect)
	mediaItemList.SetPlaceholder("Filter by title")

	itemSelector.container = container.NewBorder(container.NewVBox(mediaTypeSelector, filterContainer), nil, nil, nil, mediaItemList.GetObject())
	itemSelector.mediaTypeSelector = mediaTypeSelector
	itemSelector.mediaItemList = mediaItemList

	return itemSelector
}

func (i *ItemSelector) GetObject() fyne.CanvasObject {
	return i.container
}

// SetMediaLibrary sets the context of the media library whose data is to be shown
func (i *ItemSelector) SetMediaLibrary(ctx context.Context, mediaLibraryID int64) error {
	i.currentMediaLibraryID = mediaLibraryID

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
	if len(mediaTypeNames) == 1 {
		i.mediaTypeSelector.SetSelected(mediaTypeNames[0])
	}
	i.mediaTypeSelector.Refresh()
	return nil
}
