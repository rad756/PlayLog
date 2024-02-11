package main

import (
	"bufio"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showError(errorText string) {
	var errorPpu *widget.PopUp

	topLbl := widget.NewLabel("-- Error(s) --")
	topContent := container.New(layout.NewCenterLayout(), topLbl)
	errorLbl := widget.NewLabel(errorText)
	quitBtn := widget.NewButton("OK", func() { errorPpu.Hide() })

	content := container.New(layout.NewVBoxLayout(), topContent, errorLbl, quitBtn)

	errorPpu = widget.NewModalPopUp(content, mainWin.Canvas())
	errorPpu.Show()
}

func displayError(content *fyne.Container) {
	errorPpu := widget.NewModalPopUp(content, mainWin.Canvas())
	errorPpu.Show()
}

func startUpError(errorText string) fyne.CanvasObject {
	errorLbl := widget.NewLabel(errorText)
	quitBtn := widget.NewButton("Quit", func() { os.Exit(1) })

	return container.New(layout.NewVBoxLayout(), errorLbl, quitBtn)
}

func showServerInaccessibleError() {
	var errorPpu *widget.PopUp

	errorLbl := widget.NewLabel("Server with IP " + serverIP + ":" + serverPort + " is inaccessible")
	changeLbl := widget.NewLabel("To correct server details, type below")
	changeIPEnt := widget.NewEntry()
	changeIPEnt.PlaceHolder = "Type IP"
	changePortEnt := widget.NewEntry()
	changePortEnt.PlaceHolder = "Type Port / default 7529"
	var port string

	if changePortEnt.Text == "" {
		port = "7529"
	} else {
		port = changePortEnt.Text
	}

	changeServerBtn := widget.NewButton("Change Server Details", func() {
		if changeIPEnt.Text == "" {
			showError("IP empty")
		} else if net.ParseIP(changeIPEnt.Text) == nil {
			showError(changeIPEnt.Text + " is invalid IP")
		} else if !isServerAccessible("http://" + changeIPEnt.Text + ":" + changePortEnt.Text) {
			showError("Server with details: " + changeIPEnt.Text + ":" + port + " is inaccessible")
		} else {
			configFile := "conf.csv"
			file, err := os.Create(configFile)

			if err != nil {
				panic(err)
			} else {
				writer := bufio.NewWriter(file)

				writer.WriteString("mode,sync\n")
				writer.WriteString("ip," + changeIPEnt.Text + "\n")
				writer.WriteString("port," + port + "\n")

				writer.Flush()

			}
		}
	})

	quitBtn := widget.NewButton("Quit", func() { os.Exit(1) })

	content := container.New(layout.NewVBoxLayout(), errorLbl, changeLbl, changeIPEnt, changePortEnt, changeServerBtn, quitBtn)

	errorPpu = widget.NewModalPopUp(content, mainWin.Canvas())
	errorPpu.Show()
}
