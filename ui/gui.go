package ui

import (
	"context"
	"fmt"
	"net"
	"playlog/logic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoadGUI(MyApp *logic.MyApp) {
	if MyApp.App.Preferences().Bool("FirstRun") {
		LoadSetupUI(MyApp)
		return
	}

	if len(MyApp.Tabs) == 0 {
		LoadCreateTabUI(MyApp)
		return
	}

	if MyApp.App.Preferences().String("StorageMode") == "Local" {
		LoadMainUI(MyApp)
		return
	}

	BootSyncingUI(MyApp)
}

func BootSyncingUI(MyApp *logic.MyApp) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go logic.IsServerAccessibleBoot(MyApp, ctx, cancel, LoadMenuAfterServerBootCheck)

	checkLbl := widget.NewLabel("Checking if server is accessible...")
	centeredCheckLbl := container.NewCenter(checkLbl)
	progessBar := widget.NewProgressBarInfinite()

	desyncModeBtn := widget.NewButton("Cancel Check & Enter Desync Mode", func() {
		cancel()
	})

	vbox := container.NewVBox(layout.NewSpacer(), centeredCheckLbl, progessBar, layout.NewSpacer())

	content := container.NewBorder(nil, desyncModeBtn, nil, nil, vbox)

	MyApp.Win.SetContent(content)
}

func MakeTab(id int, MyApp *logic.MyApp) *container.TabItem {
	var tab *container.TabItem
	if MyApp.Tabs[id].Mode == "Alpha" {
		name := MyApp.Tabs[id].Name
		kind := MyApp.Tabs[id].Kind
		main := TabAlpha{Name: name, Kind: kind, ID: -1}
		second := logic.ReadAlphaKind(name+"-"+kind, MyApp)

		tab = container.NewTabItem(name+"s", NewTabAlpha(logic.ReadAlphaSlice(main.Name, MyApp), MyApp, main, second))
	} else if MyApp.Tabs[id].Mode == "Beta" {
		name := MyApp.Tabs[id].Name
		count := MyApp.Tabs[id].Count
		subCount := MyApp.Tabs[id].SubCount
		main := TabBeta{Name: name, Count: count, SubCount: subCount, ID: -1}

		tab = container.NewTabItem(name+"s", NewTabBeta(logic.ReadBetaSlice(main.Name, MyApp), MyApp, main))
	}

	return tab
}

func LoadMainUI(MyApp *logic.MyApp) {
	var tabs []*container.TabItem

	for i := 0; i < len(MyApp.Tabs); i++ {
		tabs = append(tabs, MakeTab(i, MyApp))
	}
	tabs = append(tabs, container.NewTabItem("Settings", MakeSettingsTab(MyApp)))

	content := container.NewAppTabs()
	for i := 0; i < len(tabs); i++ {
		content.Append(tabs[i])
	}

	MyApp.Win.SetContent(content)
}

func LoadStartUpServerError(MyApp *logic.MyApp) {
	topLbl := widget.NewLabel("-- Startup Error --")
	topContent := container.New(layout.NewCenterLayout(), topLbl)
	errorLbl := widget.NewLabel("Server with IP " + MyApp.App.Preferences().String("IP") + " is inaccessible\nThe app will start in Desync Mode and try to sync upon next startup\nOr you can try to enter Sync Mode by pressing Switch Mode in Settings")
	desyncModeBtn := widget.NewButton("Enter Desync Mode", func() {
		MyApp.App.Preferences().SetString("StorageMode", "Desync")

		LoadMainUI(MyApp)
	})

	content := container.New(layout.NewVBoxLayout(), topContent, errorLbl, desyncModeBtn)

	MyApp.Win.SetContent(content)
}

func LoadSyncUI(MyApp *logic.MyApp) {
	errorLbl := widget.NewLabel("-- Sync Error --")
	errorCentered := container.NewCenter(errorLbl)
	questionLbl := widget.NewLabel("Do you want to download server files OR upload local files to server?")
	questionCentered := container.NewCenter(questionLbl)

	serverFilesBtn := widget.NewButton("Download Server Files", func() {
		logic.DownloadFromServer(MyApp)
		MyApp.App.Preferences().SetString("StorageMode", "Sync")
		for {
			if MyApp.App.Preferences().String("StorageMode") != "Sync" {
				MyApp.App.Preferences().SetString("StorageMode", "Sync")
				fmt.Println("BUG: StorageMode was not set")
			} else {
				break
			}
		}
		LoadMainUI(MyApp)
	})
	orLbl := widget.NewLabel("OR")
	orCentered := container.NewCenter(orLbl)
	localFilesBtn := widget.NewButton("Upload Local Files", func() {
		logic.UploadToServer(MyApp)
		MyApp.App.Preferences().SetString("StorageMode", "Sync")
		for {
			if MyApp.App.Preferences().String("StorageMode") != "Sync" {
				MyApp.App.Preferences().SetString("StorageMode", "Sync")
				fmt.Println("BUG: StorageMode was not set")
			} else {
				break
			}
		}
		LoadMainUI(MyApp)
	})

	content := container.NewVBox(errorCentered, questionCentered, serverFilesBtn, orCentered, localFilesBtn)

	MyApp.Win.SetContent(content)
}

func ShowServerInaccessibleError(MyApp *logic.MyApp) {
	var errorPpu *widget.PopUp

	errorLbl := widget.NewLabel(fmt.Sprintf("Server with IP: %s:%s is inaccessible", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port")))
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
			dialog.ShowError(fmt.Errorf("IP empty"), MyApp.Win)
		} else if net.ParseIP(changeIPEnt.Text) == nil {
			dialog.ShowError(fmt.Errorf(changeIPEnt.Text+" is not valid IP"), MyApp.Win)
		} else if !logic.IsServerAccessible("http://" + changeIPEnt.Text + ":" + changePortEnt.Text) {
			dialog.ShowError(fmt.Errorf("Server with details: "+changeIPEnt.Text+":"+port+" is inaccessible"), MyApp.Win)
		} else {
			MyApp.App.Preferences().SetString("StorageMode", "Sync")
			MyApp.App.Preferences().SetString("IP", changeIPEnt.Text)
			MyApp.App.Preferences().SetString("Port", port)
			errorPpu.Hide()
		}
	})

	orLbl := widget.NewLabel("OR")
	centeredOrLbl := container.NewCenter(orLbl)

	backBtn := widget.NewButton("Enter Desync Mode", func() {
		MyApp.App.Preferences().SetString("StorageMode", "Desync")

		errorPpu.Hide()
	})
	backLbl := widget.NewLabel("Current change will not be pushed.\nTry it again after swithching to Server Down Mode")

	content := container.New(layout.NewVBoxLayout(), errorLbl, changeLbl, changeIPEnt, changePortEnt, changeServerBtn, centeredOrLbl, backBtn, backLbl)

	errorPpu = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	errorPpu.Show()
}

func LoadMenuAfterServerBootCheck(MyApp *logic.MyApp, err error) {
	if err != nil {
		MyApp.App.Preferences().SetString("StorageMode", "Desync")
		dialog.ShowError(fmt.Errorf("Cannot connect to server, entered Desync mode!"), MyApp.Win)
		LoadMainUI(MyApp)
	} else if logic.FileConflictCheck(MyApp) {
		LoadSyncUI(MyApp)
	} else {
		MyApp.App.Preferences().SetString("StorageMode", "Sync")
		LoadMainUI(MyApp)
	}
}

func GetLoadingPopUp(MyApp *logic.MyApp, cancel context.CancelFunc) *widget.PopUp {
	var popup *widget.PopUp
	lbl := widget.NewLabel("Checking server...")
	centeredLbl := container.NewCenter(lbl)
	pg := widget.NewProgressBarInfinite()
	btn := widget.NewButton("Cancel Check", func() {
		cancel()
		popup.Hide()
	})

	content := container.NewBorder(centeredLbl, btn, nil, nil, pg)
	popup = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	popup.Resize(fyne.NewSize(200, 0))

	return popup
}

func GetLoadingPopUpGR(MyApp *logic.MyApp, cancel context.CancelFunc) *widget.PopUp {
	var popup *widget.PopUp
	lbl := widget.NewLabel("Checking server...")
	centeredLbl := container.NewCenter(lbl)
	pg := widget.NewProgressBarInfinite()
	lbl2 := widget.NewLabel("If you cancel, change will be made locally")
	centeredLbl2 := container.NewCenter(lbl2)
	btn := widget.NewButton("Cancel Check and Enter Desync Mode", func() {
		cancel()
		MyApp.App.Preferences().SetString("StorageMode", "Desync")
		popup.Hide()
	})

	vbox := container.NewVBox(centeredLbl2, btn)

	content := container.NewBorder(centeredLbl, vbox, nil, nil, pg)
	popup = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	popup.Resize(fyne.NewSize(200, 0))

	return popup
}

func GetServerCheckPopUp(MyApp *logic.MyApp, cancel context.CancelFunc) *widget.PopUp {
	var popup *widget.PopUp
	checkLbl := widget.NewLabel("Checking if server is accessible...")
	centeredChecklbl := container.NewCenter(checkLbl)
	progressBar := widget.NewProgressBarInfinite()

	cancelBtn := widget.NewButton("Cancel Check & Enter Desync Mode", func() {
		cancel()
		popup.Hide()
	})
	warningLbl := widget.NewLabel("If canceled, you will have to retry previous action")
	centeredWarningLbl := container.NewCenter(warningLbl)

	content := container.NewVBox(layout.NewSpacer(), centeredChecklbl, progressBar, layout.NewSpacer(), cancelBtn, centeredWarningLbl)

	popup = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	popup.Resize(fyne.NewSize(200, 0))
	return popup
}

func CheckServer(MyApp *logic.MyApp) {
	if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		popup := GetLoadingPopUpGR(MyApp, cancel)

		go logic.IsServerAccessibleGR(MyApp, ctx, cancel, popup)
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			//do not load popup after delay
		default:
			popup.Show()
		}

		//If popup is not hidden, this will stop code execution until popup is hidden
		if !popup.Hidden {
			select {
			case <-ctx.Done():
			}
		}
	}
}
