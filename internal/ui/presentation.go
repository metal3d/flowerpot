package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed welcome.md
var welcomeMessage string

// PresentationView shows the welcome message. It is the first view if the
// application is launched for the first time (no petals server installtion found).
func (a *App) PresentationView() {
	a.Window.Show()
	label := widget.NewRichTextFromMarkdown(welcomeMessage)
	label.Wrapping = fyne.TextWrapWord

	okbutton := widget.NewButtonWithIcon("Yes", theme.NavigateNextIcon(), func() {
		a.InstallView()
	})
	okbutton.Importance = widget.HighImportance

	cancelbutton := widget.NewButtonWithIcon("No", theme.CancelIcon(), func() {
		a.Quit()
	})
	cancelbutton.Importance = widget.DangerImportance

	logo := canvas.NewImageFromResource(redIconFile)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(128, 128))

	a.Window.SetContent(container.NewBorder(
		nil,
		container.NewGridWithColumns(2,
			okbutton,
			cancelbutton,
		),
		nil, nil,
		container.NewBorder(
			logo,
			nil, nil, nil,
			label),
	))
}
