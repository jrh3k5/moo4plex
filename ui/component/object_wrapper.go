package component

import (
	"fyne.io/fyne/v2"
)

type ObjectWrapper interface {
	GetObject() fyne.CanvasObject
}
