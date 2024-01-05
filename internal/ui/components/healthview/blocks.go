package healthview

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

const (
	statusOnline      = "online"
	statusJoining     = "joining"
	statusUnreachable = "unreachable"
	statusOffline     = "offline"
)

var _ fyne.Widget = (*BlockView)(nil)
var _ fyne.WidgetRenderer = (*blockViewRenderer)(nil)

type BlockView struct {
	widget.BaseWidget
	row        *petalsserver.Peer
	numBlocks  int
	lineHeight float32
}

func NewBlockView(row *petalsserver.Peer, numBlocks int, lineHeight float32) *BlockView {
	b := &BlockView{
		row:        row,
		numBlocks:  numBlocks,
		lineHeight: lineHeight,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *BlockView) CreateRenderer() fyne.WidgetRenderer {
	return &blockViewRenderer{
		row:        b.row,
		numBlocks:  b.numBlocks,
		lineHeight: b.lineHeight,
		squareSize: 5,
		padding:    2,
		container:  container.NewWithoutLayout(),
	}
}

type blockViewRenderer struct {
	row        *petalsserver.Peer
	numBlocks  int
	container  *fyne.Container
	squareSize float32
	lineHeight float32
	padding    float32
}

func (b *blockViewRenderer) Destroy() {}

func (b *blockViewRenderer) Layout(size fyne.Size) {
	b.container.RemoveAll()

	row := b.row
	numBlocks := b.numBlocks
	padding := b.padding

	var posX float32

	start := row.Span.ServerInfo.StartBlock
	end := row.Span.ServerInfo.EndBlock

	color := theme.ForegroundColor()
	switch row.State {
	case statusOnline:
		color = theme.SuccessColor()
	case statusJoining:
		color = theme.WarningColor()
	case statusUnreachable:
		color = theme.ErrorColor()
	case statusOffline:
		color = theme.DisabledColor()
	}

	for i := 0; i < numBlocks; i++ {
		if i >= start && i < end { // warning: end is not included!
			rect := canvas.NewCircle(color)
			rect.FillColor = color
			rect.Resize(fyne.NewSize(b.squareSize, b.squareSize))
			rect.Move(fyne.NewPos(posX, b.lineHeight/2))
			b.container.Add(rect)
		} else {
			// draw a simple dash
			line := canvas.NewLine(theme.DisabledColor())
			line.StrokeWidth = 1
			line.Position1 = fyne.NewPos(
				posX,
				b.lineHeight/2+b.squareSize/2,
			)
			line.Position2 = fyne.NewPos(
				posX+b.squareSize,
				b.lineHeight/2+b.squareSize/2,
			)
			b.container.Add(line)
		}
		posX += b.squareSize + padding
	}
}

func (b *blockViewRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		(b.squareSize+b.padding*2)*float32(b.numBlocks),
		b.lineHeight,
	)
}

func (b *blockViewRenderer) Refresh() {
	b.container.Refresh()
}

func (b *blockViewRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{b.container}
}
