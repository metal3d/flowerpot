package healthview

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

type PeerRow struct {
	widget.BaseWidget
	row       *petalsserver.Peer
	numBlocks int
	renderer  fyne.WidgetRenderer
	peerID    *string
}

func NewPeerRow(row *petalsserver.Peer, numBlocks int, currentPeerID *string) *PeerRow {
	p := &PeerRow{
		row:       row,
		numBlocks: numBlocks,
		peerID:    currentPeerID,
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *PeerRow) CreateRenderer() fyne.WidgetRenderer {
	r := newPeerRowRenderer(p.row, p.numBlocks)
	r.currentPeerID = p.peerID
	p.renderer = r
	return r
}

func (p *PeerRow) SetXpos(x float32) {
	if r, ok := p.renderer.(*peerRowRenderer); ok {
		r.setXpos(x)
	}
}

func (p *PeerRow) LabelSize() fyne.Size {
	if p.renderer == nil {
		return fyne.Size{}
	}
	if r, ok := p.renderer.(*peerRowRenderer); ok {
		return r.peerLabel.LabelSize()
	}
	return fyne.Size{}
}

type peerRowRenderer struct {
	row           *petalsserver.Peer
	numBlocks     int
	squareSize    float32
	padding       float32
	peerLabel     *PeerLabel
	blocks        *BlockView
	background    *canvas.Rectangle
	blocksXpos    float32
	currentPeerID *string
}

func newPeerRowRenderer(row *petalsserver.Peer, numBlocks int) *peerRowRenderer {
	r := &peerRowRenderer{
		row:        row,
		numBlocks:  numBlocks,
		squareSize: 5,
		padding:    2,
		peerLabel:  NewPeerLabel(row, numBlocks),
		background: canvas.NewRectangle(color.Transparent),
	}

	height := r.peerLabel.MinSize().Height
	r.blocks = NewBlockView(row, numBlocks, height)

	return r
}

func (p *peerRowRenderer) Destroy() {}

func (p *peerRowRenderer) Layout(size fyne.Size) {
	// if this line corresponds to the current peer, color it
	if p.row.PeerID == *p.currentPeerID {
		p.background.FillColor = theme.DisabledColor()
	}
	p.background.Resize(size)

	if p.blocksXpos == 0 {
		p.blocksXpos = p.peerLabel.MinSize().Width
	}
	p.blocks.Move(fyne.NewPos(p.blocksXpos, 0))
}

func (p *peerRowRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		p.peerLabel.MinSize().Width+p.blocks.MinSize().Width,
		p.peerLabel.MinSize().Height,
	)
}

func (p *peerRowRenderer) Refresh() {
	p.background.Refresh()
	p.peerLabel.Refresh()
	p.blocks.Refresh()
}

func (p *peerRowRenderer) setXpos(x float32) {
	p.blocksXpos = x
}

func (p *peerRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{p.background, p.peerLabel, p.blocks}
}
