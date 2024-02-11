package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showError(errorText string) {
	errorLbl := widget.NewLabel(errorText)
	quitBtn := widget.NewButton("Quit", func() { os.Exit(1) })

	content := container.New(layout.NewVBoxLayout(), errorLbl, quitBtn)

	displayError(content)
}

func displayError(content *fyne.Container) {
	mainWin.SetContent(content)
}

func startUpError(errorText string) fyne.CanvasObject {
	errorLbl := widget.NewLabel(errorText)
	quitBtn := widget.NewButton("Quit", func() { os.Exit(1) })

	return container.New(layout.NewVBoxLayout(), errorLbl, quitBtn)
}
