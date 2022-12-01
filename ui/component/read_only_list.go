package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ReadOnlyList is a read-only list of data
type ReadOnlyList[V any] struct {
	allData       []V
	listContainer fyne.CanvasObject
	dataList      *widget.List
}

// NewReadOnlyList creates a new instance of ReadOnlyList
func NewReadOnlyList[V any](nameGetter func(V) string) *ReadOnlyList[V] {
	readOnlyList := &ReadOnlyList[V]{}

	dataList := widget.NewList(func() int {
		numData := len(readOnlyList.allData)
		if numData == 0 {
			return 10
		}
		return numData
	}, func() fyne.CanvasObject {
		label := widget.NewLabel("")
		label.Alignment = fyne.TextAlignLeading
		return label
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		// The list is empty and this just a templated label to help initially fill out the list
		if i >= len(readOnlyList.allData) {
			label.SetText("")
			return
		}
		datum := readOnlyList.allData[i]
		label.SetText(nameGetter(datum))
	})

	readOnlyList.listContainer = container.NewMax(dataList)
	readOnlyList.dataList = dataList

	return readOnlyList
}

// ClearData clears the data shown within this list
func (r *ReadOnlyList[V]) ClearData() {
	r.SetData(nil)
}

func (r *ReadOnlyList[V]) GetObject() fyne.CanvasObject {
	return r.listContainer
}

// SetData sets the data to be shown within this list
func (r *ReadOnlyList[V]) SetData(data []V) {
	r.allData = data
	r.dataList.Refresh()
}
