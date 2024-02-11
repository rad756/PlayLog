package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeShowTab() fyne.CanvasObject {
	var showWatchingLbl *widget.Label
	var showFinishedLbl *widget.Label
	showsList := readShowsList()
	selShowId := -1

	showLst := widget.NewList(
		func() int {
			return len(showsList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(showsList[i].name + " - S" + strconv.Itoa(showsList[i].season) + " E" + strconv.Itoa(showsList[i].episode) + " - " + showsList[i].finished)
		})

	showNameEnt := widget.NewEntry()
	showNameEnt.SetPlaceHolder("Enter Show Name")

	showSeasonEnt := widget.NewEntry()
	showSeasonEnt.PlaceHolder = "Enter Season #"
	moreSeasonBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if isNum(showSeasonEnt.Text) {
			showSeasonEnt.Text = moreFunc(showSeasonEnt.Text)
			showSeasonEnt.Refresh()
		} else if showSeasonEnt.Text == "" {
			showSeasonEnt.Text = "1"
			showSeasonEnt.Refresh()
		}
	})
	lessSeasonBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		if isNum(showSeasonEnt.Text) {
			showSeasonEnt.Text = lessFunc(showSeasonEnt.Text)
			showSeasonEnt.Refresh()
		} else if showSeasonEnt.Text == "" {
			showSeasonEnt.Text = "0"
			showSeasonEnt.Refresh()
		}
	})

	showEpisodeEnt := widget.NewEntry()
	showEpisodeEnt.PlaceHolder = "Enter Episode #"
	moreEpisodeBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if isNum(showEpisodeEnt.Text) {
			showEpisodeEnt.Text = moreFunc(showEpisodeEnt.Text)
			showEpisodeEnt.Refresh()
		} else if showEpisodeEnt.Text == "" {
			showEpisodeEnt.Text = "1"
			showEpisodeEnt.Refresh()
		}
	})
	lessEpisodeBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		if isNum(showEpisodeEnt.Text) {
			showEpisodeEnt.Text = lessFunc(showEpisodeEnt.Text)
			showEpisodeEnt.Refresh()
		} else if showEpisodeEnt.Text == "" {
			showEpisodeEnt.Text = "0"
			showEpisodeEnt.Refresh()
		}
	})

	showFinishedCck := widget.NewCheck("Finished Watching ", func(value bool) {})
	centeredShowFinishedCck := container.New(layout.NewCenterLayout(), showFinishedCck)

	showAddBtn := widget.NewButton("Add Show", func() {

		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if showNameEnt.Text != "" && isNum(showSeasonEnt.Text) && isNum(showEpisodeEnt.Text) && noComma(showNameEnt.Text) {
			seasonNum, _ := strconv.Atoi(showSeasonEnt.Text)
			episodeNum, _ := strconv.Atoi(showEpisodeEnt.Text)

			showsList = addShowFunc(showNameEnt.Text, seasonNum, episodeNum, showFinishedCck.Checked, showsList)
			saveShow(showsList)
			showWatchingLbl.Text = "Currently watching " + countWatching(showsList) + " shows"
			showWatchingLbl.Refresh()
			showFinishedLbl.Text = "Finished watching " + countWatched(showsList) + " shows"
			showFinishedLbl.Refresh()
			selShowId = -1
			showLst.UnselectAll()
			showLst.Refresh()
		}

	})
	showChangeBtn := widget.NewButton("Change Selected Show", func() {
		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if showNameEnt.Text != "" && isNum(showSeasonEnt.Text) && isNum(showEpisodeEnt.Text) && selShowId != -1 && noComma(showNameEnt.Text) {
			seasonNum, _ := strconv.Atoi(showSeasonEnt.Text)
			episodeNum, _ := strconv.Atoi(showEpisodeEnt.Text)

			showsList = deleteShowFunc(selShowId, showsList)
			showsList = addShowFunc(showNameEnt.Text, seasonNum, episodeNum, showFinishedCck.Checked, showsList)
			saveShow(showsList)
			showWatchingLbl.Text = "Currently watching " + countWatching(showsList) + " shows"
			showWatchingLbl.Refresh()
			showFinishedLbl.Text = "Finished watching " + countWatched(showsList) + " shows"
			showFinishedLbl.Refresh()
			selShowId = -1
			showLst.UnselectAll()
			showLst.Refresh()
		}
	})

	showWatchingLbl = widget.NewLabel("Currently watching " + countWatching(showsList) + " shows")
	centeredShowWatchinglbl := container.New(layout.NewCenterLayout(), showWatchingLbl)
	showFinishedLbl = widget.NewLabel("Finished watching " + countWatched(showsList) + " shows")
	centeredShowFinishedlbl := container.New(layout.NewCenterLayout(), showFinishedLbl)

	showDeleteBtn := widget.NewButton("Delete Selected Show", func() {
		if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else if selShowId != -1 {
			showsList = deleteShowFunc(selShowId, showsList)

			showWatchingLbl.Text = "Currently watching " + countWatching(showsList) + " shows"
			showWatchingLbl.Refresh()
			showFinishedLbl.Text = "Finished watching " + countWatched(showsList) + " shows"
			showFinishedLbl.Refresh()
			selShowId = -1
			showLst.UnselectAll()
			showLst.Refresh()
		}
	})

	showLst.OnSelected = func(id widget.ListItemID) {
		selShowId = id

		showNameEnt.Text = showsList[id].name
		showNameEnt.Refresh()

		showSeasonEnt.Text = strconv.Itoa(showsList[id].season)
		showSeasonEnt.Refresh()

		showEpisodeEnt.Text = strconv.Itoa(showsList[id].episode)
		showEpisodeEnt.Refresh()

		if showsList[id].finished == "Yes" {
			showFinishedCck.SetChecked(true)
		} else {
			showFinishedCck.SetChecked(false)
		}
		showFinishedCck.Refresh()
	}

	seasonBox := container.New(layout.NewBorderLayout(nil, nil, lessSeasonBtn, moreSeasonBtn), lessSeasonBtn, showSeasonEnt, moreSeasonBtn)
	episodeBox := container.New(layout.NewBorderLayout(nil, nil, lessEpisodeBtn, moreEpisodeBtn), lessEpisodeBtn, showEpisodeEnt, moreEpisodeBtn)

	showBox := container.New(layout.NewVBoxLayout(), showNameEnt, seasonBox, episodeBox, centeredShowFinishedCck, showAddBtn, layout.NewSpacer(), showChangeBtn, layout.NewSpacer(), centeredShowWatchinglbl, centeredShowFinishedlbl, layout.NewSpacer(), showDeleteBtn)

	showTab := container.NewHSplit(showLst, showBox)
	showTab.Offset = splitOffset

	return showTab
}
