package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// variables
var splitOffset = 0.6
var a = app.New()
var mainWin = a.NewWindow("PlayLog")
var firstRun bool // default false
var serverMode bool
var serverDownMode = false //if app is in offline mode
var serverIP string
var serverPort string

func main() {
	mainWin.Resize(fyne.NewSize(600, 0))
	icon, _ := fyne.LoadResourceFromPath("Icon.png")

	content := ini()
	serverSync()

	mainWin.SetContent(content)
	mainWin.SetIcon(icon)
	mainWin.ShowAndRun()
}

// checks if string contains a comma
func hasComma(s string) bool {
	return strings.Contains(s, ",")
}

func ini() fyne.CanvasObject {
	//checks if dir files exists, if not creates it
	if _, err := os.Stat("files"); os.IsNotExist(err) {
		os.Mkdir("files", 0777)
	}
	if _, err := os.Stat("conf.csv"); os.IsNotExist(err) {
		firstRun = true
	} else {
		file, err := os.Open("conf.csv")

		if err != nil {
			panic(err)
		} else {
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				s := strings.Split(scanner.Text(), ",")
				s[1] = strings.TrimSuffix(s[1], "\n")
				if s[0] == "mode" {
					if s[1] == "local" {
						serverMode = false
					} else {
						serverMode = true
					}
				}
				if s[0] == "ip" {
					serverIP = s[1]
				}
				if s[0] == "port" {
					serverPort = s[1]
				}
				if s[0] == "serverDownMode" {
					if s[1] == "1" {
						serverDownMode = true
					} else {
						serverDownMode = false
					}
				}
			}
		}
	}

	serverUP := isServerAccessible("http://" + serverIP + ":" + serverPort)

	if firstRun {
		return loadSetupUI()
	} else if serverMode && !serverUP {
		return startUpServerError()
	} else if !serverMode && serverDownMode && serverUP && fileConflictCheck() {
		return loadSyncUI()
	} else {
		return loadMainMenuUI()
	}
}

func loadMainMenuUI() fyne.CanvasObject {
	gameTab := makeGameTab()
	movieTab := makeMovieTab()
	showTab := makeShowTab()

	content := container.NewAppTabs(
		container.NewTabItem("Games", gameTab),
		container.NewTabItem("Movies", movieTab),
		container.NewTabItem("Shows", showTab))

	return content
}

func loadSyncUI() fyne.CanvasObject {
	errorLbl := widget.NewLabel("-- Sync Error --")
	errorCentered := container.NewCenter(errorLbl)
	questionLbl := widget.NewLabel("Do you want to download server files OR upload local files to server?")
	questionCentered := container.NewCenter(questionLbl)

	serverFilesBtn := widget.NewButton("Download Server Files", func() {
		downloadFromServer()
		writeConfig()
		mainWin.SetContent(loadMainMenuUI())
	})
	orLbl := widget.NewLabel("OR")
	orCentered := container.NewCenter(orLbl)
	localFilesBtn := widget.NewButton("Upload Local Files", func() {
		//Does not upload after loads main menu
		uploadToServer()
		writeConfig()
		mainWin.SetContent(loadMainMenuUI())
	})

	return container.NewVBox(errorCentered, questionCentered, serverFilesBtn, orCentered, localFilesBtn)

}

func serverSync() {
	serverUP := isServerAccessible("http://" + serverIP + ":" + serverPort)

	if serverDownMode && serverUP {
		fmt.Println(fileConflictCheck())
	} else if serverMode && serverUP {
		downloadFromServer()
	}
}

func downloadFromServer() {
	files := []string{"game.csv", "game-type.csv", "movie.csv", "movie-type.csv", "show.csv"}

	for _, v := range files {
		download(v, serverIP, serverPort)
	}
}

func uploadToServer() {
	files := []string{"game.csv", "game-type.csv", "movie.csv", "movie-type.csv", "show.csv"}

	for _, v := range files {
		upload(filepath.Join("files", v), serverIP, serverPort)
	}
}

// Checks if local and server files are different, returns true if conflict
func fileConflictCheck() bool {
	files := []string{"game.csv", "game-type.csv", "movie.csv", "movie-type.csv", "show.csv"}
	filesDownloaded := [][]string{}
	filesRead := [][]string{}

	for _, v := range files {
		filesDownloaded = append(filesDownloaded, downloadToMemory(v, serverIP, serverPort))
		//fmt.Println(downloadToMemory(v, serverIP, serverPort))
		filesRead = append(filesRead, localFileToMemory(v))
	}

	if reflect.DeepEqual(filesDownloaded, filesRead) {
		return false
	} else {
		return true
	}
}
