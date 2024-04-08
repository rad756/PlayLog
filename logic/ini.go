package logic

import (
	"reflect"

	"fyne.io/fyne/v2"
)

type MyApp struct {
	App fyne.App
	Win fyne.Window
}

func Ini(MyApp *MyApp) {
	if MyApp.App.Preferences().BoolWithFallback("FirstRun", true) {
		MyApp.App.Preferences().SetBool("FirstRun", true)
	} else {
		MyApp.App.Preferences().SetBool("FirstRun", false)
	}

	if MyApp.App.Preferences().Bool("FirstRun") {
		SaveAlphaSlice("Game", *MyApp, AlphaSlice{})
		SaveAlphaSlice("Game-Platform", *MyApp, AlphaSlice{})
		SaveAlphaSlice("Movie", *MyApp, AlphaSlice{})
		SaveAlphaSlice("Movie-Genre", *MyApp, AlphaSlice{})
		SaveBetaSlice("Show", *MyApp, BetaSlice{})
		MyApp.App.Preferences().SetFloat("GlobalOffset", 0.6)
	}
}

func ServerSetup(ip string, port string, mode string, MyApp MyApp) {
	MyApp.App.Preferences().SetString("IP", ip)
	MyApp.App.Preferences().SetString("Port", port)
	MyApp.App.Preferences().SetString("Mode", mode)
}

// Checks if local and server files are different, returns true if conflict
func FileConflictCheck(MyApp MyApp) bool {
	files := []string{"Game.json", "Game-Platform.json", "Movie.json", "Movie-Genre.json", "Show.json"}
	filesDownloaded := [][]byte{}
	filesRead := [][]byte{}

	for _, v := range files {
		filesDownloaded = append(filesDownloaded, DownloadToMemory(v, MyApp))
		filesRead = append(filesRead, LocalFileToMemory(v, MyApp))
	}

	if reflect.DeepEqual(filesDownloaded, filesRead) {
		return false
	} else {
		return true
	}

}
