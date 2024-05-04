package ui

import (
	"context"
	"fmt"
	"net"
	"playlog/logic"
	"strconv"
	"time"

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
		if MyApp.App.Preferences().String("StorageMode") == "Sync" {
			MyApp.App.Preferences().SetString("StorageMode", "Local")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		popup := GetLoadingPopUp(MyApp, cancel)

		go logic.IsServerAccessibleSwitch(MyApp, ctx, cancel, popup, LoadSyncUI)
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			//do not load popup after delay
		default:
			popup.Show()
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
