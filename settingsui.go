package main

import (
	"net"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var currentModeLbl *widget.Label
var currentIPLbl *widget.Label
var currentPortLbl *widget.Label

func makeSettingsTab() fyne.CanvasObject {
	var changeIPEnt *widget.Entry
	var changePortEnt *widget.Entry
	currentModeLbl = widget.NewLabel("Current Mode: ")

	switchModeBtn := widget.NewButton("Switch Mode", func() {
		if !serverMode { //switch from local to server sync
			if serverIP == "" || net.ParseIP(serverIP) == nil {
				showError("Server Settings Incorrect\nChange Server Details")
			} else if !isServerAccessible("http://" + serverIP + ":" + serverPort) {
				showError("Server Inaccessible\nChange Server Details")
			} else if fileConflictCheck() {
				mainWin.SetContent(loadSyncUI())
				currentModeLbl.SetText("Current Mode: Server Sync")
			} else {
				serverMode = true
				serverDownMode = false
				currentModeLbl.SetText("Current Mode: Server Sync")
				writeConfig()
			}
		} else { // switch from server sync to local
			serverMode = false
			currentModeLbl.SetText("Current Mode: Local")
			writeConfig()
		}
	})

	currentIPLbl = widget.NewLabel("Current IP: " + serverIP)
	currentPortLbl = widget.NewLabel("Current Port: " + serverPort)

	changeServerLbl := widget.NewLabel("Change Server Settings")
	centeredChangeServerLbl := container.NewCenter(changeServerLbl)
	changeIPEnt = widget.NewEntry()
	changeIPEnt.PlaceHolder = "New Server IP"
	changePortEnt = widget.NewEntry()
	changePortEnt.PlaceHolder = "New Server Port, If Empty Defaults To 7529"
	changeServerBtn := widget.NewButton("Change Server", func() {
		var err []string
		var port string

		if changeIPEnt.Text == "" || net.ParseIP(changeIPEnt.Text) == nil {
			err = append(err, "IP Empty or Invalid")
		}
		if _, err1 := strconv.Atoi(changePortEnt.Text); err1 != nil && changePortEnt.Text != "" {
			err = append(err, "Invalid Port")
		} else if _, err2 := strconv.Atoi(changePortEnt.Text); err2 == nil {
			port = changePortEnt.Text
		} else {
			port = "7529"
		}

		if len(err) != 0 {
			showError(strings.Join(err[:], "\n\n"))
		} else if !isServerAccessible("http://" + changeIPEnt.Text + ":" + port) {
			showError("Server: " + changeIPEnt.Text + ":" + port + " is inaccesslible")
		} else {
			serverIP = changeIPEnt.Text
			serverPort = port
			currentIPLbl.SetText("Current IP: " + serverIP)
			currentPortLbl.SetText("Current Port: " + serverPort)
			writeConfig()
		}
	})

	return container.NewVBox(currentModeLbl, switchModeBtn, currentIPLbl, currentPortLbl, centeredChangeServerLbl, changeIPEnt, changePortEnt, changeServerBtn)
}
