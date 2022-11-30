package media

import (
	"context"
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// ActorDetails will show details about a particular actor
type ActorDetails struct {
	serviceContainer    *services.ServiceContainer
	actorImage          *canvas.Image
	actorImageContainer *fyne.Container
	actorNameLabel      *widget.Label
	detailsContainer    fyne.CanvasObject
	movieList           *component.ReadOnlyList[*model.MediaItem]
	imageLoadCancel     context.CancelFunc
}

// NewActorDetails creates a new instance of ActorDetails
func NewActorDetails(serviceContainer *services.ServiceContainer) *ActorDetails {
	detailsView := &ActorDetails{
		serviceContainer: serviceContainer,
	}

	detailsView.actorImageContainer = container.NewMax()

	detailsView.actorNameLabel = widget.NewLabel("")
	detailsView.actorNameLabel.Alignment = fyne.TextAlignCenter

	detailsView.movieList = component.NewReadOnlyList(func(m *model.MediaItem) string {
		return m.Name
	})

	bioDetailsContainer := container.NewBorder(nil, detailsView.actorNameLabel, nil, nil, detailsView.actorImageContainer)
	detailsView.detailsContainer = container.NewGridWithRows(2, bioDetailsContainer, detailsView.movieList.GetObject())

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
func (a *ActorDetails) SetActor(ctx context.Context, actor *model.Actor) error {
	if a.imageLoadCancel != nil {
		a.imageLoadCancel()
		a.imageLoadCancel = nil
	}

	mediaItems, err := a.serviceContainer.GetActorService().GetMediaItemsForActor(ctx, actor.ID, model.Movie)
	if err != nil {
		return fmt.Errorf("failed to load movies for actor ID %d: %w", actor.ID, err)
	}
	sort.Slice(mediaItems, func(i, j int) bool {
		return mediaItems[i].Name < mediaItems[j].Name
	})
	a.movieList.SetData(mediaItems)

	// TODO: make image load cancellable?
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("ERROR: panic occurred while trying to load image for actor: %v\n", r)
			}
		}()

		cancelCtx, cancelFunc := context.WithCancel(context.Background())
		a.imageLoadCancel = cancelFunc

		if actor.ThumbnailURL != "" {
			imageResource, err := fyne.LoadResourceFromURLString(actor.ThumbnailURL)
			if err != nil {
				fmt.Printf("ERROR: failed to load image from URL '%s': %v\n", actor.ThumbnailURL, err)
				return
			}

			if ctxErr := cancelCtx.Err(); ctxErr != nil {
				// context was cancelled by another load, so stop here
				return
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
	}()
	a.actorNameLabel.SetText(actor.Name)
	return nil
}
