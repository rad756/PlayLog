package ui

import (
	"fmt"
	"net"
	"playlog/logic"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func MakeSettingsTab(MyApp logic.MyApp) fyne.CanvasObject {
	var changeIPEnt *widget.Entry
	var changePortEnt *widget.Entry
	ModeBind := binding.BindPreferenceString("StorageMode", MyApp.App.Preferences())

	currentModeLbl := widget.NewLabelWithData(binding.NewSprintf("Current mode: %s", ModeBind))

	switchModeBtn := widget.NewButton("Switch Mode", func() {
		var errStr []string
		if MyApp.App.Preferences().String("StorageMode") == "Sync" {
			MyApp.App.Preferences().SetString("StorageMode", "Local")
			return
		}

		if MyApp.App.Preferences().String("StorageMode") == "Local" && logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
			if !logic.FileConflictCheck(MyApp) {
				MyApp.App.Preferences().SetString("StorageMode", "Sync")
				return
			} else {
				LoadSyncUI(MyApp)
				return
			}
		} else {
			errStr = append(errStr, "Cannot connect to server, check details or if server is running")
		}

		if MyApp.App.Preferences().String("StorageMode") == "Desync" && logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
			if !logic.FileConflictCheck(MyApp) {
				MyApp.App.Preferences().SetString("StorageMode", "Sync")
				fmt.Println("a")
				return
			} else {
				fmt.Println("b")
				MyApp.Win.SetContent(LoadSyncUI(MyApp))
				return
			}
		} else {
			errStr = append(errStr, "Cannot switch to Sync Mode, check server details or if server is running")
		}

		if len(errStr) != 0 {
			dialog.NewError(logic.BuildError(errStr), MyApp.Win)
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
		var errStr []string
		var port string

		ip := changeIPEnt.Text

		if ip == "" || net.ParseIP(ip) == nil {
			errStr = append(errStr, "IP Empty or Invalid")
		}
		if _, err1 := strconv.Atoi(changePortEnt.Text); err1 != nil && changePortEnt.Text != "" {
			errStr = append(errStr, "Invalid Port")
		} else if _, err2 := strconv.Atoi(changePortEnt.Text); err2 == nil {
			port = changePortEnt.Text
		} else {
			port = "7529"
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
		} else if !logic.IsServerAccessible("http://" + ip + ":" + port) {
			dialog.ShowError(fmt.Errorf("Server: "+ip+":"+port+" is inaccesslible"), MyApp.Win)
		} else {
			MyApp.App.Preferences().SetString("IP", ip)
			MyApp.App.Preferences().SetString("Port", port)
		}
	})

	return container.NewVBox(currentModeLbl, switchModeBtn, currentIPLbl, currentPortLbl, centeredChangeServerLbl, changeIPEnt, changePortEnt, changeServerBtn)
}
