package ui

import (
	"fmt"
	"playlog/logic"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TabBeta struct {
	Name     string
	Count    string
	SubCount string
	ID       int
}

func NewTabBeta(betaSlice *logic.BetaSlice, MyApp *logic.MyApp, tabBeta TabBeta) fyne.CanvasObject {
	var currentLbl *widget.Label
	var finishedLbl *widget.Label
	tabBeta.ID = -1

	lst := widget.NewList(
		func() int {
			return len(betaSlice.Slice)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			var fin string

			if betaSlice.Slice[i].Finished {
				fin = "Yes"
			} else {
				fin = "No"
			}

			o.(*widget.Label).SetText(fmt.Sprintf("%s - S%d - E%d - %s", betaSlice.Slice[i].Name, betaSlice.Slice[i].Count, betaSlice.Slice[i].SubCount, fin))
		})

	nameEnt := widget.NewEntry()
	nameEnt.SetPlaceHolder(fmt.Sprintf("Enter %s Name", tabBeta.Name))

	countEnt := widget.NewEntry()
	countEnt.PlaceHolder = fmt.Sprintf("Enter %s #", tabBeta.Count)
	moreCountBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		countEnt.Text = logic.MoreBeta(countEnt.Text)
		countEnt.Refresh()
	})
	lessCountBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		countEnt.Text = logic.LessBeta(countEnt.Text)
		countEnt.Refresh()
	})

	subCountEnt := widget.NewEntry()
	subCountEnt.PlaceHolder = fmt.Sprintf("Enter %s #", tabBeta.SubCount)
	moreSubCountBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		subCountEnt.Text = logic.MoreBeta(subCountEnt.Text)
		subCountEnt.Refresh()
	})
	lessSubCountBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		subCountEnt.Text = logic.LessBeta(subCountEnt.Text)
		subCountEnt.Refresh()
	})

	finishedCck := widget.NewCheck("Finished", func(value bool) {})
	centeredFinishedCck := container.NewCenter(finishedCck)

	addBtn := widget.NewButton("Add "+tabBeta.Name, func() {
		var errStr []string

		if nameEnt.Text == "" {
			errStr = append(errStr, fmt.Sprintf("%s name empty", tabBeta.Name))
		}
		if !logic.IsNum(countEnt.Text) {
			errStr = append(errStr, fmt.Sprintf("%s %s must contain a number", tabBeta.Name, strings.ToLower(tabBeta.Count)))
		}
		if !logic.IsNum(subCountEnt.Text) {
			errStr = append(errStr, fmt.Sprintf("%s %s must contain a number", tabBeta.Name, strings.ToLower(tabBeta.SubCount)))
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
		}

		CheckServer(MyApp)

		betaSlice.AddBeta(nameEnt.Text, countEnt.Text, subCountEnt.Text, finishedCck.Checked, MyApp, tabBeta.Name)
		tabBeta.ID = -1
		lst.UnselectAll()
		lst.Refresh()
		currentLbl.SetText(fmt.Sprintf("Current %ss: %d", tabBeta.Name, betaSlice.CountCurrent()))
		finishedLbl.SetText(fmt.Sprintf("Finished %ss: %d", tabBeta.Name, betaSlice.CountFinished()))
	})

	changeBtn := widget.NewButton("Change Selected "+tabBeta.Name, func() {
		var errStr []string

		if nameEnt.Text == "" {
			errStr = append(errStr, fmt.Sprintf("%s name empty", tabBeta.Name))
		}
		if !logic.IsNum(countEnt.Text) {
			errStr = append(errStr, fmt.Sprintf("%s %s must contain a number", tabBeta.Name, strings.ToLower(tabBeta.Count)))
		}
		if !logic.IsNum(subCountEnt.Text) {
			errStr = append(errStr, fmt.Sprintf("%s %s must contain a number", tabBeta.Name, strings.ToLower(tabBeta.SubCount)))
		}
		if tabBeta.ID == -1 {
			errStr = append(errStr, fmt.Sprintf("No %s was selected to be changed", strings.ToLower(tabBeta.Name)))
		}

		if len(errStr) != 0 {
			dialog.ShowError(logic.BuildError(errStr), MyApp.Win)
		}

		CheckServer(MyApp)

		betaSlice.DeleteBeta(tabBeta.ID, MyApp, tabBeta.Name)
		betaSlice.AddBeta(nameEnt.Text, countEnt.Text, subCountEnt.Text, finishedCck.Checked, MyApp, tabBeta.Name)
		lst.UnselectAll()
		lst.Refresh()
		currentLbl.SetText(fmt.Sprintf("Current %ss: %d", tabBeta.Name, betaSlice.CountCurrent()))
		finishedLbl.SetText(fmt.Sprintf("Finished %ss: %d", tabBeta.Name, betaSlice.CountFinished()))
	})

	currentLbl = widget.NewLabel(fmt.Sprintf("Current %ss: %d", tabBeta.Name, betaSlice.CountCurrent()))
	centeredCurrentLbl := container.NewCenter(currentLbl)
	finishedLbl = widget.NewLabel(fmt.Sprintf("Finished %ss: %d", tabBeta.Name, betaSlice.CountFinished()))
	centeredFinishedLbl := container.NewCenter(finishedLbl)

	deleteBtn := widget.NewButton("Delete Selected "+tabBeta.Name, func() {
		if tabBeta.ID == -1 {
			dialog.ShowError(fmt.Errorf("No %s was selected to be deleted", strings.ToLower(tabBeta.Name)), MyApp.Win)
		}

		CheckServer(MyApp)

		betaSlice.DeleteBeta(tabBeta.ID, MyApp, tabBeta.Name)
		tabBeta.ID = -1
		lst.UnselectAll()
		lst.Refresh()
		currentLbl.SetText(fmt.Sprintf("Current %ss: %d", tabBeta.Name, betaSlice.CountCurrent()))
		finishedLbl.SetText(fmt.Sprintf("Finished %ss: %d", tabBeta.Name, betaSlice.CountFinished()))
	})

	lst.OnSelected = func(id widget.ListItemID) {
		tabBeta.ID = id

		nameEnt.Text = betaSlice.Slice[id].Name
		nameEnt.Refresh()

		countEnt.Text = strconv.Itoa(betaSlice.Slice[id].Count)
		countEnt.Refresh()

		subCountEnt.Text = strconv.Itoa(betaSlice.Slice[id].SubCount)
		subCountEnt.Refresh()

		if betaSlice.Slice[id].Finished {
			finishedCck.SetChecked(true)
		} else {
			finishedCck.SetChecked(false)
		}
	}

	countBorder := container.NewBorder(nil, nil, lessCountBtn, moreCountBtn, countEnt)
	subCountBorder := container.NewBorder(nil, nil, lessSubCountBtn, moreSubCountBtn, subCountEnt)

	if MyApp.Mobile {
		mobileChangeBtn := widget.NewButton("Change "+tabBeta.Name+" List", func() {
			makeMobileBetaChangePopUp(MyApp, nameEnt, countBorder, subCountBorder, centeredFinishedCck, addBtn, changeBtn, deleteBtn)
		})

		tab := container.NewBorder(nil, mobileChangeBtn, nil, nil, lst)

		return tab
	} else {
		vbox := container.NewVBox(nameEnt, countBorder, subCountBorder, centeredFinishedCck, addBtn, layout.NewSpacer(), changeBtn, layout.NewSpacer(), centeredCurrentLbl, centeredFinishedLbl, layout.NewSpacer(), deleteBtn)
		tab := container.NewHSplit(lst, vbox)
		tab.Offset = MyApp.App.Preferences().Float("GlobalOffset")

		return tab
	}
}

func makeMobileBetaChangePopUp(MyApp *logic.MyApp, nameEnt fyne.Widget, countBorder *fyne.Container, subCountBorder *fyne.Container, centeredFinishedCck *fyne.Container, addBtn fyne.Widget, changeBtn fyne.Widget, deleteBtn fyne.Widget) {
	var changeBetaPopUp *widget.PopUp

	exitBtn := widget.NewButton("Exit", func() { changeBetaPopUp.Hide() })

	content := container.NewVBox(nameEnt, countBorder, subCountBorder, centeredFinishedCck, addBtn, changeBtn, deleteBtn, exitBtn)

	changeBetaPopUp = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	changeBetaPopUp.Resize(fyne.NewSize(250, 0))
	changeBetaPopUp.Show()
}
