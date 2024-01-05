package healthview

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

const reloadInterval = 10 * time.Second

type HealthView struct {
	widget.BaseWidget
	status        *petalsserver.Status
	models        []*ModelHealth
	container     *fyne.Container
	refreshButton *widget.Button
	title         *canvas.Text
	peerID        *string
	onProblem     func()
}

func NewHealthView(status *petalsserver.Status, onProblem func(), peerID *string) *HealthView {
	h := &HealthView{
		status:    status,
		container: container.NewVBox(),
		title:     canvas.NewText("", theme.ForegroundColor()),
		peerID:    peerID,
		onProblem: onProblem,
	}
	h.ExtendBaseWidget(h)

	h.title.TextStyle = fyne.TextStyle{Monospace: false, Bold: true}
	h.title.Alignment = fyne.TextAlignCenter
	h.title.TextSize = theme.TextSize() * 1.25

	remaining := reloadInterval / time.Second
	refreshChan := make(chan struct{})
	refreshButton := widget.NewButtonWithIcon(
		fmt.Sprintf("Refresh (auto in %ds)", remaining),
		theme.ViewRefreshIcon(),
		func() {
			refreshChan <- struct{}{}
		},
	)
	h.refreshButton = refreshButton
	h.LoadStatus()

	go func() {
		reload := func(force bool) {
			defer refreshButton.Enable()
			refreshButton.SetText(fmt.Sprintf("Refresh (auto in %ds)", remaining))
			if force || remaining <= 0 {
				refreshButton.Disable()
				refreshButton.SetText("Refreshing...")
				h.LoadStatus()
				remaining = reloadInterval / time.Second
				time.Sleep(1 * time.Second)
				refreshButton.SetText(fmt.Sprintf("Refresh (auto in %ds)", remaining))
			}
		}

		for {
			select {
			case <-refreshChan:
				reload(true)
				remaining = reloadInterval / time.Second
			case <-time.Tick(1 * time.Second):
				remaining--
				reload(false)
			}
		}
	}()
	return h
}

func (h *HealthView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			h.title,
			h.refreshButton,
			nil,
			nil,
			container.NewVScroll(h.container),
		),
	)
}

func (h *HealthView) LoadStatus() {
	var stateColors = map[string]color.Color{
		"healthy": theme.SuccessColor(),
		"broken":  theme.ErrorColor(),
	}
	status, err := petalsserver.GetStatus()
	if err != nil {
		log.Printf("error getting status: %v", err)
		return
	}

	// detect problem on the server
	if h.peerID != nil {
		for _, model := range status.ModelReports {
			for _, row := range model.ServerRows {
				if row.State == "unreachable" && row.PeerID == *h.peerID {
					h.onProblem()
					break
				}
			}
		}
	}

	h.container.RemoveAll()
	now := time.Now()
	h.title.Text = "Status refreshed on " + now.Format("2006-01-02 15:04:05")

	h.status = status
	h.models = make([]*ModelHealth, len(h.status.ModelReports))
	for i, model := range h.status.ModelReports {
		modelNameTitle := canvas.NewText(model.Name, theme.ForegroundColor())
		modelNameTitle.TextStyle = fyne.TextStyle{Monospace: false, Bold: true}

		statusTitle := canvas.NewText(model.State, stateColors[model.State])
		statusTitle.TextStyle = fyne.TextStyle{Monospace: false, Bold: false}

		h.container.Add(modelNameTitle)
		h.container.Add(statusTitle)

		h.models[i] = NewModelHealth(model.ServerRows, model.NumBlocks, h.peerID)
		h.container.Add(h.models[i])
	}
	h.container.Refresh()
	h.status = status
	h.Refresh()
}

var _ fyne.WidgetRenderer = (*healthRenderer)(nil)
var _ fyne.Widget = (*ModelHealth)(nil)

type ModelHealth struct {
	widget.BaseWidget
	rows          []petalsserver.Peer
	numBlocks     int
	currentPeerID *string
}

func NewModelHealth(rows []petalsserver.Peer, numBlocks int, current *string) *ModelHealth {
	h := &ModelHealth{
		rows:          rows,
		numBlocks:     numBlocks,
		currentPeerID: current,
	}
	h.ExtendBaseWidget(h)
	return h
}

func (h *ModelHealth) CreateRenderer() fyne.WidgetRenderer {
	return newHealthRenderer(h.rows, h.numBlocks, h.currentPeerID)
}

var _ fyne.WidgetRenderer = (*healthRenderer)(nil)

type healthRenderer struct {
	rows       []petalsserver.Peer
	numBlocks  int
	squareSize float32
	labelWidth float32
	height     float32
	padding    float32
	objects    []fyne.CanvasObject
}

func newHealthRenderer(rows []petalsserver.Peer, numBlocks int, currentID *string) *healthRenderer {
	h := &healthRenderer{
		rows:       rows,
		numBlocks:  numBlocks,
		squareSize: 5,
		padding:    2,
		objects:    []fyne.CanvasObject{},
	}
	var (
		y              float32
		maxPeerIDWidth float32
	)
	for i := range h.rows {
		row := h.rows[i]
		peerStatus := NewPeerRow(&row, h.numBlocks, currentID)
		peerStatus.Move(fyne.NewPos(0, y))
		y += peerStatus.MinSize().Height + h.padding
		maxPeerIDWidth = func() float32 {
			width := peerStatus.LabelSize().Width
			if width > maxPeerIDWidth {
				return width
			}
			return maxPeerIDWidth
		}()

		h.objects = append(h.objects, peerStatus)
	}

	h.labelWidth = maxPeerIDWidth
	h.height = y

	return h
}

func (h *healthRenderer) Destroy() {
}

// Layout the components of the widget.
//
// Implements WidgetRenderer.Layout()
func (h *healthRenderer) Layout(size fyne.Size) {
	for i := range h.objects {
		obj := h.objects[i]
		obj.Resize(fyne.NewSize(size.Width, obj.MinSize().Height))
		if row, ok := obj.(*PeerRow); ok {
			row.SetXpos(h.labelWidth)
			obj.Refresh()
		}
	}
}

func (h *healthRenderer) MinSize() fyne.Size {
	return fyne.Size{
		Width:  h.labelWidth + float32(h.numBlocks)*(h.squareSize+h.padding),
		Height: (h.padding*2 + h.height),
	}
}

func (h *healthRenderer) Refresh() {
	for i := range h.objects {
		obj := h.objects[i]
		obj.Refresh()
	}
}

func (h *healthRenderer) Objects() []fyne.CanvasObject {
	return h.objects
}
