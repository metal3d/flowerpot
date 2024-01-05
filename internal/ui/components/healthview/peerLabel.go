package healthview

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

type PeerLabel struct {
	widget.BaseWidget
	row       *petalsserver.Peer
	numBlocks int
	renderer  fyne.WidgetRenderer
}

func NewPeerLabel(row *petalsserver.Peer, numBlocks int) *PeerLabel {
	p := &PeerLabel{
		row:       row,
		numBlocks: numBlocks,
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *PeerLabel) CreateRenderer() fyne.WidgetRenderer {
	r := newPeerLabelRenderer(p.numBlocks, p.row)
	p.renderer = r
	return r
}

func (p *PeerLabel) LabelSize() fyne.Size {
	if p.renderer == nil {
		return fyne.Size{}
	}
	if r, ok := p.renderer.(*peerLabelRenderer); ok {
		return r.idLabel.MinSize()
	}
	return fyne.Size{}
}

type peerLabelRenderer struct {
	row       *petalsserver.Peer
	numBlocks int
	idLabel   *widget.RichText
}

func newPeerLabelRenderer(numBlocks int, row *petalsserver.Peer) *peerLabelRenderer {
	p := &peerLabelRenderer{
		row:       row,
		numBlocks: numBlocks,
	}

	p.idLabel = widget.NewRichText()
	peer := "..." + row.PeerID[len(row.PeerID)-6:]
	md := "`" + peer + "`"
	if row != nil && row.Span.ServerInfo.PublicName != "" {
		u, err := url.Parse(row.Span.ServerInfo.PublicName)
		if err == nil && u.Hostname() != "" {
			md = fmt.Sprintf("%s :: [%s](%s)", md, u.Hostname(), u.String())
		} else {
			md = fmt.Sprintf("%s :: %s", md, row.Span.ServerInfo.PublicName)
		}
	}
	p.idLabel.ParseMarkdown(md)
	return p
}

func (p *peerLabelRenderer) Destroy() {}

func (p *peerLabelRenderer) Layout(size fyne.Size) {
	p.idLabel.Resize(size)
}

func (p *peerLabelRenderer) MinSize() fyne.Size {
	return p.idLabel.MinSize()
}

func (p *peerLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{p.idLabel}
}

func (p *peerLabelRenderer) Refresh() {
	p.idLabel.Refresh()
}
