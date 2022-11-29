package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ItemEditor is a component that allows an item to be selected and subsequently edited
type ItemEditor struct {
	container    fyne.CanvasObject
	itemSelector *ItemSelector
}

// NewEditor creates a new instance of ItemEditor
func NewItemEditor(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window) *ItemEditor {
	itemEditor := &ItemEditor{}

	itemEditActionContainer := NewItemEditActionContainer(ctx, serviceContainer, parentWindow)
	itemSelector := NewItemSelector(ctx, serviceContainer, parentWindow, func(m *model.MediaItem) {
		if setErr := itemEditActionContainer.SetItem(ctx, m); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set action container to media item '%s': %w", m.Name, setErr), *parentWindow)
		}
	})
	itemEditor.container = container.NewGridWithRows(2, itemSelector.GetObject(), itemEditActionContainer.GetObject())
	itemEditor.itemSelector = itemSelector

	return itemEditor
}

func (i *ItemEditor) GetObject() fyne.CanvasObject {
	return i.container
}

// SetMediaLibrary sets the media library in context for this component
func (i *ItemEditor) SetMediaLibrary(ctx context.Context, mediaLibraryID int64) error {
	return i.itemSelector.SetMediaLibrary(ctx, mediaLibraryID)
}
