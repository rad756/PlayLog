package main

import (
	"bufio"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// variables
var splitOffset = 0.6
var a = app.New()
var mainWin = a.NewWindow("PlayLog")
var firstRun bool // default false
var serverMode = true
var serverIP string
var serverPort string

func main() {
	//ini()
	mainWin.Resize(fyne.NewSize(600, 0))
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	var content fyne.CanvasObject

	content = ini()

	mainWin.SetContent(content)
	mainWin.SetIcon(icon)
	mainWin.ShowAndRun()
}

// checks if string contains a comma
func noComma(s string) bool {
	return !strings.Contains(s, ",")
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
			}
		}
	}
	if serverMode && !firstRun && isServerAccessible("http://"+serverIP+":"+serverPort) {
		files := []string{"game.csv", "game-type.csv", "movie.csv", "movie-type.csv", "show.csv"}

		for _, v := range files {
			download(v, serverIP, serverPort)
		}
	}

	if firstRun {
		return loadSetupUI()
	} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
		return startUpError("Server with IP " + serverIP + " is inaccessible")
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
