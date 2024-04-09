package ui

import (
	"fmt"
	"net"
	"playlog/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func LoadGUI(MyApp logic.MyApp) fyne.CanvasObject {
	var content fyne.CanvasObject

	if MyApp.App.Preferences().Bool("FirstRun") {
		content = LoadSetupUI(MyApp)
	} else if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		if !logic.IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
			content = LoadStartUpServerError(MyApp)
		} else if logic.FileConflictCheck(MyApp) {
			content = LoadSyncUI(MyApp)
		} else {
			content = LoadMainUI(MyApp)
		}
	} else { //Runs in Local Mode
		content = LoadMainUI(MyApp)
	}

	return content
}

func LoadMainUI(MyApp logic.MyApp) *container.AppTabs {
	gameTab := &TabAlpha{Name: "Game", Kind: "Platform", ID: -1}
	gameKind := logic.ReadAlphaKind("Game-Platform", MyApp)
	movieTab := &TabAlpha{Name: "Movie", Kind: "Genre", ID: -1}
	movieKind := logic.ReadAlphaKind("Movie-Genre", MyApp)
	showTab := &TabBeta{Name: "Show", Count: "Season", SubCount: "Episode", Action: "Watch", ID: -1}

	content := container.NewAppTabs(
		container.NewTabItem("Games", NewTabAlpha(logic.ReadAlphaSlice(gameTab.Name, MyApp), MyApp, *gameTab, gameKind)),
		container.NewTabItem("Movies", NewTabAlpha(logic.ReadAlphaSlice(movieTab.Name, MyApp), MyApp, *movieTab, movieKind)),
		container.NewTabItem("Shows", NewTabBeta(logic.ReadBetaSlice(showTab.Name, MyApp), MyApp, *showTab)),
		container.NewTabItem("Settings", MakeSettingsTab(MyApp)))

	return content
}

func LoadStartUpServerError(MyApp logic.MyApp) fyne.CanvasObject {
	topLbl := widget.NewLabel("-- Startup Error --")
	topContent := container.New(layout.NewCenterLayout(), topLbl)
	errorLbl := widget.NewLabel("Server with IP " + MyApp.App.Preferences().String("IP") + " is inaccessible\nThe app will start in Desync Mode and try to sync upon next startup\nOr you can try to enter Sync Mode by pressing Switch Mode in Settings")
	desyncModeBtn := widget.NewButton("Enter Desync Mode", func() {
		MyApp.App.Preferences().SetString("StorageMode", "Desync")

		MyApp.Win.SetContent(LoadMainUI(MyApp))
	})

	return container.New(layout.NewVBoxLayout(), topContent, errorLbl, desyncModeBtn)
}

func LoadSyncUI(MyApp logic.MyApp) fyne.CanvasObject {
	errorLbl := widget.NewLabel("-- Sync Error --")
	errorCentered := container.NewCenter(errorLbl)
	questionLbl := widget.NewLabel("Do you want to download server files OR upload local files to server?")
	questionCentered := container.NewCenter(questionLbl)

	serverFilesBtn := widget.NewButton("Download Server Files", func() {
		logic.DownloadFromServer(MyApp)
		MyApp.App.Preferences().SetString("StorageMode", "Sync")
		MyApp.Win.SetContent(LoadMainUI(MyApp))
	})
	orLbl := widget.NewLabel("OR")
	orCentered := container.NewCenter(orLbl)
	localFilesBtn := widget.NewButton("Upload Local Files", func() {
		logic.UploadToServer(MyApp)
		MyApp.App.Preferences().SetString("StorageMode", "Sync")
		MyApp.Win.SetContent(LoadMainUI(MyApp))
	})

	return container.NewVBox(errorCentered, questionCentered, serverFilesBtn, orCentered, localFilesBtn)

}

func ShowError(errorText string, MyApp logic.MyApp) {
	var errorPpu *widget.PopUp

	topLbl := widget.NewLabel("-- Error(s) --")
	topContent := container.New(layout.NewCenterLayout(), topLbl)
	errorLbl := widget.NewLabel(errorText)
	backBtn := widget.NewButton("OK", func() { errorPpu.Hide() })

	content := container.New(layout.NewVBoxLayout(), topContent, errorLbl, backBtn)

	errorPpu = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	errorPpu.Show()
}

func ShowServerInaccessibleError(MyApp logic.MyApp) {
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
			ShowError("IP empty", MyApp)
		} else if net.ParseIP(changeIPEnt.Text) == nil {
			ShowError(changeIPEnt.Text+" is not valid IP", MyApp)
		} else if !logic.IsServerAccessible("http://" + changeIPEnt.Text + ":" + changePortEnt.Text) {
			ShowError("Server with details: "+changeIPEnt.Text+":"+port+" is inaccessible", MyApp)
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
