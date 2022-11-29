package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ActorListMediaItem is a list of actors for a particular media item
type ActorListMediaItem struct {
	serviceContainer *services.ServiceContainer
	container        fyne.CanvasObject
	actorsList       *component.ClickableList[*model.Actor]
	actorDetails     *ActorDetails
}

func NewActorListMediaItem(serviceContainer *services.ServiceContainer, parentWindow *fyne.Window) *ActorListMediaItem {
	item := &ActorListMediaItem{
		serviceContainer: serviceContainer,
	}

	actorDetails := NewActorDetails()

	actorList := component.NewClickableList(func(a *model.Actor) string {
		return a.Name
	}, func(a *model.Actor) {
		if setErr := actorDetails.SetActor(a); setErr != nil {
			dialog.ShowError(fmt.Errorf("failed to set details for actor '%s'", a.Name), *parentWindow)
		}
	})

	item.container = container.NewGridWithColumns(2, actorList.GetObject(), actorDetails.GetObject())
	item.actorsList = actorList
	item.actorDetails = actorDetails

	return item
}

func (a *ActorListMediaItem) GetObject() fyne.CanvasObject {
	return a.container
}

// SetMediaItem sets the media item to be in contxt for this list
func (a *ActorListMediaItem) SetMediaItem(ctx context.Context, mediaItemID int64) error {
	actors, err := a.serviceContainer.GetActorService().GetActorsForItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("unable to get actors for media item ID %d: %w", mediaItemID, err)
	}
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Name < actors[j].Name
	})
	a.actorsList.SetData(actors)
	a.actorDetails.ClearActor()
	return nil
}
