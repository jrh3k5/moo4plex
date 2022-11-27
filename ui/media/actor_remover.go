package media

import (
	"fyne.io/fyne/v2"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type ActorRemover struct {
	serviceContainer *services.ServiceContainer
	editorContainer  fyne.CanvasObject
}

func NewActorRemover(serviceContainer *services.ServiceContainer) *ActorRemover {
	actorRemover := &ActorRemover{
		serviceContainer: serviceContainer,
	}

	// TODO: implement this
	actorRemover.editorContainer = fyne.NewContainer()

	return actorRemover
}

func (a *ActorRemover) GetObject() fyne.CanvasObject {
	return a.editorContainer
}
