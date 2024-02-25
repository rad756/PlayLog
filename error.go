package main

import (
	"net"

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
	backBtn := widget.NewButton("OK", func() { errorPpu.Hide() })

	content := container.New(layout.NewVBoxLayout(), topContent, errorLbl, backBtn)

	errorPpu = widget.NewModalPopUp(content, mainWin.Canvas())
	errorPpu.Show()
}

func startUpServerError() fyne.CanvasObject {
	topLbl := widget.NewLabel("-- Startup Error --")
	topContent := container.New(layout.NewCenterLayout(), topLbl)
	errorLbl := widget.NewLabel("Server with IP " + serverIP + " is inaccessible\nThe app will start in Server Down Mode and try to sync upon next startup")
	offlineModeBtn := widget.NewButton("Enter Server Down Mode", func() {
		serverDownMode = true
		serverMode = false

		writeConfig()

		mainWin.SetContent(loadMainMenuUI())
	})

	return container.New(layout.NewVBoxLayout(), topContent, errorLbl, offlineModeBtn)
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
			showError(changeIPEnt.Text + " is not valid IP")
		} else if !isServerAccessible("http://" + changeIPEnt.Text + ":" + changePortEnt.Text) {
			showError("Server with details: " + changeIPEnt.Text + ":" + port + " is inaccessible")
		} else {
			serverMode = true
			serverDownMode = false
			serverIP = changeIPEnt.Text
			serverPort = changePortEnt.Text
			writeConfig()
		}
	})

	orLbl := widget.NewLabel("OR")
	centeredOrLbl := container.NewCenter(orLbl)

	backBtn := widget.NewButton("Enter Server Down Mode", func() {
		serverDownMode = true
		serverMode = false

		writeConfig()

		errorPpu.Hide()
	})
	backLbl := widget.NewLabel("Current change will not be pushed.\nTry it again after swithching to Server Down Mode")

	content := container.New(layout.NewVBoxLayout(), errorLbl, changeLbl, changeIPEnt, changePortEnt, changeServerBtn, centeredOrLbl, backBtn, backLbl)

	errorPpu = widget.NewModalPopUp(content, mainWin.Canvas())
	errorPpu.Show()
}
