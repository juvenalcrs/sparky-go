package sparky

import "fyne.io/fyne/v2"

// Renderer defines sparky renderer interface.
type Renderer interface {
	// CreateContent creates the render content.
	// This function is called only the first time the render is created.
	CreateContent() *fyne.Container
	// Destroy destroys the render content.
	Destroy()
	// Refresh refreshes the render content.
	Refresh()
}

// CreateRenderer creates a widget renderer from a sparky.Renderer.
func CreateRenderer(r Renderer) fyne.WidgetRenderer {
	return &rendererImpl{
		render:  r,
		objects: []fyne.CanvasObject{r.CreateContent()},
	}
}

// rendererImpl should implement fyne.WidgetRenderer
var _ (fyne.WidgetRenderer) = (*rendererImpl)(nil)

// rendererImpl defines a simple renderer implementation.
type rendererImpl struct {
	render  Renderer
	objects []fyne.CanvasObject
}

// Destroy implements fyne.WidgetRenderer.
func (r *rendererImpl) Destroy() {
	r.render.Destroy()
}

// Layout implements fyne.WidgetRenderer.
func (r *rendererImpl) Layout(size fyne.Size) {
	r.objects[0].Resize(size)
}

// MinSize implements fyne.WidgetRenderer.
func (r *rendererImpl) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

// Objects implements fyne.WidgetRenderer.
func (r *rendererImpl) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh implements fyne.WidgetRenderer.
func (r *rendererImpl) Refresh() {
	r.render.Refresh()
}
