package main

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeMovieTab() fyne.CanvasObject {
	var movieFinishedLbl *widget.Label
	moviesList := readMoviesList()
	genreList := readGenreList()
	selMovieId := -1

	movieLst := widget.NewList(
		func() int {
			return len(moviesList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(moviesList[i].name + " - " + moviesList[i].genre)
		})

	movieNameEnt := widget.NewEntry()
	movieNameEnt.SetPlaceHolder("Enter Movie Name")

	movieGenreDdl := widget.NewSelect(genreList, nil)
	movieGenreDdl.PlaceHolder = "Select Genre"

	movieAddBtn := widget.NewButton("Add Movie", func() {
		var err = []string{}
		if movieNameEnt.Text == "" {
			err = append(err, "Movie name empty")
		}
		if movieGenreDdl.Selected == "" {
			err = append(err, "Genre empty - You can add genres by pressing Change Genres button")
		}
		if hasComma(movieNameEnt.Text) {
			err = append(err, "Movie name cannot contain commas")
		}

		if len(err) != 0 {
			showError(strings.Join(err[:], "\n\n"))
		} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else {
			moviesList = addMovieFunc(movieNameEnt.Text, movieGenreDdl.Selected, moviesList)
			saveMovie(moviesList)
			movieFinishedLbl.SetText(strconv.Itoa(len(moviesList)) + " Movies Watched")
			selMovieId = -1
			movieLst.UnselectAll()
			movieLst.Refresh()
		}

	})

	movieChangeBtn := widget.NewButton("Change Selected Movie", func() {
		var err = []string{}
		if movieNameEnt.Text == "" {
			err = append(err, "Movie name empty")
		}
		if hasComma(movieNameEnt.Text) {
			err = append(err, "Movie name cannot contain commas")
		}
		if selMovieId == -1 {
			err = append(err, "No movie was selected to delete")
		}

		if len(err) != 0 {
			showError(strings.Join(err[:], "\n\n"))
		} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else {
			moviesList = deleteMovieFunc(selMovieId, moviesList)
			moviesList = addMovieFunc(movieNameEnt.Text, movieGenreDdl.Selected, moviesList)
			saveMovie(moviesList)
			selMovieId = -1
			movieLst.UnselectAll()
			movieLst.Refresh()
		}
	})

	movieChangeGenreBtn := widget.NewButton("Change Genres", func() {
		var genrePpu *widget.PopUp
		var genreDdl *widget.Select

		genreEnt := widget.NewEntry()
		genreEnt.SetPlaceHolder("Enter name of genre")
		genreAddBtn := widget.NewButton("Add Platform", func() {
			var err = []string{}
			if genreEnt.Text == "" {
				err = append(err, "Genre name empty")
			}
			if hasComma(genreEnt.Text) {
				err = append(err, "Genre name cannot contain commas")
			}

			if len(err) != 0 {
				showError(strings.Join(err[:], "\n\n"))
			} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
				showServerInaccessibleError()
			} else {
				genreList = addGenreFunc(genreEnt.Text, genreList)
				saveGenre(genreList)
				movieGenreDdl.Options = genreList
				genreDdl.Options = genreList
			}

		})

		genreDdl = widget.NewSelect(genreList, nil)
		genreDeleteBtn := widget.NewButton("Delete Selected Genre", func() {
			if genreDdl.SelectedIndex() == -1 {
				showError("Select genre to delete")
			} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
				showServerInaccessibleError()
			} else {
				genreList = deleteGenreFunc(genreDdl.SelectedIndex(), genreList)
				saveGenre(genreList)
				genreDdl.Options = genreList
				movieGenreDdl.Options = genreList
				genreDdl.ClearSelected()
			}
		})

		exitBtn := widget.NewButton("Exit", func() {
			genrePpu.Hide()
		})
		content := container.New(layout.NewVBoxLayout(), genreEnt, genreAddBtn, genreDdl, genreDeleteBtn, exitBtn)
		genrePpu = widget.NewModalPopUp(content, mainWin.Canvas())
		genrePpu.Resize(fyne.NewSize(200, 0))
		genrePpu.Show()
	})

	movieFinishedLbl = widget.NewLabel(strconv.Itoa(len(moviesList)) + " Movies Watched")
	centeredMovieFinishedLbl := container.New(layout.NewCenterLayout(), movieFinishedLbl)

	movieDeleteBtn := widget.NewButton("Delete Selected Movie", func() {
		if selMovieId == -1 {
			showError("No movie selected to delete")
		} else if serverMode && !isServerAccessible("http://"+serverIP+":"+serverPort) {
			showServerInaccessibleError()
		} else {
			moviesList = deleteMovieFunc(selMovieId, moviesList)
			saveMovie(moviesList)
			selMovieId = -1
			movieLst.UnselectAll()
			movieLst.Refresh()
			movieFinishedLbl.SetText(strconv.Itoa(len(moviesList)) + " Movies Watched")
		}
	})

	movieLst.OnSelected = func(id widget.ListItemID) {
		selMovieId = id

		//Copies selected game name to entry
		movieNameEnt.Text = moviesList[id].name
		movieNameEnt.Refresh()

		//Sets selection in platform dropdown
		for i := range genreList {
			if genreList[i] == moviesList[id].genre {
				movieGenreDdl.SetSelectedIndex(i)
			}
		}
	}

	movieBox := container.New(layout.NewVBoxLayout(), movieNameEnt, movieGenreDdl, movieAddBtn, layout.NewSpacer(), movieChangeBtn, movieChangeGenreBtn, layout.NewSpacer(), centeredMovieFinishedLbl, layout.NewSpacer(), movieDeleteBtn)

	movieTab := container.NewHSplit(movieLst, movieBox)
	movieTab.Offset = splitOffset
	return movieTab
}
