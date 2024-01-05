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

func (a *App) MainView() {
	a.Window.Show()
	tabs := container.NewAppTabs()

	logsicon := theme.ListIcon()
	settingsicon := theme.SettingsIcon()
	healthicon := theme.HelpIcon()

	// the terminal output
	terminal := createTerminalOutput(a)
	// drop cache button
	dropCacheButton := createDropCacheButton(a)
	// start/stop button
	startStopButton := createStopStartButton(a, terminal, dropCacheButton)

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
				a.StopServer()
				time.Sleep(10 * time.Second)
				a.StartServer()
			},
			a.peerID,
		),
	))

	IsAutoStarted := a.IsAutoStarted()
	var onLoginButton *widget.Button

	if IsAutoStarted {
		onLoginButton = widget.NewButton("Disable autostart", func() {
			if err := a.UninstallAutoStart(); err != nil {
				dialog.ShowError(err, a.Window)
				return
			}
			onLoginButton.SetText("Enable autostart")
		})
	} else {
		onLoginButton = widget.NewButton("Enable autostart", func() {
			if err := a.InstallAutoStart(); err != nil {
				dialog.ShowError(err, a.Window)
				return
			}
			onLoginButton.SetText("Disable autostart")
		})
	}

	tabs.Append(container.NewTabItemWithIcon(
		"Settings",
		settingsicon,
		container.NewBorder(
			nil, onLoginButton,
			nil, nil,
			components.Settings(
				a.Options,
				func() {
					a.Save()
				},
			)),
	))

	// make the version string
	version := fmt.Sprintf("%s-%d", a.Metadata().Version, a.Metadata().Build)
	aboutContent = strings.ReplaceAll(aboutContent, "{{version}}", version)

	rt := widget.NewRichTextFromMarkdown(aboutContent)
	rt.Wrapping = fyne.TextWrapWord
	tabs.Append(container.NewTabItemWithIcon(
		"About",
		theme.InfoIcon(),
		container.NewVScroll(rt),
	))

	a.Window.SetContent(tabs)
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
				a.StopServer()
				setStart()
			case false:
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
