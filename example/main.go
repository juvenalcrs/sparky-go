package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fpabl0/sparky-go"
	"github.com/fpabl0/sparky-go/scont"
)

const (
	keyFirstName sparky.ValueKey = iota
)

func main() {
	a := app.New()
	// a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("hello")
	w.Resize(fyne.NewSize(500, 400))

	ctx := sparky.NewContext(w)
	ctx.PutValue(keyFirstName, "Pablo")

	l := ctx.ShowLoader("Loading...")
	go func() {
		time.Sleep(1 * time.Second)
		<-l.Error("Un simple error muy pero muy largo debe de salir mal ahora!!!!!!")
		ctx.ShowInfo("Atención!", "Hay una nueva versión disponible")
		firstName := ctx.GetValue(keyFirstName).(string)
		fmt.Printf("Listo %s\n", firstName)
	}()

	list := widget.NewList(
		func() int {
			return 10
		},
		func() fyne.CanvasObject {
			return newTitleItem()
		},
		func(index widget.ListItemID, obj fyne.CanvasObject) {
			tile := obj.(*tileItem)
			tile.Name = fmt.Sprintf("Name %d", index)
			tile.LastName = fmt.Sprintf("LastName %d", index)
			tile.Cellphone = "Cellphone"
			tile.Refresh()
		},
	)

	toggle := false
	w.SetContent(
		scont.NewFrame(2, 2,
			container.NewBorder(widget.NewButton("Press", func() {
				if toggle {
					ctx.ShowSuccess("Listo!", "El conductor se ha creado")
				} else {
					go func() {
						resp := <-ctx.ShowPasswordInput("Nueva contraseña", "Ingrese su nueva contraseña", "Confirmar")
						if resp == nil {
							return
						}
						fmt.Println(*resp)
					}()
				}
				toggle = !toggle
			}), nil, nil, nil, list),
		),
	)

	w.ShowAndRun()
}

// ===============================================================
// Title Item
// ===============================================================

type tileItem struct {
	widget.BaseWidget
	Name      string
	LastName  string
	Cellphone string
}

func newTitleItem() *tileItem {
	p := &tileItem{}
	p.ExtendBaseWidget(p)
	return p
}

func (p *tileItem) CreateRenderer() fyne.WidgetRenderer {
	p.ExtendBaseWidget(p)
	return sparky.CreateRenderer(&titleItemRenderer{widget: p})
}

type titleItemRenderer struct {
	name      *canvas.Text
	lastname  *canvas.Text
	cellphone *canvas.Text

	widget *tileItem
}

func (r *titleItemRenderer) CreateContent() *fyne.Container {
	r.name = canvas.NewText(r.widget.Name, theme.ForegroundColor())
	r.name.TextSize = theme.CaptionTextSize()
	r.lastname = canvas.NewText(r.widget.LastName, theme.ForegroundColor())
	r.lastname.TextSize = theme.CaptionTextSize()
	r.cellphone = canvas.NewText(r.widget.Cellphone, theme.ForegroundColor())
	r.cellphone.TextSize = theme.CaptionTextSize()

	return container.NewVBox(
		r.name,
		r.lastname,
		r.cellphone,
	)
}

func (r *titleItemRenderer) Destroy() {}

func (r *titleItemRenderer) Refresh() {
	r.name.Text = r.widget.Name
	r.name.Refresh()

	r.lastname.Text = r.widget.LastName
	r.lastname.Refresh()

	r.cellphone.Text = r.widget.Cellphone
	r.cellphone.Refresh()
}
