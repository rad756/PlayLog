package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeGameTab() fyne.CanvasObject {
	var gameFinishedLbl *widget.Label
	gamesList := readGamesList()
	platformList := readPlatformList()
	selGameId := -1

	gameLst := widget.NewList(
		func() int {
			return len(gamesList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(gamesList[i].name + " - " + gamesList[i].platform)
		})

	gameNameEnt := widget.NewEntry()
	gameNameEnt.SetPlaceHolder("Enter Game Name")

	gamePlatformDdl := widget.NewSelect(platformList, nil)
	gamePlatformDdl.PlaceHolder = "Select Platform"

	gameAddBtn := widget.NewButton("Add Game", func() {
		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if gameNameEnt.Text != "" && gamePlatformDdl.Selected != "" && noComma(gameNameEnt.Text) {
			gamesList = addGameFunc(gameNameEnt.Text, gamePlatformDdl.Selected, gamesList)
			saveGame(gamesList)
			gameFinishedLbl.SetText(strconv.Itoa(len(gamesList)) + " Games Finished")
			selGameId = -1
			gameLst.UnselectAll()
			gameLst.Refresh()
		}

	})

	gameChangeBtn := widget.NewButton("Change Selected Game", func() {
		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if gameNameEnt.Text != "" && gamePlatformDdl.Selected != "" && selGameId != -1 && noComma(gameNameEnt.Text) {
			gamesList = deleteGameFunc(selGameId, gamesList)
			gamesList = addGameFunc(gameNameEnt.Text, gamePlatformDdl.Selected, gamesList)
			saveGame(gamesList)
			selGameId = -1
			gameLst.UnselectAll()
			gameLst.Refresh()
		}
	})

	gameChangePlatformBtn := widget.NewButton("Change Platforms", func() {
		var platformPpu *widget.PopUp
		var platformDdl *widget.Select

		platformEnt := widget.NewEntry()
		platformEnt.SetPlaceHolder("Enter name of platform")
		platformAddBtn := widget.NewButton("Add Platform", func() {
			if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
				showServerInaccessibleError()
			} else if platformEnt.Text != "" && noComma(platformEnt.Text) {
				platformList = addPlatformFunc(platformEnt.Text, platformList)
				savePlatform(platformList)
				gamePlatformDdl.Options = platformList
				platformDdl.Options = platformList
			}

		})

		platformDdl = widget.NewSelect(platformList, nil)
		platformDeleteBtn := widget.NewButton("Delete selected platform", func() {
			if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
				showServerInaccessibleError()
			} else {
				platformList = deletePlatformFunc(platformDdl.SelectedIndex(), platformList)
				savePlatform(platformList)
				platformDdl.Options = platformList
				gamePlatformDdl.Options = platformList
				platformDdl.ClearSelected()
			}
		})

		exitBtn := widget.NewButton("Exit", func() {
			platformPpu.Hide()
		})
		content := container.New(layout.NewVBoxLayout(), platformEnt, platformAddBtn, platformDdl, platformDeleteBtn, exitBtn)
		platformPpu = widget.NewModalPopUp(content, mainWin.Canvas())
		platformPpu.Resize(fyne.NewSize(200, 0))
		platformPpu.Show()
	})

	gameFinishedLbl = widget.NewLabel(strconv.Itoa(len(gamesList)) + " Games Finished")
	centeredGameFinishedLbl := container.New(layout.NewCenterLayout(), gameFinishedLbl)

	gameDeleteBtn := widget.NewButton("Delete Selected Game", func() {
		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if selGameId != -1 {
			gamesList = deleteGameFunc(selGameId, gamesList)
			saveGame(gamesList)
			selGameId = -1
			gameLst.UnselectAll()
			gameLst.Refresh()
			gameFinishedLbl.SetText(strconv.Itoa(len(gamesList)) + " Games Finished")
		}
	})

	gameLst.OnSelected = func(id widget.ListItemID) {
		selGameId = id

		//Copies selected game name to entry
		gameNameEnt.Text = gamesList[id].name
		gameNameEnt.Refresh()

		//Sets selection in platform dropdown
		for i := range platformList {
			if platformList[i] == gamesList[id].platform {
				gamePlatformDdl.SetSelectedIndex(i)
			}
		}
	}

	gameBox := container.New(layout.NewVBoxLayout(), gameNameEnt, gamePlatformDdl, gameAddBtn, layout.NewSpacer(), gameChangeBtn, gameChangePlatformBtn, layout.NewSpacer(), centeredGameFinishedLbl, layout.NewSpacer(), gameDeleteBtn)

	gameTab := container.NewHSplit(gameLst, gameBox)
	gameTab.Offset = splitOffset
	return gameTab
}
