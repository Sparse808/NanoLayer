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
	inspectorTitle := widget.NewLabel("Inspector")
	inspectedRect := widget.NewLabel("rect")

	//side panel layout
	toolsSection := container.NewVBox()
	inspecterSection := container.NewVBox(inspectorTitle, inspectedRect)

	//toolsSection
	title := widget.NewLabel("Tools")
	but1 := widget.NewButton("box", func() {
		fmt.Println("pressed")

		noLayout.Objects = append(noLayout.Objects, NewDraggableRect(70, 70, *gridBackground, inspectedRect))
	})

	toolsSection.Objects = append(toolsSection.Objects, title)
	toolsSection.Objects = append(toolsSection.Objects, but1)

	sidePanel := container.NewHBox(layout.NewSpacer(), toolsSection, inspecterSection)

	myLayout := container.NewHBox(layout.NewSpacer(), noLayout, layout.NewSpacer(), sidePanel)

	w.SetContent(myLayout)
	w.Resize(fyne.NewSize(700, 500))
	w.ShowAndRun()
}

type inspector struct {
	selected bool
	title    widget.Label
	rect     DraggableRect
	rectPos  widget.Label
}

func NewInspector() *inspector {
	insp := &inspector{
		selected: false,
		title:    *widget.NewLabel("object"),
		rectPos:  *widget.NewLabel("00"),
	}

	insp.rect.ExtendBaseWidget()
	return insp
}

type DraggableRect struct {
	widget.BaseWidget
	rect        *canvas.Rectangle
	gridGrounds canvas.Rectangle
	inspector   *widget.Label
}

func NewDraggableRect(x, y float32, grounds canvas.Rectangle, inspector *widget.Label) *DraggableRect {
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
	r.inspector.SetText(fmt.Sprintf("%p", r))
}

// Dragged interface implementations

func (r *DraggableRect) Dragged(e *fyne.DragEvent) {
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
}

func (r *DraggableRect) DragEnd() {}
