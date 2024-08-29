package ui

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

const spinnerMillis = 200

type UI struct {
	spinner *spinner.Spinner
	output  *os.File
	verbose bool
}

var infoFormat = color.New(color.FgHiBlue).Sprint
var erroFormat = color.New(color.FgRed).Sprint

func NewUI(output *os.File, verbose bool) *UI {
	ui := &UI{
		spinner: spinner.New(spinner.CharSets[11], spinnerMillis*time.Millisecond, spinner.WithWriterFile(output)),
		output:  output,
		verbose: verbose,
	}

	log.SetOutput(output)

	_ = ui.spinner.Color("white", "bold")
	ui.spinner.Prefix = " "
	ui.spinner.Start()

	return ui
}

func (ui *UI) Update(msg string) {
	if ui.verbose {
		ui.Info(msg)
	} else {
		ui.spinner.Suffix = " " + msg
	}
}

func (ui *UI) println(msg string) {
	ui.spinner.Stop()
	fmt.Fprintln(ui.output, msg)
	ui.spinner.Start()
}

func (ui *UI) Info(msg string) {
	ui.println(" " + infoFormat(msg))
}

func (ui *UI) Error(msg string) {
	ui.println(" " + erroFormat(msg))
}

func (ui *UI) Stop() {
	ui.spinner.Disable()
}
