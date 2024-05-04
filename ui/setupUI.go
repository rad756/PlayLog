package ui

import (
	"fmt"
	"net"
	"playlog/logic"
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoadSetupUI(MyApp logic.MyApp) {
	questionLbl := widget.NewLabel("Select which mode app will run in, Local or Server Sync")
	localBtn := widget.NewButton("Local Setup", func() {
		MyApp.App.Preferences().SetString("StorageMode", "Local")
		MyApp.App.Preferences().SetBool("FirstRun", false)
		LoadMainUI(MyApp)
	})

	serverIpLbl := widget.NewLabel("Enter Server IP Address")
	serverIpEnt := widget.NewEntry()
	serverIpEnt.SetPlaceHolder("Enter Server IP")
	serverPortLbl := widget.NewLabel("Enter Server Port below")
	serverPortEnt := widget.NewEntry()
	serverPortEnt.SetPlaceHolder("Default is 7529")
	serverBtn := widget.NewButton("Server Sync", func() {
		var port string
		var errStr []string
		ip := serverIpEnt.Text
		if serverPortEnt.Text == "" {
			port = "7529"
		} else if _, errStr := strconv.Atoi(serverPortEnt.Text); errStr == nil {
			port = serverPortEnt.Text
		}

		if ip == "" {
			errStr = append(errStr, "IP Empty")
		}
		if len(errStr) == 0 && net.ParseIP(ip) == nil {
			errStr = append(errStr, fmt.Sprintf("%s is not a valid IP", ip))
		}
		if len(errStr) == 0 && !logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", ip, port)) {
			errStr = append(errStr, fmt.Sprintf("Server with details: %s:%s is inaccessible", ip, port))
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
		} else {
			logic.ServerSetup(ip, port, "Sync", MyApp)
			MyApp.App.Preferences().SetBool("FirstRun", false)
			LoadMainUI(MyApp)
		}

	})

	content := container.New(layout.NewVBoxLayout(), questionLbl, localBtn, serverIpLbl, serverIpEnt, serverPortLbl, serverPortEnt, serverBtn)

	MyApp.Win.SetContent(content)
}
