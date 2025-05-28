package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"fmt"
)

func main() {
	a := app.New()
	w := a.NewWindow("Draggable Rectangle")

	//screen background
	gridBackground := canvas.NewRectangle(color.Black)
	gridBackground.Resize(fyne.NewSize(100, 300))

	noLayout := container.NewWithoutLayout(gridBackground)

	//inspector section
	inspector := NewInspector()

	//side panel layout
	toolsSection := container.NewVBox()

	//toolsSection
	title := widget.NewLabel("Tools")
	but1 := widget.NewButton("box", func() {
		fmt.Println("pressed")

		noLayout.Objects = append(noLayout.Objects, NewDraggableRect(70, 70, *gridBackground, inspector))
	})

	toolsSection.Objects = append(toolsSection.Objects, title)
	toolsSection.Objects = append(toolsSection.Objects, but1)

	sidePanel := container.NewHBox(layout.NewSpacer(), toolsSection, inspector)

	myLayout := container.NewHBox(layout.NewSpacer(), noLayout, layout.NewSpacer(), sidePanel)

	w.SetContent(myLayout)
	w.Resize(fyne.NewSize(700, 500))
	w.ShowAndRun()
}

// Inspector Code
type Inspector struct {
	widget.BaseWidget
	selected bool
	title    widget.Label
	rect     *DraggableRect
	rectPos  widget.Label
}

func NewInspector() *Inspector {
	insp := &Inspector{
		selected: false,
		title:    *widget.NewLabel("object"),
		rectPos:  *widget.NewLabel("00"),
	}

	insp.ExtendBaseWidget(insp)
	return insp
}

func (insp *Inspector) CreateRenderer() fyne.WidgetRenderer {
	return &inspectorRenderer{
		inspector: insp,
		objects:   []fyne.CanvasObject{&insp.title, &insp.rectPos},
	}
}

type inspectorRenderer struct {
	inspector *Inspector
	objects   []fyne.CanvasObject
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

// DraggableRect widget code
type DraggableRect struct {
	widget.BaseWidget
	rect        *canvas.Rectangle
	gridGrounds canvas.Rectangle
	inspector   *Inspector
}

func NewDraggableRect(x, y float32, grounds canvas.Rectangle, inspector *Inspector) *DraggableRect {
	dr := &DraggableRect{
		rect:        canvas.NewRectangle(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
		gridGrounds: grounds,
		inspector:   inspector,
	}
	dr.ExtendBaseWidget(dr)
	dr.Resize(fyne.NewSize(80, 60))
	dr.Move(fyne.NewPos(x, y))
	return dr
}

// Renderer
func (r *DraggableRect) CreateRenderer() fyne.WidgetRenderer {
	// The rectangle is always positioned at (0,0) inside the widget
	return &draggableRectRenderer{
		rect:    r,
		objects: []fyne.CanvasObject{r.rect},
	}
}

type draggableRectRenderer struct {
	rect    *DraggableRect
	objects []fyne.CanvasObject
}

func (r *draggableRectRenderer) Layout(size fyne.Size) {
	// Fill the entire widget space
	r.rect.rect.Resize(size)
	r.rect.rect.Move(fyne.NewPos(0, 0))
}

func (r *draggableRectRenderer) MinSize() fyne.Size {
	return r.rect.Size()
}

func (r *draggableRectRenderer) Refresh() {
	r.rect.rect.FillColor = color.NRGBA{R: 255, G: 100, B: 100, A: 255}
	r.rect.rect.Refresh()
}

func (r *draggableRectRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *draggableRectRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *draggableRectRenderer) Destroy() {}

func (r *DraggableRect) Tapped(e *fyne.PointEvent) {
	//r.inspector.SetText(fmt.Sprintf("%p", r))
	r.inspector.rectPos.SetText(fmt.Sprintf("(%.0f, %.0f)", r.Position().X, r.Position().Y))
	r.inspector.rect = r
	r.inspector.selected = true
	r.inspector.Refresh()
}

// Dragged interface implementations

func (r *DraggableRect) Dragged(e *fyne.DragEvent) {
	if r.inspector.rect != r {
		return
	}
	y := r.gridGrounds.Position().Y + r.gridGrounds.Size().Height
	x := r.gridGrounds.Position().X + r.gridGrounds.Size().Width

	newPos := r.Position().Add(e.Dragged)

	if newPos.X+r.Size().Width > x {
		newPos.X = x - r.Size().Width
	} else if newPos.X < r.gridGrounds.Position().X {
		newPos.X = r.gridGrounds.Position().X
	}
	if newPos.Y+r.Size().Height > y {
		newPos.Y = y - r.Size().Height
	} else if newPos.Y < r.gridGrounds.Position().Y {
		newPos.Y = r.gridGrounds.Position().Y
	}

	r.Move(newPos) // Move the widget itself
	r.inspector.Refresh()
}

func (r *DraggableRect) DragEnd() {}
