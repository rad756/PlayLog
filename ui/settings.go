package ui

import (
	"fmt"
	"net"
	"playlog/logic"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func MakeSettingsTab(MyApp logic.MyApp) fyne.CanvasObject {
	var changeIPEnt *widget.Entry
	var changePortEnt *widget.Entry
	ModeBind := binding.BindPreferenceString("Mode", MyApp.App.Preferences())

	currentModeLbl := widget.NewLabelWithData(binding.NewSprintf("Current mode: %s", ModeBind))

	switchModeBtn := widget.NewButton("Switch Mode", func() {
		if MyApp.App.Preferences().String("Mode") == "Sync" {
			MyApp.App.Preferences().SetString("Mode", "Local")
			return
		}

		if MyApp.App.Preferences().String("Mode") == "Local" && logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
			if !logic.FileConflictCheck(MyApp) {
				MyApp.App.Preferences().SetString("Mode", "Sync")
				return
			} else {
				LoadSyncUI(MyApp)
				return
			}
		} else {
			ShowError("Cannot connect to server, check details or if server is running", MyApp)
		}

		if MyApp.App.Preferences().String("Mode") == "Desync" && logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
			if !logic.FileConflictCheck(MyApp) {
				MyApp.App.Preferences().SetString("Mode", "Sync")
				return
			} else {
				LoadSyncUI(MyApp)
				return
			}
		} else {
			ShowError("Cannot switch to Sync Mode, check server details or if server is running", MyApp)
		}

	})

	IpBind := binding.BindPreferenceString("IP", MyApp.App.Preferences())
	currentIPLbl := widget.NewLabelWithData(binding.NewSprintf("Current IP: %s", IpBind))
	PortBind := binding.BindPreferenceString("Port", MyApp.App.Preferences())
	currentPortLbl := widget.NewLabelWithData(binding.NewSprintf("Current Port: %s", PortBind))

	changeServerLbl := widget.NewLabel("Change Server Settings")
	centeredChangeServerLbl := container.NewCenter(changeServerLbl)
	changeIPEnt = widget.NewEntry()
	changeIPEnt.PlaceHolder = "New Server IP"
	changePortEnt = widget.NewEntry()
	changePortEnt.PlaceHolder = "New Server Port, Defaults To 7529"
	changeServerBtn := widget.NewButton("Change Server", func() {
		var err []string
		var port string

		ip := changeIPEnt.Text

		if ip == "" || net.ParseIP(ip) == nil {
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
			ShowError(strings.Join(err[:], "\n\n"), MyApp)
		} else if !logic.IsServerAccessible("http://" + ip + ":" + port) {
			ShowError("Server: "+ip+":"+port+" is inaccesslible", MyApp)
		} else {
			MyApp.App.Preferences().SetString("IP", ip)
			MyApp.App.Preferences().SetString("Port", port)
		}
	})

	return container.NewVBox(currentModeLbl, switchModeBtn, currentIPLbl, currentPortLbl, centeredChangeServerLbl, changeIPEnt, changePortEnt, changeServerBtn)
}