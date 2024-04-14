package ui

import (
	"fmt"
	"playlog/logic"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TabAlpha struct {
	Name      string
	Kind      string
	ID        int
	MultiKind bool
}

func NewTabAlpha(alphaSlice *logic.AlphaSlice, MyApp logic.MyApp, tabAlpha TabAlpha, kind *logic.Kind) fyne.CanvasObject {
	var finishedCountLbl *widget.Label
	tabAlpha.MultiKind = false

	lst := widget.NewList(
		func() int {
			return len(alphaSlice.Slice)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s - %s", alphaSlice.Slice[i].Name, alphaSlice.Slice[i].Kind))
		})

	nameEnt := widget.NewEntry()
	nameEnt.SetPlaceHolder(fmt.Sprintf("Enter %s Name", tabAlpha.Name))

	kindSel := widget.NewSelect(kind.Slice, nil)
	kindSel.PlaceHolder = fmt.Sprintf("Select %s", tabAlpha.Kind)
	moreKindBtn := widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		makeMoreKindPopUp(MyApp, tabAlpha, kind, kindSel, tabAlpha)
	})
	clearKindBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		kindSel.SetOptions(kind.Slice)
		kindSel.ClearSelected()
		kindSel.Refresh()
		tabAlpha.ID = -1
	})
	kindBorder := container.NewBorder(nil, nil, clearKindBtn, moreKindBtn, kindSel)

	addBtn := widget.NewButton("Add "+tabAlpha.Name, func() {
		var err []string

		if nameEnt.Text == "" {
			err = append(err, fmt.Sprintf("%s name empty", tabAlpha.Name))
		}
		if kindSel.SelectedIndex() == -1 {
			err = append(err, fmt.Sprintf("%s empty", tabAlpha.Kind))
		}

		if len(err) != 0 {
			ShowError(strings.Join(err[:], "\n\n"), MyApp)
		} else if logic.IsInSyncModeAndServerInaccessible(MyApp) {
			ShowServerInaccessibleError(MyApp)
		} else {
			alphaSlice.AddAlpha(nameEnt.Text, kindSel.Selected, MyApp, tabAlpha.Name)
			tabAlpha.ID = -1
			lst.UnselectAll()
			lst.Refresh()
			finishedCountLbl.SetText(fmt.Sprintf("%d %ss Finished", len(alphaSlice.Slice), tabAlpha.Name))
		}

	})

	changeBtn := widget.NewButton("Change Selected "+tabAlpha.Name, func() {
		var err = []string{}

		if nameEnt.Text == "" {
			err = append(err, fmt.Sprintf("%s name empty", tabAlpha.Name))
		}
		if tabAlpha.ID == -1 {
			err = append(err, fmt.Sprintf("No %s was selected to change", strings.ToLower(tabAlpha.Name)))
		}

		if len(err) != 0 {
			ShowError(strings.Join(err[:], "\n\n"), MyApp)
		} else if logic.IsInSyncModeAndServerInaccessible(MyApp) {
			ShowServerInaccessibleError(MyApp)
		} else {
			alphaSlice.DeleteAlpha(tabAlpha.ID, MyApp, tabAlpha.Name)
			alphaSlice.AddAlpha(nameEnt.Text, kindSel.Selected, MyApp, tabAlpha.Name)
			tabAlpha.ID = -1
			lst.UnselectAll()
			lst.Refresh()
			finishedCountLbl.SetText(fmt.Sprintf("%d %ss Finished", len(alphaSlice.Slice), tabAlpha.Name))
		}
	})

	changeKindBtn := widget.NewButton(fmt.Sprintf("Change %ss", tabAlpha.Kind), func() { makeChangeKindPopUp(MyApp, tabAlpha, kind, kindSel) })

	finishedCountLbl = widget.NewLabel(fmt.Sprintf("%d %ss Finished", len(alphaSlice.Slice), tabAlpha.Name))
	centeredFinishedCountLbl := container.NewCenter(finishedCountLbl)

	deleteBtn := widget.NewButton("Delete Selected "+tabAlpha.Name, func() {
		if tabAlpha.ID == -1 {
			ShowError(fmt.Sprintf("No %s was selected to be deleted", strings.ToLower(tabAlpha.Name)), MyApp)
		} else if logic.IsInSyncModeAndServerInaccessible(MyApp) {
			ShowServerInaccessibleError(MyApp)
		} else {
			alphaSlice.DeleteAlpha(tabAlpha.ID, MyApp, tabAlpha.Name)
			tabAlpha.ID = -1
			lst.UnselectAll()
			lst.Refresh()
			finishedCountLbl.SetText(fmt.Sprintf("%d %ss Finished", len(alphaSlice.Slice), tabAlpha.Name))
		}
	})

	lst.OnSelected = func(id widget.ListItemID) {
		tabAlpha.ID = id
		isInKind := false

		nameEnt.Text = alphaSlice.Slice[id].Name
		nameEnt.Refresh()

		for i := range kind.Slice {
			if kind.Slice[i] == alphaSlice.Slice[id].Kind {
				isInKind = true
				tabAlpha.MultiKind = false
				kindSel.SetOptions(kind.Slice)
				kindSel.SetSelectedIndex(i)
				return
			}
		}

		if !isInKind {
			kindSel.SetOptions([]string{alphaSlice.Slice[id].Kind})
			kindSel.SetSelectedIndex(0)
			tabAlpha.MultiKind = true
		}
	}

	vBox := container.NewVBox(nameEnt, kindBorder, addBtn, layout.NewSpacer(), changeBtn, changeKindBtn, layout.NewSpacer(), centeredFinishedCountLbl, layout.NewSpacer(), deleteBtn)
	tab := container.NewHSplit(lst, vBox)
	tab.Offset = MyApp.App.Preferences().Float("GlobalOffset")

	return tab
}

func makeChangeKindPopUp(MyApp logic.MyApp, ta TabAlpha, k *logic.Kind, tks *widget.Select) {
	var tabKindPopUp *widget.PopUp
	var kindSel *widget.Select

	tabKindEnt := widget.NewEntry()
	tabKindEnt.SetPlaceHolder(fmt.Sprintf("Enter %s Name", ta.Kind))

	addKindBtn := widget.NewButton("Add "+ta.Kind, func() {
		k.AddKind(tabKindEnt.Text, (ta.Name + "-" + ta.Kind), MyApp)

		kindSel.Options = k.Slice
		tks.Options = k.Slice

	})

	kindSel = widget.NewSelect(k.Slice, nil)
	deleteKindBtn := widget.NewButton("Delete Selected "+ta.Kind, func() {
		if kindSel.SelectedIndex() == -1 {
			return
		} else {
			k.DeleteKind(kindSel.SelectedIndex(), (ta.Name + "-" + ta.Kind), MyApp)
			kindSel.Options = k.Slice
			kindSel.ClearSelected()
			tks.Options = k.Slice
			tks.ClearSelected()

		}
	})

	exitBtn := widget.NewButton("Exit", func() { tabKindPopUp.Hide() })

	content := container.NewVBox(tabKindEnt, addKindBtn, kindSel, deleteKindBtn, layout.NewSpacer(), exitBtn)

	tabKindPopUp = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	tabKindPopUp.Resize(fyne.NewSize(200, 0))
	tabKindPopUp.Show()
}

func makeMoreKindPopUp(MyApp logic.MyApp, ta TabAlpha, k *logic.Kind, tks *widget.Select, tabAlpha TabAlpha) {
	var moreKindPopUP *widget.PopUp
	var checkGroup *widget.CheckGroup
	var selectedKind []string

	titleLbl := widget.NewLabel(fmt.Sprintf("---Multi %s Selection---", ta.Kind))
	centeredTitle := container.NewCenter(titleLbl)
	selectedLbl := widget.NewLabel("")
	clearBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		checkGroup.SetSelected([]string{})
		selectedLbl.SetText(strings.Join(selectedKind, " "))
		selectedLbl.Refresh()
	})
	selectedBorder := container.NewBorder(nil, nil, nil, clearBtn, selectedLbl)

	backBtn := widget.NewButton("Back", func() { moreKindPopUP.Hide() })
	saveSelectionBtn := widget.NewButton("Save Selection", func() {
		if len(selectedKind) > 1 {
			tks.SetOptions([]string{strings.Join(selectedKind, " ")})
			tks.SetSelectedIndex(0)
			tks.Refresh()
			moreKindPopUP.Hide()
			tabAlpha.MultiKind = true
		} else {
			ShowError("Need to select more than 1", MyApp)
		}
	})

	checkGroup = widget.NewCheckGroup(k.Slice, func(s []string) {
		selectedKind = s
		selectedLbl.SetText(strings.Join(selectedKind, " "))
		selectedLbl.Refresh()
	})

	scroll := container.NewScroll(checkGroup)
	topVbox := container.NewVBox(centeredTitle, selectedBorder)
	bottomGrid := container.NewAdaptiveGrid(2, backBtn, saveSelectionBtn)

	content := container.NewBorder(topVbox, bottomGrid, nil, nil, scroll)

	moreKindPopUP = widget.NewModalPopUp(content, MyApp.Win.Canvas())
	moreKindPopUP.Resize(fyne.NewSize(250, 350))
	moreKindPopUP.Show()
}
