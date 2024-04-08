package ui

import (
	"fmt"
	"net"
	"playlog/logic"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoadSetupUI(MyApp logic.MyApp) fyne.CanvasObject {
	questionLbl := widget.NewLabel("Select which mode app will run in, Local or Server Sync")
	localBtn := widget.NewButton("Local Setup", func() {
		MyApp.App.Preferences().SetString("Mode", "Local")
		MyApp.App.Preferences().SetBool("FirstRun", false)
		MyApp.Win.SetContent(LoadMainUI(MyApp))
	})

	serverIpLbl := widget.NewLabel("Enter Server IP Address")
	serverIpEnt := widget.NewEntry()
	serverIpEnt.SetPlaceHolder("Enter Server IP")
	serverPortLbl := widget.NewLabel("Enter Server Port below")
	serverPortEnt := widget.NewEntry()
	serverPortEnt.SetPlaceHolder("Default is 7529")
	serverBtn := widget.NewButton("Server Sync", func() {
		var port string
		var err []string
		ip := serverIpEnt.Text
		if serverPortEnt.Text == "" {
			port = "7529"
		} else if _, err := strconv.Atoi(serverPortEnt.Text); err == nil {
			port = serverPortEnt.Text
		}

		if ip == "" {
			err = append(err, "IP Empty")
		}
		if len(err) == 0 && net.ParseIP(ip) == nil {
			err = append(err, fmt.Sprintf("%s is not a valid IP", ip))
		}
		if len(err) == 0 && !logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", ip, port)) {
			err = append(err, fmt.Sprintf("Server with details: %s:%s is inaccessible", ip, port))
		}

		if len(err) != 0 {
			ShowError(strings.Join(err[:], "\n\n"), MyApp)
		} else {
			logic.ServerSetup(ip, port, "Sync", MyApp)
			MyApp.App.Preferences().SetBool("FirstRun", false)
			MyApp.Win.SetContent(LoadMainUI(MyApp))
		}

	})

	return container.New(layout.NewVBoxLayout(), questionLbl, localBtn, serverIpLbl, serverIpEnt, serverPortLbl, serverPortEnt, serverBtn)
}
