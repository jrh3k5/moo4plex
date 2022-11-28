package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ActorListMediaItem is a list of actors for a particular media item
type ActorListMediaItem struct {
	serviceContainer *services.ServiceContainer
	container        fyne.CanvasObject
	actorsList       *component.ClickableList[*model.Actor]
}

func NewActorListMediaItem(serviceContainer *services.ServiceContainer) *ActorListMediaItem {
	item := &ActorListMediaItem{
		serviceContainer: serviceContainer,
	}

	actorList := component.NewClickableList(func(a *model.Actor) string {
		return a.Name
	}, func(a *model.Actor) {
		fmt.Printf("selected '%s'\n", a.Name)
	})
	item.actorsList = actorList

	item.container = container.NewGridWithColumns(2, actorList.GetObject())

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
	return nil
}
