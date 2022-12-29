package component

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ClickableList is a reusable component that shows a list of clickable items
type ClickableList[V any] struct {
	listContainer *fyne.Container
	dataList      *widget.List
	listFilter    *widget.Entry
	nameGetter    func(V) string
	allData       []V
	currentData   []V
}

// NewClickableList creates a new instance of ClickableList
func NewClickableList[V any](nameGetter func(V) string, clickHandler func(V), withFilter bool) *ClickableList[V] {
	clickableList := &ClickableList[V]{
		nameGetter: nameGetter,
	}

	dataList := widget.NewList(func() int {
		numData := len(clickableList.currentData)
		if numData == 0 {
			return 10
		}
		return numData
	}, func() fyne.CanvasObject {
		button := widget.NewButton("", func() {})
		button.Alignment = widget.ButtonAlignLeading
		button.Disable()
		return button
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		button := o.(*widget.Button)
		// The list is empty and this just a templated button to help initially fill out the list
		if i >= len(clickableList.currentData) {
			button.SetText("")
			button.Disable()
			return
		}
		datum := clickableList.currentData[i]
		button.SetText(nameGetter(datum))
		button.OnTapped = func() {
			clickHandler(datum)
		}
		button.Enable()
	})

	if withFilter {
		listFilter := widget.NewEntry()
		listFilter.Disable()
		listFilter.SetPlaceHolder("Filter")
		listFilter.OnChanged = func(v string) {
			clickableList.applyFilter(v)
			if clickableList.listFilter != nil {
				clickableList.listFilter.Refresh()
			}
		}
		clickableList.listContainer = container.NewBorder(listFilter, nil, nil, nil, container.NewMax(dataList))
		clickableList.listFilter = listFilter
	} else {
		clickableList.listContainer = container.NewMax(dataList)
	}

	clickableList.dataList = dataList

	return clickableList
}

// ClearData clears all clickable data from the list
func (c *ClickableList[V]) ClearData() {
	c.allData = nil
	c.currentData = nil
	c.dataList.Refresh()
}

func (c *ClickableList[V]) GetObject() fyne.CanvasObject {
	return c.listContainer
}

// SetData sets the data to be displayed within the list
func (c *ClickableList[V]) SetData(data []V) {
	c.allData = data
	c.dataList.Refresh()

	if c.listFilter != nil {
		c.applyFilter(c.listFilter.Text)
		c.listFilter.Enable()
	} else {
		c.applyFilter("")
	}
}

// SetPlaceholder sets the placeholder text to be shown
func (c *ClickableList[V]) SetPlaceholder(placeholderText string) {
	if c.listFilter != nil {
		c.listFilter.SetPlaceHolder(placeholderText)
	}
}

func (c *ClickableList[V]) applyFilter(textFilter string) {
	var currentData []V
	for _, datum := range c.allData {
		if strings.Contains(strings.ToLower(c.nameGetter(datum)), strings.ToLower(textFilter)) {
			currentData = append(currentData, datum)
		}
	}
	c.currentData = currentData
	c.dataList.Refresh()
}
