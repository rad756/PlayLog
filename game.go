package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var gameFile = filepath.Join("files", "game.csv")
var platformFile = filepath.Join("files", "game-type.csv")

type game struct {
	name     string
	platform string
}

func addGameFunc(gameName string, gamePlatform string, gamesList []game) []game {
	gamesList = append(gamesList, game{name: gameName, platform: gamePlatform})

	return orderGamesList(gamesList)
}

func addPlatformFunc(platformName string, platformList []string) []string {
	platformList = append(platformList, platformName)

	return orderPlatformList(platformList)
}

func deleteGameFunc(id int, gamesList []game) []game {
	if id >= 0 {
		gamesList = append(gamesList[:id], gamesList[id+1:]...)
		return gamesList
	}

	return gamesList
}

func deletePlatformFunc(id int, platformList []string) []string {
	if id >= 0 {
		platformList = append(platformList[:id], platformList[id+1:]...)
		return platformList
	}

	return platformList
}

func readGamesList() []game {
	gamesList := []game{}

	file, err := os.Open(gameFile)

	if err != nil {
		fmt.Println(err)
		return gamesList
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ",")
			gamesList = append(gamesList, game{name: s[0], platform: s[1]})
		}
		return gamesList
	}
}

func readPlatformList() []string {
	platformList := []string{}

	file, err := os.Open(platformFile)

	if err != nil {
		fmt.Println(err)
		return platformList
	} else {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			platformList = append(platformList, scanner.Text())
		}
		return platformList
	}
}

func saveGame(gamesList []game) {
	file, _ := os.Create(gameFile)
	writer := bufio.NewWriter(file)

	for i := 0; i < len(gamesList); i++ {
		writer.WriteString(gamesList[i].name + "," + gamesList[i].platform + "\n")
	}

	writer.Flush()
}

func savePlatform(platformList []string) {
	file, _ := os.Create(platformFile)
	writer := bufio.NewWriter(file)

	for i := 0; i < len(platformList); i++ {
		writer.WriteString(platformList[i] + "\n")
	}

	writer.Flush()
}

func orderGamesList(gamesList []game) []game {
	sort.Slice(gamesList, func(i, j int) bool {
		return strings.ToLower(gamesList[i].name) < strings.ToLower(gamesList[j].name)
	})

	return gamesList
}

func orderPlatformList(platformList []string) []string {
	sort.Slice(platformList, func(i, j int) bool {
		return strings.ToLower(platformList[i]) < strings.ToLower(platformList[j])
	})

	return platformList
}
