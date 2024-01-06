package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func (app *App) SetupSystray() {
	// the systray icon is not shown by default
	app.menu = fyne.NewMenu("File",
		fyne.NewMenuItem("...", func() {}), // will be replaced by the start/stop menu entry with SetMenuEntry
		fyne.NewMenuItem(
			"Show interface",
			func() {
				app.Window.Show()
			}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(
			"Quit",
			func() {
				app.Quit()
			},
		),
	)
	if desk, ok := app.App.(desktop.App); ok {
		desk.SetSystemTrayMenu(app.menu)
	}

	app.SetMenuEntry("Start", func() {
		app.StartServer()
	})

}
func (app *App) SetMenuEntry(label string, action func()) {

	if desk, ok := app.App.(desktop.App); ok {
		log.Println("Setting menu entry: ", label)
		switch label {
		case "Start":
			desk.SetSystemTrayIcon(greenIconFile)
		case "Stop":
			desk.SetSystemTrayIcon(redIconFile)
		}
	}

	app.menu.Items[0] = fyne.NewMenuItem(
		label,
		action,
	)
	app.menu.Refresh()
}
