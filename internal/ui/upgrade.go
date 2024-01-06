package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"fyne.io/fyne/v2/dialog"
)

const gitRelease = "https://api.github.com/repos/metal3d/flowerpot/releases/latest"

var forbidden = []string{
	"alpha",
	"beta",
	"rc",
	"dev",
	"test",
	"pre",
	"preview",
}

func (app *App) checkLatestVersion() {
	resp, err := http.Get(gitRelease)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Latest version: ", release.TagName)
	for _, f := range forbidden {
		if strings.Contains(release.TagName, f) {
			return
		}
	}
	version := app.getVersion()
	version = "v" + version
	if release.TagName != version {
		dialog.ShowInformation(
			"New version available",
			fmt.Sprintf("A new version is available on GitHub: %s", release.TagName),
			app.Window,
		)
	}
}
