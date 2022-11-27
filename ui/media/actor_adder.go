package media

import (
	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type ActorAdder struct {
	serviceContainer *services.ServiceContainer
	editorConainer   fyne.CanvasObject
}

func NewActorAdder(serviceContainer *services.ServiceContainer) *ActorAdder {
	actorAdder := &ActorAdder{
		serviceContainer: serviceContainer,
	}

	// TODO: actually implement this
	actorAdder.editorConainer = fyne.NewContainer()

	return actorAdder
}

func (a *ActorAdder) GetObject() fyne.CanvasObject {
	return a.editorConainer
}
