package media

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ItemEditActionContainer is a component that allows for the editing of a component
type ItemEditActionContainer struct {
	editorContainer  fyne.CanvasObject
	serviceContainer *services.ServiceContainer
}

func NewItemEditActionContainer(serviceContainer *services.ServiceContainer) *ItemEditActionContainer {
	actionContainer := &ItemEditActionContainer{
		serviceContainer: serviceContainer,
	}

	actorRemover := NewActorRemover(serviceContainer)
	actorAdder := NewActorAdder(serviceContainer)

	editorAppTabs := container.NewAppTabs(
		container.NewTabItem("Add Actor", actorAdder.GetObject()),
		container.NewTabItem("Remove Actor", actorRemover.GetObject()),
	)
	editorAppTabs.SetTabLocation(container.TabLocationBottom)

	actionContainer.editorContainer = editorAppTabs

	return actionContainer
}

func (i *ItemEditActionContainer) GetObject() fyne.CanvasObject {
	return i.editorContainer
}

func (i *ItemEditActionContainer) SetItem(ctx context.Context, mediaItem *model.MediaItem) error {
	// TODO: implement
	return nil
}
