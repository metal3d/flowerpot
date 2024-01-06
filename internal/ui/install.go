package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/coreos/go-semver/semver"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

var (
	//go:embed need-python-message.md
	needPythonMessage string
	//go:embed need-nvidia-smi.md
	needNvidiaSMIMessage string
)

// InstallView shows the installation progress.
func (a *App) InstallView() {
	progress := widget.NewProgressBarInfinite()

	// check if python is <3.12
	pythonpath, err := findPythonLesserThan("3.12.0")
	if err != nil {
		a.Window.Show()
		label := widget.NewRichTextFromMarkdown(needPythonMessage)
		label.Wrapping = fyne.TextWrapWord
		backbutton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
			a.PresentationView()
		})
		a.Window.SetContent(container.NewBorder(
			nil,
			backbutton,
			nil,
			nil,
			label,
		))
		return
	}

	// check if nvidia-smi is installed
	_, err = findNvidiaSMI()
	if err != nil {
		a.Window.Show()
		label := widget.NewRichTextFromMarkdown(needNvidiaSMIMessage)
		label.Wrapping = fyne.TextWrapWord
		backbutton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
			a.PresentationView()
		})
		a.Window.SetContent(container.NewBorder(
			nil,
			backbutton,
			nil,
			nil,
			label,
		))
		return
	}

	label := widget.NewLabel("Installing Petals Server using" + pythonpath)
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}

	output := widget.NewLabel("")
	output.Importance = widget.HighImportance
	output.TextStyle = fyne.TextStyle{Monospace: true}
	outscroll := container.NewVScroll(output)
	output.Wrapping = fyne.TextWrapWord

	nextbutton := widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), func() {
		a.MainView()
	})

	box := container.NewBorder(
		label,
		progress,
		nil,
		nil,
		outscroll,
	)

	a.Window.SetContent(box)
	a.Window.Show()
	go func() {
		out, err := petalsserver.PipInstallPetals(pythonpath)
		if err != nil {
			log.Println("error installing petals server:", err)
		}
		for line := range out {
			output.SetText(output.Text + line + "\n")
			outscroll.ScrollToBottom()
		}
		progress.Stop()
		box = container.NewBorder(
			label,
			nextbutton,
			nil,
			nil,
			outscroll,
		)
		a.Window.SetContent(box)
	}()
}

func findPythonLesserThan(maxversion string) (string, error) {
	python := []string{
		"python3.11",
		"python3.10",
		"python3",
		"python",
	}
	for _, pyexec := range python {
		pythonpath, err := exec.LookPath(pyexec)
		if err != nil {
			continue
		}
		dest := pythonpath
		// find the real path of the python executable
		for err == nil {
			dest, err = os.Readlink(dest)
			if err != nil {
				continue
			}
			dest = filepath.Join(
				filepath.Dir(pythonpath),
				dest,
			)
			pythonpath = dest
		}
		py := "import sys;print(f'{sys.version_info.major}.{sys.version_info.minor}.{sys.version_info.micro}')"
		cmd := exec.Command(pythonpath, "-c", py)
		out, err := cmd.Output()
		if err != nil {
			log.Println("err", err)
			continue
		}
		version := strings.TrimSpace(string(out))
		installedVersion := semver.New(version)
		targetVersion := semver.New(maxversion)
		if installedVersion.LessThan(*targetVersion) {
			return pythonpath, nil
		}
	}
	return "", fmt.Errorf("no python found")
}

func findNvidiaSMI() (string, error) {
	nvidiasmipath, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return "", err
	}
	return nvidiasmipath, nil
}
