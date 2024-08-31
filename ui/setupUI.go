package ui

import (
	"errors"
	"fmt"
	"net"
	"playlog/logic"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func LoadSetupUI(MyApp *logic.MyApp) {
	questionLbl := widget.NewLabel("Select which mode app will run in, Local or Server Sync")
	localBtn := widget.NewButton("Local Setup", func() {
		MyApp.App.Preferences().SetString("StorageMode", "Local")
		MyApp.App.Preferences().SetBool("FirstRun", false)
		LoadCreateTabUI(MyApp)
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
			//Make if statement to check if tabs are created, if true go to main menu, if false go to create tab menu
			LoadMainUI(MyApp)
		}

	})

	content := container.New(layout.NewVBoxLayout(), questionLbl, localBtn, serverIpLbl, serverIpEnt, serverPortLbl, serverPortEnt, serverBtn)

	MyApp.Win.SetContent(content)
}

func LoadCreateTabUI(MyApp *logic.MyApp) {
	var tabLst *widget.List
	listID := -1

	upBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		if listID > 0 {
			logic.MoveTabUp(listID, MyApp)
			listID = listID - 1
			tabLst.Select(listID)
			tabLst.Refresh()
			logic.CreateTabsFile(MyApp)
		}
	})

	tabLst = widget.NewList(
		func() int {
			return len(MyApp.Tabs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			var text string
			if MyApp.Tabs[i].Mode == "Alpha" {
				text = MyApp.Tabs[i].Name + " - Alpha"
			} else if MyApp.Tabs[i].Mode == "Beta" {
				text = MyApp.Tabs[i].Name + " - Beta"
			}
			o.(*widget.Label).SetText(text)
		})

	downBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		if listID > -1 && listID < len(MyApp.Tabs)-1 {
			logic.MoveTabDown(listID, MyApp)
			listID = listID + 1
			tabLst.Select(listID)
			tabLst.Refresh()
			logic.CreateTabsFile(MyApp)
		}
	})
	addAlphaBtn := widget.NewButton("Add Type A Tab", func() {
		LoadCreateAlphaTabUI(tabLst, MyApp)
		listID = -1
		tabLst.UnselectAll()
	})
	addBetaBtn := widget.NewButton("Add Type B Tab", func() {
		LoadCreateBetaTabUI(tabLst, MyApp)
		listID = -1
		tabLst.UnselectAll()
	})
	deleteBtb := widget.NewButton("Delete Selected", func() {
		if listID == -1 {
			dialog.ShowError(errors.New("No tab was selected to be deleted!"), MyApp.Win)
			return
		}

		logic.DeleteTab(listID, MyApp)
		tabLst.UnselectAll()
		tabLst.Refresh()
		listID = -1
	})
	mainMenuBtn := widget.NewButton("Enter Main Menu", func() {
		LoadMainUI(MyApp)
	})

	tabLst.OnSelected = func(id widget.ListItemID) {
		listID = id
	}

	vbox := container.NewVBox(downBtn, layout.NewSpacer(), addAlphaBtn, addBetaBtn, layout.NewSpacer(), deleteBtb, layout.NewSpacer(), mainMenuBtn)

	content := container.NewBorder(upBtn, vbox, nil, nil, tabLst) //have this inside scroll

	MyApp.Win.SetContent(content)
}

func LoadCreateAlphaTabUI(lst *widget.List, MyApp *logic.MyApp) {
	var popup *widget.PopUp

	titleLbl := widget.NewLabel("Type A Tab Creator")
	nameEnt := widget.NewEntry()
	nameEnt.PlaceHolder = "Name of Tab Items"
	kindEnt := widget.NewEntry()
	kindEnt.PlaceHolder = "Kind/Type of Tab Items"
	addBtn := widget.NewButton("Add Tab", func() {
		var errStr []string
		var name = nameEnt.Text
		var kind = kindEnt.Text

		if name == "" {
			errStr = append(errStr, "Tab name cannot be empty")
		}
		if kind == "" {
			errStr = append(errStr, "Tab kind/type cannot be empty")
		}

		if logic.ContainsComma(kind) {
			errStr = append(errStr, "Tab kind/type cannot contain a comma")
		}

		if logic.PathExists(name+".json", MyApp) {
			errStr = append(errStr, name+" already exists")
		}

		if logic.PathExists(name+"-"+kind+".json", MyApp) {
			errStr = append(errStr, name+"-"+kind+" already exists")
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
			return
		}

		MyApp.Tabs = append(MyApp.Tabs, logic.Tab{Mode: "Alpha", Name: name, Kind: kind})
		logic.CreateTabsFile(MyApp)
		logic.SaveAlphaSlice(name, MyApp, logic.AlphaSlice{})
		logic.SaveAlphaKind(name+"-"+kind, MyApp, logic.Kind{})
		lst.Refresh()
	})
	exitBtn := widget.NewButton("Exit", func() { popup.Hide() })

	content := container.NewVBox(titleLbl, nameEnt, kindEnt, addBtn, exitBtn)

	popup = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	popup.Resize(fyne.NewSize(200, 0))
	popup.Show()
}

func LoadCreateBetaTabUI(lst *widget.List, MyApp *logic.MyApp) {
	var popup *widget.PopUp

	titleLbl := widget.NewLabel("Type B Tab Creator")
	nameEnt := widget.NewEntry()
	nameEnt.PlaceHolder = "Name of Tab Items"
	countEnt := widget.NewEntry()
	countEnt.PlaceHolder = "Count of Tab Items"
	subCountEnt := widget.NewEntry()
	subCountEnt.PlaceHolder = "Sub-Count of Tab Items"
	addBtn := widget.NewButton("Add Tab", func() {
		var errStr []string
		var name = nameEnt.Text
		var count = countEnt.Text
		var subCount = subCountEnt.Text

		if name == "" {
			errStr = append(errStr, "Tab name cannot be empty")
		}

		if count == "" {
			errStr = append(errStr, "Count cannot be empty")
		}

		if subCount == "" {
			errStr = append(errStr, "Sub-Count cannot be empty")
		}

		if logic.PathExists(name+".json", MyApp) {
			errStr = append(errStr, name+" already exists")
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
			return
		}

		MyApp.Tabs = append(MyApp.Tabs, logic.Tab{Mode: "Beta", Name: name, Count: count, SubCount: subCount})
		logic.CreateTabsFile(MyApp)
		logic.SaveBetaSlice(name, MyApp, logic.BetaSlice{})
		lst.Refresh()
	})
	exitBtn := widget.NewButton("Exit", func() { popup.Hide() })

	content := container.NewVBox(titleLbl, nameEnt, countEnt, subCountEnt, addBtn, exitBtn)

	popup = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	popup.Resize(fyne.NewSize(200, 0))
	popup.Show()
}
