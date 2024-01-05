package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

func (a *App) UpdateView() {
	a.Window.Show()
	label := widget.NewLabel("An update is available. Would you like to install it?")
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter

	updateButton := widget.NewButtonWithIcon("Update", theme.ConfirmIcon(), func() {
		progress := widget.NewProgressBarInfinite()
		label := widget.NewLabel("Downloading update...")
		label.Wrapping = fyne.TextWrapWord
		label.Alignment = fyne.TextAlignCenter
		box := container.NewBorder(
			nil,
			progress,
			nil,
			nil,
			label,
		)
		a.Window.SetContent(box)

		if err := petalsserver.UpdatePetals(); err != nil {
			label.SetText("Error downloading update: " + err.Error())
			return
		}
		label.SetText("Update downloaded")
		nextButton := widget.NewButtonWithIcon("Continue", theme.ConfirmIcon(), func() {
			a.MainView()
		})
		nextButton.Importance = widget.HighImportance
		a.Window.SetContent(container.NewBorder(
			nil,
			nextButton,
			nil,
			nil,
			label,
		))

	})
	updateButton.Importance = widget.HighImportance

	skipButton := widget.NewButtonWithIcon("Skip", theme.NavigateNextIcon(), func() {
		a.MainView()
	})
	skipButton.Importance = widget.WarningImportance
	a.Window.SetContent(container.NewBorder(
		nil,
		container.NewGridWithColumns(2,
			updateButton,
			skipButton,
		),
		nil,
		nil,
		label,
	))
}
