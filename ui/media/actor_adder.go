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

// ActorAdder is a component to be used to add an actor to a media item
type ActorAdder struct {
	serviceContainer   *services.ServiceContainer
	container          fyne.CanvasObject
	addActorButton     *widget.Button
	actorList          *component.ClickableList[*model.Actor]
	actorDetails       *ActorDetails
	currentMediaItemID int64
	mediaItemActors    []*model.Actor
	currentActor       *model.Actor
}

// NewActorAdder creates a new instance of ActorAdder
func NewActorAdder(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, onSave func()) *ActorAdder {
	actorAdder := &ActorAdder{
		serviceContainer: serviceContainer,
	}

	addActorButton := widget.NewButton("Add Actor", func() {
		dialog.ShowConfirm("Confirm Actor Addition", fmt.Sprintf("You are about to add the actor '%s' to this media item. Do you wish to continue?", actorAdder.currentActor.Name), func(confirmed bool) {
			if !confirmed {
				return
			}

			if addErr := serviceContainer.GetActorService().AddActorToItem(ctx, actorAdder.currentMediaItemID, actorAdder.currentActor.ID); addErr != nil {
				dialog.ShowError(fmt.Errorf("failed to add actor to media item: %w", addErr), *parentWindow)
				return
			}

			onSave()
		}, *parentWindow)
	})
	addActorButton.Disable()

	actorList := component.NewClickableList(func(a *model.Actor) string {
		return a.Name
	}, func(a *model.Actor) {
		actorAdder.actorDetails.SetActor(ctx, a)
		actorAdder.currentActor = a

		for _, itemActor := range actorAdder.mediaItemActors {
			if itemActor.ID == a.ID {
				addActorButton.Disable()
				return
			}
		}

		addActorButton.Enable()
	})
	actorDetails := NewActorDetails(serviceContainer)

	actionContainer := container.NewBorder(nil, addActorButton, nil, nil, actorDetails.GetObject())

	actorAdder.container = container.NewGridWithColumns(2, actorList.GetObject(), actionContainer)
	actorAdder.addActorButton = addActorButton
	actorAdder.actorList = actorList
	actorAdder.actorDetails = actorDetails

	return actorAdder
}

func (a *ActorAdder) GetObject() fyne.CanvasObject {
	return a.container
}

// RefreshMediaItem refreshes the media item data in this control
func (a *ActorAdder) RefreshMediaItem(ctx context.Context) error {
	if a.currentMediaItemID > 0 {
		if setErr := a.SetMediaItem(ctx, a.currentMediaItemID); setErr != nil {
			return fmt.Errorf("failed to set the media to item ID %d: %w", a.currentMediaItemID, setErr)
		}
	}
	a.currentActor = nil
	a.actorDetails.ClearActor()
	a.addActorButton.Disable()
	return nil
}

// SetMediaLibrary sets the media library in context for this component
func (a *ActorAdder) SetMediaLibrary(ctx context.Context, mediaLibraryID int64) error {
	actors, err := a.serviceContainer.GetActorService().GetActorsForMediaLibrary(ctx, mediaLibraryID)
	if err != nil {
		return fmt.Errorf("failed to load actors for media library ID %d: %w", mediaLibraryID, err)
	}
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Name < actors[j].Name
	})
	a.actorList.SetData(actors)
	return nil
}

// SetMediaItem sets the media item in context for this component
func (a *ActorAdder) SetMediaItem(ctx context.Context, mediaItemID int64) error {
	mediaItemActors, err := a.serviceContainer.GetActorService().GetActorsForItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("failed to load actors for media item ID %d: %w", mediaItemID, err)
	}
	a.currentMediaItemID = mediaItemID
	a.mediaItemActors = mediaItemActors
	return nil
}
