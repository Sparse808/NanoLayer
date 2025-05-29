package main

import (
	"image/color"

	"fyne.io/fyne/v2"
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
	return &inspectorRenderer{
		inspector: insp,
		objects:   []fyne.CanvasObject{&insp.title, &insp.rectPos},
	}
}

func (r *inspectorRenderer) Layout(size fyne.Size) {
	// Fill the entire widget space
	r.objects[0].Move(fyne.NewPos(0, 0))              // title
	r.objects[0].Resize(fyne.NewSize(size.Width, 20)) // resize as needed

	r.objects[1].Move(fyne.NewPos(0, 25)) // rectPos
	r.objects[1].Resize(fyne.NewSize(size.Width, 20))

}

func (r *inspectorRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *inspectorRenderer) Refresh() {
	if r.inspector.rect != nil {
		r.inspector.rectPos.SetText(fmt.Sprintf("(%.0f, %.0f)", r.inspector.rect.Position().X, r.inspector.rect.Position().Y))
		r.objects[1].Refresh()
	}

}

func (r *inspectorRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *inspectorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *inspectorRenderer) Destroy() {}
