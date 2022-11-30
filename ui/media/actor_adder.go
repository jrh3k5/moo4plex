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

// ActorAdder is a component to be used to add an actor to a media item
type ActorAdder struct {
	serviceContainer *services.ServiceContainer
	container        fyne.CanvasObject
	actorList        *component.ClickableList[*model.Actor]
	actorDetails     *ActorDetails
	mediaItemActors  []*model.Actor
	currentActor     *model.Actor
}

// NewActorAdder creates a new instance of ActorAdder
func NewActorAdder(ctx context.Context, serviceContainer *services.ServiceContainer) *ActorAdder {
	actorAdder := &ActorAdder{
		serviceContainer: serviceContainer,
	}

	addActorButton := widget.NewButton("Add Actor", func() {
		fmt.Printf("adding actor '%s'\n", actorAdder.currentActor.Name)
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
	actorAdder.actorList = actorList
	actorAdder.actorDetails = actorDetails

	return actorAdder
}

func (a *ActorAdder) GetObject() fyne.CanvasObject {
	return a.container
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
	a.mediaItemActors = mediaItemActors
	return nil
}
