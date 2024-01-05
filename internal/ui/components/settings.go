package components

import (
	"fmt"
	"strings"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/metal3d/flowerpot/internal/petalsserver"
)

//go:embed settings.md
var settingsMessage string

var models = []string{
	"petals-team/StableBeluga2",
	"tiiuae/falcon-180B-chat",
	"meta-llama/Llama-2-70b-chat-hf",
	"meta-llama/Llama-2-70b-hf",
	"bigscience/bloom-560m",
}

func Settings(options *petalsserver.RunOptions, onchange func(), setAtLogin func(), getAtStartup func() bool) fyne.CanvasObject {
	// install flowerpot at login
	getLabel := func() string {
		atlogin := "Start Flowerpot at login"
		if getAtStartup() {
			atlogin = "Do not start Flowerpot at login"
		}
		return atlogin
	}
	var installAtLogin *widget.Button
	installAtLogin = widget.NewButton(getLabel(), func() {
		setAtLogin()
		installAtLogin.SetText(getLabel())
	})

	nameEntry := widget.NewEntry()
	nameEntry.SetText(options.PublicName)
	atStartup := widget.NewCheck("", func(value bool) {
		// TODO
	})
	atStartup.SetChecked(options.AutoStart)

	explanation := widget.NewRichTextFromMarkdown(settingsMessage)
	explanation.Wrapping = fyne.TextWrapWord

	maxDiskSize := widget.NewSlider(0, 1024)
	maxDiskSize.Step = 5
	maxdiskLabel := widget.NewLabel(fmt.Sprintf("%d GiB", options.MaxDiskSize))
	maxDiskSize.OnChanged = func(value float64) {
		if value < 1 {
			maxdiskLabel.SetText("auto")
			return
		}
		maxdiskLabel.SetText(fmt.Sprintf("%d GiB", int(value)))
	}
	maxDiskSize.SetValue(float64(options.MaxDiskSize))
	maxDiskSize.OnChanged(float64(options.MaxDiskSize))

	numBlocks := widget.NewSlider(0, 1024)
	numBlocksLabel := widget.NewLabel(fmt.Sprintf("%d", options.NumBlocks))
	numBlocks.OnChanged = func(value float64) {
		if value < 1 {
			numBlocksLabel.SetText("auto")
			return
		}
		numBlocksLabel.SetText(fmt.Sprintf("%d", int(value)))
	}
	numBlocks.SetValue(float64(options.NumBlocks))
	numBlocks.OnChanged(float64(options.NumBlocks))

	modelList := widget.NewSelect(models, func(value string) {
		options.ModelName = value
	})
	if options.ModelName == "" {
		options.ModelName = models[0]
	}
	modelList.SetSelected(options.ModelName)

	processList := widget.NewMultiLineEntry()
	processList.SetText(strings.Join(options.StopOnProcess, ","))

	form := widget.NewForm(
		widget.NewFormItem("Public name", nameEntry),
		widget.NewFormItem("Model to serve", modelList),
		widget.NewFormItem("Max disk size", container.NewBorder(
			nil, nil, nil, maxdiskLabel, maxDiskSize,
		)),
		widget.NewFormItem("Number of blocks", container.NewBorder(
			nil, nil, nil, numBlocksLabel, numBlocks,
		)),
		//widget.NewFormItem("Memory threshold", container.NewBorder(
		//	nil, nil, nil, memorySliderLabel, memorySlider,
		//)),
		widget.NewFormItem("Stop on process", processList),
		widget.NewFormItem("Launch at startup", atStartup),
		widget.NewFormItem("Click to launch Flowerpot at desktop login", installAtLogin),
	)

	form.OnSubmit = func() {
		options.StopOnProcess = strings.Split(processList.Text, ",")
		for i, name := range options.StopOnProcess {
			options.StopOnProcess[i] = strings.TrimSpace(name)
		}
		options.PublicName = nameEntry.Text
		options.AutoStart = atStartup.Checked
		options.MaxDiskSize = int(maxDiskSize.Value)
		options.NumBlocks = int(numBlocks.Value)
		//options.Threshold = memorySlider.Value
		onchange()
	}
	form.SubmitText = "Save"

	return container.NewBorder(
		form,
		nil,
		nil,
		nil,
		container.NewVScroll(explanation),
	)
}
