package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type ActorAdder struct {
	serviceContainer *services.ServiceContainer
	actorList        *component.ClickableList[*model.Actor]
}

func NewActorAdder(serviceContainer *services.ServiceContainer) *ActorAdder {
	actorAdder := &ActorAdder{
		serviceContainer: serviceContainer,
	}

	actorList := component.NewClickableList(func(a *model.Actor) string {
		return a.Name
	}, func(a *model.Actor) {
		fmt.Printf("Removing actor '%s'\n", a.Name)
	})
	actorAdder.actorList = actorList

	return actorAdder
}

func (a *ActorAdder) GetObject() fyne.CanvasObject {
	return a.actorList.GetObject()
}

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

func (a *ActorAdder) SetMediaItem(ctx context.Context, mediaItemID int64) error {
	// TODO: implement
	return nil
}
