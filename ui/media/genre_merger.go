package media

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jrh3k5/moo4plex/model"
	"github.com/jrh3k5/moo4plex/ui/services"
)

type GenreMerger struct {
	serviceContainer *services.ServiceContainer
	mergerComponent  fyne.CanvasObject
	mergeButton      *widget.Button
	toMergeList      *widget.List
	mergeTarget      *model.Genre
	toMerge          []*model.Genre
}

func NewGenreMerger(ctx context.Context, parentWindow *fyne.Window, serviceContainer *services.ServiceContainer) *GenreMerger {
	merger := &GenreMerger{
		serviceContainer: serviceContainer,
	}

	mergeButton := widget.NewButton("Merge Genres", func() {
		dialog.ShowConfirm("Confirm Merge", fmt.Sprintf("You are about to merge %d genres into the genre '%s'. Do you wish to continue?", len(merger.toMerge), merger.mergeTarget.Name), func(confirmed bool) {
			if confirmed {
				if mergeErr := serviceContainer.GetGenreService().MergeGenres(ctx, merger.mergeTarget, merger.toMerge); mergeErr != nil {
					// TODO: report error
					fmt.Printf("failed to merge genres: %v\n", mergeErr)
				}
			}
		}, *parentWindow)
	})
	mergeButton.Disable()
	merger.mergeButton = mergeButton

	toMergeList := widget.NewList(func() int {
		toMergeCount := len(merger.toMerge)
		if toMergeCount < 10 {
			return 10
		}
		return toMergeCount
	}, func() fyne.CanvasObject {
		button := widget.NewButton("", func() {})
		button.Alignment = widget.ButtonAlignLeading
		button.Disable()
		return button
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		button := o.(*widget.Button)
		// The list is empty and this just a templated button to help initially fill out the list
		if i >= len(merger.toMerge) {
			button.SetText("")
			button.Disable()
			return
		}
		genre := merger.toMerge[i]
		button.SetText(genre.Name)
		button.OnTapped = func() {
			merger.removeMerge(genre)
		}
		button.Enable()
	})
	merger.toMergeList = toMergeList

	merger.mergerComponent = container.NewBorder(mergeButton, nil, nil, nil, container.NewMax(toMergeList))

	return merger
}

func (g *GenreMerger) GetObject() fyne.CanvasObject {
	return g.mergerComponent
}

// AddMerge adds the given genre to be merged into the configured merge target
func (g *GenreMerger) AddMerge(genre *model.Genre) {
	if genre.ID == g.mergeTarget.ID {
		return
	}

	for _, mergeable := range g.toMerge {
		if mergeable.ID == genre.ID {
			return
		}
	}

	g.toMerge = append(g.toMerge, genre)
	g.toMergeList.Refresh()

	g.mergeButton.Enable()
}

// SetMergeTarget sets the target genre into which the selected genres will be merged
func (g *GenreMerger) SetMergeTarget(genre *model.Genre) {
	g.mergeTarget = genre
}

func (g *GenreMerger) removeMerge(genre *model.Genre) {
	for i := len(g.toMerge) - 1; i >= 0; i-- {
		if g.toMerge[i].ID == genre.ID {
			g.toMerge = append(g.toMerge[0:i], g.toMerge[i+1:]...)
		}
	}
	g.toMergeList.Refresh()
	if len(g.toMerge) == 0 {
		g.mergeButton.Disable()
	}
}
