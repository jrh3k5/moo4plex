package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ItemEditActionContainer is a component that allows for the editing of a component
type ItemEditActionContainer struct {
	editorContainer  fyne.CanvasObject
	editorLabel      *widget.Label
	serviceContainer *services.ServiceContainer
	actorList        *ActorListMediaItem
	actorAdder       *ActorAdder
	actorRemover     *ActorRemover
}

func NewItemEditActionContainer(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window) *ItemEditActionContainer {
	actionContainer := &ItemEditActionContainer{
		serviceContainer: serviceContainer,
	}

	editorLabel := widget.NewLabel("Item:")
	actionContainer.editorLabel = editorLabel

	actorRemover := NewActorRemover(ctx, serviceContainer, parentWindow, func() {
		if refreshErr := actionContainer.refreshData(ctx); refreshErr != nil {
			dialog.ShowError(fmt.Errorf("failed to refresh data after removal: %w", refreshErr), *parentWindow)
		}
	})
	actorAdder := NewActorAdder(serviceContainer)
	actorList := NewActorListMediaItem(serviceContainer, parentWindow)

	editorAppTabs := container.NewAppTabs(
		container.NewTabItem("Actor List", actorList.GetObject()),
		container.NewTabItem("Add Actor", actorAdder.GetObject()),
		container.NewTabItem("Remove Actor", actorRemover.GetObject()),
	)
	editorAppTabs.SetTabLocation(container.TabLocationBottom)

	actionContainer.editorContainer = editorAppTabs
	actionContainer.actorAdder = actorAdder
	actionContainer.actorRemover = actorRemover
	actionContainer.actorList = actorList

	return actionContainer
}

func (i *ItemEditActionContainer) GetObject() fyne.CanvasObject {
	return i.editorContainer
}

func (i *ItemEditActionContainer) refreshData(ctx context.Context) error {
	if removeRefreshErr := i.actorRemover.RefreshMediaItem(ctx); removeRefreshErr != nil {
		return fmt.Errorf("failed to refresh data within actor remover: %w", removeRefreshErr)
	}
	return nil
}

// SetItem sets the media item to be used in the context of this component
func (i *ItemEditActionContainer) SetItem(ctx context.Context, mediaItem *model.MediaItem) error {
	i.editorLabel.SetText(fmt.Sprintf("Item: %s", mediaItem.Name))
	i.actorAdder.SetMediaItem(ctx, mediaItem.ID)
	i.actorRemover.SetMediaItem(ctx, mediaItem.ID)
	i.actorList.SetMediaItem(ctx, mediaItem.ID)
	return nil
}
