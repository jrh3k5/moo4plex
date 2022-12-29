package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/component"
	"github.com/jrh3k5/moo4plex/ui/services"
)

// GenreSplitter is a UI component used to 'split' a genre into other, existing genres
type GenreSplitter struct {
	serviceContainer        *services.ServiceContainer
	splitButton             *widget.Button
	splitSource             *model.Genre
	splittableTargetsList   *GenreList
	selectedSplitTargets    []*model.Genre
	selectedSplitTargetList *component.ClickableList[*model.Genre]
	container               fyne.CanvasObject
}

func NewGenreSplitter(ctx context.Context, serviceContainer *services.ServiceContainer, parentWindow *fyne.Window, onSaveCallback func()) *GenreSplitter {
	splitter := &GenreSplitter{
		serviceContainer: serviceContainer,
	}

	splittableTargetsList := NewGenreList(serviceContainer, splitter.addSplitTarget)
	splitter.splittableTargetsList = splittableTargetsList

	splitButton := widget.NewButton("Split Genres", func() {
		dialog.ShowConfirm("Split Genres", fmt.Sprintf("This will split the genre '%s' into %d genres. Do you wish to continue?", splitter.splitSource.Name, len(splitter.selectedSplitTargets)), func(confirmed bool) {
			if !confirmed {
				return
			}

			if splitErr := serviceContainer.GetGenreService().SplitGenres(ctx, splitter.splitSource, splitter.selectedSplitTargets); splitErr != nil {
				dialog.ShowError(fmt.Errorf("failed to split genres: %w", splitErr), *parentWindow)
			}
			onSaveCallback()
		}, *parentWindow)
	})
	splitButton.Disable()
	splitter.splitButton = splitButton

	splitTargetList := component.NewClickableList(func(v *model.Genre) string {
		return v.Name
	}, splitter.removeSplitTarget, false)
	splitter.selectedSplitTargetList = splitTargetList

	splitTargetsContainer := container.NewBorder(splitButton, nil, nil, nil, splitTargetList.GetObject())

	splitter.container = container.NewGridWithColumns(2, splittableTargetsList.GetObject(), splitTargetsContainer)

	return splitter
}

// ClearSplits clears the context from this splitter component
func (g *GenreSplitter) ClearSplits() {
	g.selectedSplitTargets = nil
	g.selectedSplitTargetList.ClearData()
	g.splitSource = nil

	g.splitButton.Disable()
}

func (g *GenreSplitter) GetObject() fyne.CanvasObject {
	return g.container
}

// SetSplitSource sets the genre that is to be split into other targets
func (g *GenreSplitter) SetSplitSource(splitSource *model.Genre) {
	g.splitSource = splitSource
}

// SetTargettableGenres sets the genres to which a genre can be split
func (g *GenreSplitter) SetTargettableGenres(ctx context.Context, mediaLibraryID int64) error {
	return g.splittableTargetsList.SetGenres(ctx, mediaLibraryID)
}

func (g *GenreSplitter) addSplitTarget(splitTarget *model.Genre) {
	// don't let a genre be split into itself
	if splitTarget.ID == g.splitSource.ID {
		return
	}

	for _, existingSplitTarget := range g.selectedSplitTargets {
		if existingSplitTarget.ID == splitTarget.ID {
			// don't re-add a split target
			return
		}
	}

	g.selectedSplitTargets = append(g.selectedSplitTargets, splitTarget)
	g.selectedSplitTargetList.SetData(g.selectedSplitTargets)

	g.splitButton.Enable()
}

func (g *GenreSplitter) removeSplitTarget(splitTarget *model.Genre) {
	for i := len(g.selectedSplitTargets) - 1; i >= 0; i-- {
		if g.selectedSplitTargets[i].ID == splitTarget.ID {
			g.selectedSplitTargets = append(g.selectedSplitTargets[:i], g.selectedSplitTargets[i+1:]...)
		}
	}
	g.selectedSplitTargetList.SetData(g.selectedSplitTargets)
}
