package media

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
)

// ActorDetails will show details about a particular actor
type ActorDetails struct {
	actorImage          *canvas.Image
	actorImageContainer *fyne.Container
	actorNameLabel      *widget.Label
	detailsContainer    fyne.CanvasObject
}

func NewActorDetails() *ActorDetails {
	detailsView := &ActorDetails{}

	detailsView.actorImageContainer = container.NewMax()

	detailsView.actorNameLabel = widget.NewLabel("")
	detailsView.actorNameLabel.Alignment = fyne.TextAlignCenter

	detailsView.detailsContainer = container.NewBorder(nil, detailsView.actorNameLabel, nil, nil, detailsView.actorImageContainer)

	return detailsView
}

func (a *ActorDetails) GetObject() fyne.CanvasObject {
	return a.detailsContainer
}

// ClearActor clears the actor currently being shown
func (a *ActorDetails) ClearActor() {
	a.actorImageContainer.RemoveAll()
	a.actorImage = nil
	a.actorNameLabel.SetText("")
}

// SetActor sets the actor details to be shown
func (a *ActorDetails) SetActor(actor *model.Actor) error {
	if actor.ThumbnailURL != "" {
		imageResource, err := fyne.LoadResourceFromURLString(actor.ThumbnailURL)
		if err != nil {
			return fmt.Errorf("failed to load image from URL '%s': %w", actor.ThumbnailURL, err)
		}
		if a.actorImage == nil {
			a.actorImage = canvas.NewImageFromResource(imageResource)
			a.actorImage.FillMode = canvas.ImageFillContain
			a.actorImageContainer.Add(a.actorImage)
			a.actorImageContainer.Refresh()
		} else {
			a.actorImage.Resource = imageResource
			a.actorImage.Refresh()
		}
	}
	a.actorNameLabel.SetText(actor.Name)
	return nil
}
