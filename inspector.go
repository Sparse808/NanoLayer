package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"fmt"
)

type Inspector struct {
	widget.BaseWidget
	selectedRect *DraggableRect
	selected     bool
}

func NewInspector() *Inspector {
	insp := &Inspector{
		selected: false,
	}

	insp.ExtendBaseWidget(insp)
	return insp
}

type inspectorRenderer struct {
	insp    *Inspector
	objects []fyne.CanvasObject
}

func (insp *Inspector) CreateRenderer() fyne.WidgetRenderer {
	title := widget.NewLabel("Inspector")
	title.Alignment = fyne.TextAlignCenter

	if !insp.selected {
		blank := widget.NewLabel("No Rect Selected")
		c := container.NewBorder(title, nil, nil, nil, blank)
		return widget.NewSimpleRenderer(c)
	}
	selectedName := widget.NewLabel("Name: " + fmt.Sprintf("%p", insp.selectedRect))

	name := container.NewHBox(selectedName)

	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, nil, nil))
}

func (r *inspectorRenderer) Layout(size fyne.Size) {
	// Fill the entire widget space

}

func (r *inspectorRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *inspectorRenderer) Refresh() {

}

func (r *inspectorRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *inspectorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *inspectorRenderer) Destroy() {}
