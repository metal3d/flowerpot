package components

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TerminalOutput struct {
	widget.BaseWidget
	logChan  chan []byte
	label    *widget.Label
	scroll   *container.Scroll
	onOutput func()
	serverID string
}

func NewTerminalOutput(onOutput func()) *TerminalOutput {
	t := &TerminalOutput{
		onOutput: onOutput,
	}

	t.label = widget.NewLabel("")
	t.label.TextStyle = fyne.TextStyle{Monospace: true, TabWidth: 2}
	t.label.Wrapping = fyne.TextWrapWord

	t.scroll = container.NewVScroll(t.label)

	t.ExtendBaseWidget(t)
	return t
}

func (t *TerminalOutput) StartLogs(logChan chan []byte) {
	go t.writeLogs(logChan)
}

func (t *TerminalOutput) SetText(text string) {
	t.label.SetText(text)
}

func (t *TerminalOutput) Text() string {
	return t.label.Text
}

func (t *TerminalOutput) writeLogs(logs chan []byte) {
	lines := []byte{}
	for bytes := range logs {
		newLine := false
		for _, b := range bytes {
			switch b {
			case '\r':
				if lastNewLine := strings.LastIndex(string(lines), "\r"); lastNewLine > 0 {
					lines = lines[:lastNewLine]
				}
			case '\n':
				newLine = true
				// ensure we don't remove the logs on new textual process bar
				bytes = append(bytes, '\r')
			case '\b', ']':
				// we manage the textual progress bar
				newLine = true
			}
		}
		lines = append(lines, bytes...)
		if newLine {
			t.label.SetText(string(lines))
			t.scroll.ScrollToBottom()
			if t.onOutput != nil {
				t.onOutput()
			}
		}

	}
	t.label.SetText("Server stopped")
}

func (t *TerminalOutput) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.scroll)
}
