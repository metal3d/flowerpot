package components

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

const (
	lowMemoryUsage    float64 = 0.0
	mediumMemoryUsage         = 75.0
	highMemoryUsage           = 95.0
)

var _ fyne.Widget = (*GPULoad)(nil)

type GPULoad struct {
	widget.BaseWidget
	renderer fyne.WidgetRenderer
}

func NewGPULoad() *GPULoad {
	g := &GPULoad{}
	g.ExtendBaseWidget(g)
	g.stat()
	return g
}

func (g *GPULoad) CreateRenderer() fyne.WidgetRenderer {
	r := newGPULoadRenderer()
	g.renderer = r
	return r
}

func (g *GPULoad) stat() {
	go func() {
		for {
			usage := petalsserver.GetFreeMemory()
			if g.renderer != nil {
				g.renderer.(*gpuLoadRenderer).setUsage(usage)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

var _ fyne.WidgetRenderer = (*gpuLoadRenderer)(nil)

type gpuLoadRenderer struct {
	level      *canvas.Rectangle
	background *canvas.Rectangle
	usage      float64
	text       *canvas.Text
	colors     map[float64]color.Color
}

func newGPULoadRenderer() *gpuLoadRenderer {
	g := &gpuLoadRenderer{
		level:      canvas.NewRectangle(color.RGBA{0, 0, 0, 0}),
		background: canvas.NewRectangle(theme.DisabledColor()),
	}

	t := fyne.CurrentApp().Settings().Theme()
	variant := fyne.CurrentApp().Settings().ThemeVariant()
	g.colors = map[float64]color.Color{
		highMemoryUsage:   t.Color(theme.ColorRed, variant),
		mediumMemoryUsage: t.Color(theme.ColorOrange, variant),
		lowMemoryUsage:    t.Color(theme.ColorGreen, variant),
	}

	g.level.CornerRadius = 8

	return g
}

func (g *gpuLoadRenderer) Destroy() {}

func (g *gpuLoadRenderer) Layout(size fyne.Size) {

	// make the background fill the whole widget
	g.background.Resize(size)
	g.background.Move(fyne.NewPos(0, 0))

	// make the rectangle "usage" percent of the width
	g.level.Resize(fyne.NewSize(size.Width*float32(g.usage/100), 20))
	if g.usage > highMemoryUsage {
		g.level.FillColor = g.colors[highMemoryUsage]
	} else if g.usage > mediumMemoryUsage {
		g.level.FillColor = g.colors[mediumMemoryUsage]
	} else {
		g.level.FillColor = g.colors[lowMemoryUsage]
	}

	g.level.Move(fyne.NewPos(0, 0))

	g.text = canvas.NewText(
		fmt.Sprintf("GPU memory usage %.2f %%", g.usage),
		theme.ForegroundColor(),
	)
	g.text.Alignment = fyne.TextAlignCenter
	g.text.Move(fyne.NewPos(size.Width/2, 0))
}

func (g *gpuLoadRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 20)
}

func (g *gpuLoadRenderer) Refresh() {
	g.level.Refresh()
}

func (g *gpuLoadRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{g.background, g.level, g.text}
}

func (g *gpuLoadRenderer) setUsage(usage float64) {
	g.usage = usage
	g.Refresh()
}
