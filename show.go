package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var showFile = filepath.Join("files", "show.csv")

type show struct {
	name     string
	season   int
	episode  int
	finished string
}

func readShowsList() []show {
	showsList := []show{}

	file, err := os.Open(showFile)

	if err != nil {
		fmt.Println(err)
		return showsList
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ",")
			showSeason, _ := strconv.Atoi(s[1])
			showEpisode, _ := strconv.Atoi(s[2])
			showsList = append(showsList, show{name: s[0], season: showSeason, episode: showEpisode, finished: s[3]})
		}
		return showsList
	}
}

func saveShow(showsList []show) {
	file, _ := os.Create(showFile)
	writer := bufio.NewWriter(file)

	for i := 0; i < len(showsList); i++ {
		writer.WriteString(showsList[i].name + "," + strconv.Itoa(showsList[i].season) + "," + strconv.Itoa(showsList[i].episode) + "," + showsList[i].finished + "\n")
	}

	writer.Flush()

	if serverMode {
		upload(showFile, serverIP, serverPort)
	}
}

func addShowFunc(showName string, showSeason int, showEpisode int, showFinished bool, showsList []show) []show {
	fin := "No"
	if showFinished {
		fin = "Yes"
	}

	showsList = append(showsList, show{name: showName, season: showSeason, episode: showEpisode, finished: fin})

	return orderShowList(showsList)
}

func deleteShowFunc(id int, showsList []show) []show {
	if id >= 0 {
		showsList = append(showsList[:id], showsList[id+1:]...)
		return showsList
	}

	return showsList
}

func orderShowList(showsList []show) []show {
	sort.Slice(showsList, func(i, j int) bool {
		return strings.ToLower(showsList[i].name) < strings.ToLower(showsList[j].name)
	})

	return showsList
}

func isNum(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	} else {
		return false
	}

}

func countWatching(showsList []show) string {
	x := 0

	for i := range showsList {
		if showsList[i].finished == "No" {
			x++
		}
	}

	return strconv.Itoa(x)
}

func countWatched(showsList []show) string {
	x := 0

	for i := range showsList {
		if showsList[i].finished == "Yes" {
			x++
		}
	}

	return strconv.Itoa(x)
}

func moreFunc(s string) string {
	i, _ := strconv.Atoi(s)
	i++
	return strconv.Itoa(i)
}

func lessFunc(s string) string {
	i, _ := strconv.Atoi(s)
	i--
	return strconv.Itoa(i)
}
