package ui

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
	"github.com/metal3d/flowerpot/internal/ui/components"
	"github.com/metal3d/flowerpot/internal/ui/components/healthview"
)

//go:embed about.md
var aboutContent string

func (app *App) MainView() {
	app.Window.Show()
	app.SetupSystray()

	tabs := container.NewAppTabs()

	logsicon := theme.ListIcon()
	settingsicon := theme.SettingsIcon()
	healthicon := theme.HelpIcon()

	// the terminal output
	terminal := createTerminalOutput(app)
	// drop cache button
	dropCacheButton := createDropCacheButton(app)
	// start/stop button
	startStopButton := createStopStartButton(app, terminal, dropCacheButton)

	status, err := petalsserver.GetStatus()
	if err != nil {
		log.Println(err)
	}

	// gpu load widget
	gpuLoad := components.NewGPULoad()

	// create tabs and add them to the window
	tabs.Append(container.NewTabItemWithIcon(
		"Logs",
		logsicon,
		container.NewBorder(
			gpuLoad,
			container.NewVBox(
				startStopButton,
				dropCacheButton,
			),
			nil,
			nil,
			terminal,
		),
	))

	tabs.Append(container.NewTabItemWithIcon(
		"Swarm Health",
		healthicon,
		healthview.NewHealthView(
			status,
			func() {
				log.Println("Problem detected, restarting server")
				terminal.SetText("Problem detected, restarting server in 10 seconds\n")
				app.StopServer()
				time.Sleep(10 * time.Second)
				app.StartServer()
			},
			app.peerID,
		),
	))

	tabs.Append(container.NewTabItemWithIcon(
		"Settings",
		settingsicon,
		components.Settings(
			app.Options,
			func() { // onchange callback
				app.Save()
			},
			func() { // set at login button callback
				isAutoStarted := app.IsAutoStarted()
				if isAutoStarted {
					if err := app.UninstallAutoStart(); err != nil {
						dialog.ShowError(err, app.Window)
						return
					}
				} else {
					if err := app.InstallAutoStart(); err != nil {
						dialog.ShowError(err, app.Window)
						return
					}
				}
			},
			app.IsAutoStarted, // the function to check if the app is auto started
		),
	),
	)

	// make the version string
	version := fmt.Sprintf("%s-%d", app.Metadata().Version, app.Metadata().Build)
	aboutContent = strings.ReplaceAll(aboutContent, "{{version}}", version)

	rt := widget.NewRichTextFromMarkdown(aboutContent)
	rt.Wrapping = fyne.TextWrapWord
	tabs.Append(container.NewTabItemWithIcon(
		"About",
		theme.InfoIcon(),
		container.NewVScroll(rt),
	))

	app.Window.SetContent(tabs)
}

func createTerminalOutput(a *App) *components.TerminalOutput {
	var terminal *components.TerminalOutput
	terminal = components.NewTerminalOutput(
		func() { // on output, check for peer ID if it is not already set
			if a.peerID != nil && *a.peerID != "" {
				return
			}
			reg := regexp.MustCompile(`\[.*\/p2p\/(.*)',`)
			text := terminal.Text()
			lines := strings.Split(text, "\n")
			for _, line := range lines {
				matches := reg.FindStringSubmatch(line)
				if len(matches) > 1 {
					log.Println("Found peer ID: ", matches[1])
					*a.peerID = matches[1]
				}
			}
		},
	)
	return terminal
}

func createStopStartButton(a *App, terminal *components.TerminalOutput, dropCacheButton *widget.Button) *widget.Button {
	// start/stop button
	var startStopButton *widget.Button

	setStart := func() {
		dropCacheButton.Enable()
		dropCacheButton.SetText("Drop Cache")
		startStopButton.SetIcon(theme.MediaPlayIcon())
		startStopButton.SetText("Start")
		startStopButton.Importance = widget.SuccessImportance
		startStopButton.Refresh()
	}

	setStop := func() {
		dropCacheButton.Disable()
		dropCacheButton.SetText("Drop Cache (server is running)")
		startStopButton.SetIcon(theme.MediaStopIcon())
		startStopButton.SetText("Stop")
		startStopButton.Importance = widget.DangerImportance
		startStopButton.Refresh()
	}

	startStopButton = widget.NewButtonWithIcon(
		"Start",
		theme.MediaPlayIcon(),
		func() {
			switch a.started {
			case true:
				a.manuallyStopped = true
				a.StopServer()
				setStart()
			case false:
				a.manuallyStopped = false
				a.StartServer()
				setStop()
			}
		},
	)
	setStart()

	// register callbacks
	a.onStarted = func(output chan []byte, err error) {
		setStop()
		if err != nil {
			log.Println(err)
			terminal.SetText("Error starting server: " + err.Error())
			return
		}
		terminal.StartLogs(output)
	}

	a.onStopped = func() {
		setStart()
	}
	return startStopButton
}

func createDropCacheButton(a *App) *widget.Button {
	dropCacheButton := widget.NewButtonWithIcon(
		"Drop Cache",
		theme.DeleteIcon(),
		func() {
			dialog.ShowConfirm(
				"Drop Cache",
				"Are you sure you want to drop the cache?",
				func(ok bool) {
					if !ok {
						return
					}
					progress := dialog.NewCustomWithoutButtons(
						"Dropping Cache",
						container.NewBorder(
							widget.NewLabel("Please wait..."),
							nil, nil, nil,
							widget.NewProgressBarInfinite(),
						),
						a.Window,
					)
					defer progress.Hide()
					progress.Show()
					if err := petalsserver.EmptyCache(); err != nil {
						log.Println(err)
					}
					time.Sleep(1 * time.Second)
				},
				a.Window,
			)
		})
	dropCacheButton.Importance = widget.LowImportance

	return dropCacheButton
}
