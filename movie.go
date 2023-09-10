package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var movieFile = filepath.Join("files", "movie.csv")
var genreFile = filepath.Join("files", "movie-type.csv")

type movie struct {
	name  string
	genre string
}

func addMovieFunc(movieName string, movieGenre string, moviesList []movie) []movie {
	moviesList = append(moviesList, movie{name: movieName, genre: movieGenre})

	return orderMoviesList(moviesList)
}

func addGenreFunc(genreName string, genreList []string) []string {
	genreList = append(genreList, genreName)

	return orderGenreList(genreList)
}

func deleteMovieFunc(id int, moviesList []movie) []movie {
	if id >= 0 {
		moviesList = append(moviesList[:id], moviesList[id+1:]...)
		return moviesList
	}

	return moviesList
}

func deleteGenreFunc(id int, genreList []string) []string {
	if id >= 0 {
		genreList = append(genreList[:id], genreList[id+1:]...)
		return genreList
	}

	return genreList
}

func readMoviesList() []movie {
	moviesList := []movie{}

	file, err := os.Open(movieFile)

	if err != nil {
		fmt.Println(err)
		return moviesList
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ",")
			moviesList = append(moviesList, movie{name: s[0], genre: s[1]})
		}
		return moviesList
	}
}

func readGenreList() []string {
	genreList := []string{}

	file, err := os.Open(genreFile)

	if err != nil {
		fmt.Println(err)
		return genreList
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			genreList = append(genreList, scanner.Text())
		}
		return genreList
	}
}

func saveMovie(moviesList []movie) {
	file, _ := os.Create(movieFile)
	writer := bufio.NewWriter(file)

	for i := 0; i < len(moviesList); i++ {
		writer.WriteString(moviesList[i].name + "," + moviesList[i].genre + "\n")
	}

	writer.Flush()

	if serverMode {
		upload(movieFile, serverIP, serverPort)
	}
}

func saveGenre(genreList []string) {
	file, _ := os.Create(genreFile)
	writer := bufio.NewWriter(file)

	for i := 0; i < len(genreList); i++ {
		writer.WriteString(genreList[i] + "\n")
	}

	writer.Flush()

	if serverMode {
		upload(platformFile, serverIP, serverPort)
	}
}

func orderMoviesList(moviesList []movie) []movie {
	sort.Slice(moviesList, func(i, j int) bool {
		return strings.ToLower(moviesList[i].name) < strings.ToLower(moviesList[j].name)
	})

	return moviesList
}

func orderGenreList(genreList []string) []string {
	sort.Slice(genreList, func(i, j int) bool {
		return strings.ToLower(genreList[i]) < strings.ToLower(genreList[j])
	})

	return genreList
}
