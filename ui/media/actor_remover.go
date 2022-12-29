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

type ActorRemover struct {
	serviceContainer   *services.ServiceContainer
	editorContainer    fyne.CanvasObject
	actorList          *component.ClickableList[*model.Actor]
	actorDetails       *ActorDetails
	removeActorButton  *widget.Button
	currentMediaItemID int64
	currentActor       *model.Actor
}

func NewActorRemover(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, onSave func()) *ActorRemover {
	actorRemover := &ActorRemover{
		serviceContainer: serviceContainer,
	}

	actorDetails := NewActorDetails(serviceContainer)
	actorList := component.NewClickableList(func(a *model.Actor) string {
		return a.Name
	}, func(a *model.Actor) {
		if setErr := actorDetails.SetActor(ctx, a); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set details for actor '%s'", a.Name), *parentWindow)
		}
		actorRemover.removeActorButton.Enable()
		actorRemover.currentActor = a
	}, true)
	removeActorButton := widget.NewButton("Remove Actor", func() {
		dialog.ShowConfirm("Confirm Actor Removal", fmt.Sprintf("You are about to remove the actor '%s'. Do you wish to continue?", actorRemover.currentActor.Name), func(confirmed bool) {
			if !confirmed {
				return
			}

			if removeErr := serviceContainer.GetActorService().RemoveActorFromItem(ctx, actorRemover.currentMediaItemID, actorRemover.currentActor.ID); removeErr != nil {
				dialog.ShowError(fmt.Errorf("failed to remove actor ID '%d' from media item ID '%d': %w", actorRemover.currentActor.ID, actorRemover.currentMediaItemID, removeErr), *parentWindow)
				return
			}

			onSave()
		}, *parentWindow)
	})
	removeActorButton.Disable()

	detailsContainer := container.NewBorder(nil, removeActorButton, nil, nil, actorDetails.GetObject())
	actorRemover.editorContainer = container.NewGridWithColumns(2, actorList.GetObject(), detailsContainer)
	actorRemover.actorDetails = actorDetails
	actorRemover.removeActorButton = removeActorButton
	actorRemover.actorList = actorList

	return actorRemover
}

func (a *ActorRemover) GetObject() fyne.CanvasObject {
	return a.editorContainer
}

// RefreshMediaItem will refresh the media item data within this control
func (a *ActorRemover) RefreshMediaItem(ctx context.Context) error {
	if a.currentMediaItemID > 0 {
		if setErr := a.SetMediaItem(ctx, a.currentMediaItemID); setErr != nil {
			return fmt.Errorf("failed to set the media to item ID %d: %w", a.currentMediaItemID, setErr)
		}
	}
	a.currentActor = nil
	a.actorDetails.ClearActor()
	a.removeActorButton.Disable()
	return nil
}

// SetMediaItem sets the media item to be in contxt for this list
func (a *ActorRemover) SetMediaItem(ctx context.Context, mediaItemID int64) error {
	actors, err := a.serviceContainer.GetActorService().GetActorsForItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("unable to get actors for media item ID %d: %w", mediaItemID, err)
	}
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Name < actors[j].Name
	})
	a.actorList.SetData(actors)
	a.actorDetails.ClearActor()
	a.removeActorButton.Disable()

	a.currentMediaItemID = mediaItemID

	return nil
}
