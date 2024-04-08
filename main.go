package main

import (
	"playlog/logic"
	"playlog/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	MyApp := &logic.MyApp{App: app.New()}
	MyApp.Win = MyApp.App.NewWindow("Playlog")

	if MyApp.App.Metadata().Release {
		MyApp.App = app.NewWithID("com.github.rad756.playlog")
	} else {
		MyApp.App = app.NewWithID("com.github.rad756.playlog.testing")
	}

	MyApp.Win.Resize(fyne.NewSize(600, 400))

	logic.Ini(MyApp)

	MyApp.Win.SetContent(ui.LoadGUI(*MyApp))
	MyApp.Win.ShowAndRun()
}
