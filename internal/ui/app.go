package ui

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

const (
	SettingPublicName  = "public_name"
	SettingMaxDiskSize = "max_disk_size"
	SettingNumBlocks   = "num_blocks"
	SettingAutoStart   = "auto_start"
	//SettingThreshold     = "threshold"
	SettingStopOnProcess = "process_names"
)

// embed icons, logos
var (
	//go:embed logo-green.png
	greenIconPngBytes []byte
	greenIconFile     fyne.Resource

	//go:embed logo-red.png
	redIconPngBytes []byte
	redIconFile     fyne.Resource
)

func init() {
	// init the icons as static resources
	greenIconFile = fyne.NewStaticResource("logo-green.png", greenIconPngBytes)
	redIconFile = fyne.NewStaticResource("logo-red.png", redIconPngBytes)
}

// App is the main application of FlowerPot.
type App struct {
	fyne.App
	Window          fyne.Window
	menu            *fyne.Menu
	Options         *petalsserver.RunOptions
	cancelServer    context.CancelFunc
	started         bool
	manuallyStopped bool
	peerID          *string
	onStarted       func(chan []byte, error)
	onStopped       func()
	currentCmd      *exec.Cmd
}

// NewApp creates a new application
func NewApp(id string) *App {

	app := &App{
		App:     app.NewWithID(id),
		started: false,
		peerID:  new(string),
	}

	app.DesktopIntegration()

	prefs := app.Preferences()
	app.Options = &petalsserver.RunOptions{
		PublicName:  prefs.StringWithFallback(SettingPublicName, ""),
		MaxDiskSize: prefs.IntWithFallback(SettingMaxDiskSize, 0),
		NumBlocks:   prefs.IntWithFallback(SettingNumBlocks, 0),
		AutoStart:   prefs.BoolWithFallback(SettingAutoStart, false),
		//Threshold:   app.Preferences().FloatWithFallback(SettingThreshold, 0.20),
		StopOnProcess: app.Preferences().StringListWithFallback(SettingStopOnProcess, []string{
			"SteamLaunch",
			"Blender",
		}),
	}

	// the window
	app.Window = app.NewWindow("FlowerPot for Petals Server")
	app.SetIcon(redIconFile) // the officon icon is red

	// do not close the window, just hide it
	// as we want to keep the server running and
	// the systray icon alive
	app.Window.SetCloseIntercept(app.Window.Hide)

	// default size and hide the window (we will show it later)
	app.Window.Resize(fyne.NewSize(800, 600))

	// hide by default
	app.Window.Hide()

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

	return app

}

// Run the application
func (a *App) Run() {
	// check if petals server is installed, if not, show the presentation view
	if !petalsserver.IsPetalsServerInstalled() {
		a.PresentationView()
		a.App.Run()
		return
	}

	// check if the installed version is the same as the git version
	thisVersion, err := petalsserver.GetInstalledGITSHA()
	if err != nil {
		log.Println("Error getting installed version: ", err)
	} else {
		gitVersion, err := petalsserver.GetLatestGitCommitSHA()
		if err != nil {
			log.Println("Error getting git version: ", err)
		} else if thisVersion != gitVersion {
			a.UpdateView()
			a.App.Run()
			return
		}
	}

	// normal startup
	a.ServerLoop()
	a.MainView()
	a.App.Run()
}

// Save the preferences.
func (a *App) Save() error {
	prefs := a.Preferences()
	prefs.SetString(SettingPublicName, a.Options.PublicName)
	prefs.SetInt(SettingMaxDiskSize, a.Options.MaxDiskSize)
	prefs.SetInt(SettingNumBlocks, a.Options.NumBlocks)
	prefs.SetBool(SettingAutoStart, a.Options.AutoStart)
	prefs.SetStringList(SettingStopOnProcess, a.Options.StopOnProcess)

	return nil
}

func (a *App) Quit() {
	a.StopServer()
	a.App.Quit()
}

func (a *App) StartServer() error {
	var err error
	var cmd *exec.Cmd
	outchan := make(chan []byte, 1024)
	a.cancelServer, cmd, err = petalsserver.LaunchPetalsServer(a.Options, outchan)
	if err != nil {
		return err
	}
	a.started = true
	a.currentCmd = cmd

	a.SetMenuEntry("Stop", func() {
		a.StopServer()
	})

	if a.onStarted != nil {
		a.onStarted(outchan, err)
		return err
	}
	return nil
}

func (a *App) StopServer() {

	if a.currentCmd != nil {
		log.Println("Stopping server")
		a.currentCmd.Process.Signal(os.Interrupt)
	}

	if a.cancelServer != nil {
		a.cancelServer()
		// wait for the process to exit
		if err := a.currentCmd.Wait(); err != nil {
			log.Printf("Process finished with error: %v", err)
		}
	} else {
		log.Println("cancelServer is nil")
	}

	// BUG: this is a hack to avoid the server to remain running in background
	if petalsserver.IsPetalsServerRunning() {
		defer petalsserver.ForceKill() // TODO: avoid this
		log.Println("server still running, force kill")
	}
	a.started = false
	*a.peerID = ""
	a.SetMenuEntry("Start", func() {
		a.StartServer()
	})
	if a.onStopped != nil {
		a.onStopped()
	}
}

func (a *App) StartIfNeeded() {
	if a.started {
		log.Println("server already started")
		return
	}

	// the server should start if the auto start option is enabled
	// and if the user did not manually stop the server
	if !a.Options.AutoStart || a.manuallyStopped {
		log.Println("auto start disabled")
		return
	}

	// is the process already running?
	if petalsserver.IsPetalsServerRunning() {
		log.Println("server already running")
		return
	}

	if ok, reason := canStart(a.Options); !ok {
		log.Println("server should not start because: ", reason)
		return
	}

	a.StartServer()
}

func (a *App) StopServerIfNeeded() {
	if !a.started {
		return
	}

	if !petalsserver.IsPetalsServerRunning() {
		return
	}

	if ok, reason := canStart(a.Options); !ok {
		log.Println("stopping server because: ", reason)
		a.StopServer()
	}
}

func (a *App) ServerLoop() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !a.started {
				a.StartIfNeeded()
			} else {
				a.StopServerIfNeeded()
			}
		}
	}()
}

func canStart(options *petalsserver.RunOptions) (bool, string) {

	status := petalsserver.GPUStatus
	for _, process := range status.GPU[0].Processes {
		for _, stopOnProcess := range options.StopOnProcess {
			if strings.Contains(process.Name, stopOnProcess) {
				return false, "process " + process.Name + " is running"
			}
		}
	}

	computeProcess := petalsserver.GetComputeProcessCount()
	if computeProcess > 1 {
		return false, "more than one compute process"
	}

	return true, "That's OK"
}

func (a *App) SetMenuEntry(label string, action func()) {

	if desk, ok := a.App.(desktop.App); ok {
		log.Println("Setting menu entry: ", label)
		switch label {
		case "Start":
			desk.SetSystemTrayIcon(greenIconFile)
		case "Stop":
			desk.SetSystemTrayIcon(redIconFile)
		}
	}

	a.menu.Items[0] = fyne.NewMenuItem(
		label,
		action,
	)
	a.menu.Refresh()
}
