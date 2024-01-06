package ui

import "fmt"

var version = ""

func (app *App) getVersion() string {
	return fmt.Sprintf("%s-%d", app.Metadata().Version, app.Metadata().Build)
}
