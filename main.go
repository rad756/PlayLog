package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// variables
var splitOffset = 0.6
var a = app.New()
var mainWin = a.NewWindow("PlayLog")

func main() {
	mainWin.Resize(fyne.NewSize(600, 400))
	icon, _ := fyne.LoadResourceFromPath("Icon.png")

	gameTab := makeGameTab()
	movieTab := makeMovieTab()
	showTab := makeShowTab()

	mainTab := container.NewAppTabs(
		container.NewTabItem("Games", gameTab),
		container.NewTabItem("Movies", movieTab),
		container.NewTabItem("Shows", showTab))

	mainWin.SetContent(mainTab)
	mainWin.SetIcon(icon)
	mainWin.ShowAndRun()
}

// checks if string contains a comma
func noComma(s string) bool {
	return !strings.Contains(s, ",")
}
