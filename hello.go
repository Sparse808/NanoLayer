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

	//inspector section

	inspector := NewInspector()
	edit := newEditor(inspector)

	editorLayout := container.NewWithoutLayout(edit)

	//side panel layout
	toolsSection := container.NewVBox()

	//toolsSection
	title := widget.NewLabel("Tools")
	but1 := widget.NewButton("box", func() {
		fmt.Println("pressed")

		editorLayout.Objects = append(editorLayout.Objects, NewDraggableRect(70, 70, *edit.background, inspector))
	})

	toolsSection.Objects = append(toolsSection.Objects, title)
	toolsSection.Objects = append(toolsSection.Objects, but1)

	sidePanel := container.NewHBox(layout.NewSpacer(), toolsSection, inspector)

	myLayout := container.NewHBox(layout.NewSpacer(), editorLayout, layout.NewSpacer(), sidePanel)

	w.SetContent(myLayout)
	w.Resize(fyne.NewSize(700, 500))
	w.ShowAndRun()
}

type editor struct {
	widget.BaseWidget
	background *canvas.Rectangle
	sb         *selectedBorder
	inspector  *Inspector
}

func newEditor(insp *Inspector) *editor {
	newEditor := &editor{
		background: canvas.NewRectangle(color.Black),
		sb:         newSelectedBorder(),
		inspector:  insp,
	}

	newEditor.background.Resize(fyne.NewSize(100, 300))
	newEditor.ExtendBaseWidget(newEditor)
	return newEditor
}

func (edit *editor) CreateRenderer() fyne.WidgetRenderer {
	return &editorRenderer{
		edit:    edit,
		objects: []fyne.CanvasObject{edit.background, &edit.sb.top},
	}
}

type editorRenderer struct {
	edit    *editor
	objects []fyne.CanvasObject
}

func (r *editorRenderer) Layout(size fyne.Size) {
	if r.edit.inspector.rect != nil {
		r.edit.sb.top.Position1 = r.edit.inspector.rect.Position()
		r.edit.sb.top.Position2 = fyne.NewPos(r.edit.inspector.rect.Position().X, r.edit.inspector.rect.Position().Y+r.edit.inspector.rect.Size().Height)
		r.edit.sb.top.StrokeColor = color.White
		r.edit.sb.top.StrokeWidth = 5
		r.Refresh()
	}

}

func (r *editorRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 300)
}

func (r *editorRenderer) Refresh() {

}

func (r *editorRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *editorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *editorRenderer) Destroy() {}

// selected border code
type selectedBorder struct {
	top   canvas.Line
	bot   canvas.Line
	left  canvas.Line
	right canvas.Line
}

func (sb *selectedBorder) turnOff() {
	sb.top.Hide()
	sb.bot.Hide()
	sb.left.Hide()
	sb.right.Hide()
}

func (sb *selectedBorder) turnOn() {
	sb.top.Show()
	sb.bot.Show()
	sb.left.Show()
	sb.right.Show()
}

func (sb *selectedBorder) setPosition(rect fyne.CanvasObject) {
	sb.top.Move(rect.Position())
	//sb.bot.Move(fyne.NewPos(rect.Position().X, rect.Position().Y+rect.Size().Width))
	//sb.left.Move(rect.Position())
	//sb.right.Move(fyne.NewPos(rect.Position().X+rect.MinSize().Width, rect.Position().Y))
}

func (sb *selectedBorder) setLengths(rect fyne.CanvasObject) {
	//rectX := rect.Position().X
	//rectY := rect.Position().Y
	//sb.top.Resize(fyne.NewPos(rectX, rectY + rect.Size().Width))
	sb.top.Resize(fyne.NewSize(rect.MinSize().Width, 10))
	//sb.bot.Resize(fyne.NewPos(rectX + rect.Size().Height, rectY + rect.Size().Width))

}

func newSelectedBorder() *selectedBorder {
	newsb := &selectedBorder{
		top: *canvas.NewLine(color.White),
	}

	//newsb.ExtendBaseWidget(newsb)

	return newsb
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
