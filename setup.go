package main

import (
	"bufio"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func setupui() fyne.CanvasObject {
	questionLbl := widget.NewLabel("Select which mode app will run in, Local or Server Sync")
	localBtn := widget.NewButton("Local Setup", func() { localSetup() })

	serverIpLbl := widget.NewLabel("Enter Server IP Address below MANDITORY FOR SERVER SYNC")
	serverIpEnt := widget.NewEntry()
	serverIpEnt.SetPlaceHolder("Enter Server IP")
	serverPortLbl := widget.NewLabel("Enter Server Port below")
	serverPortEnt := widget.NewEntry()
	serverPortEnt.SetPlaceHolder("Default is 7529")
	serverBtn := widget.NewButton("Server Sync", func() { serverSetup(serverIpEnt.Text, serverPortEnt.Text) })
	infoLbl := widget.NewLabel("After selection app will close, opening the app again will open in chosen mode")

	return container.New(layout.NewVBoxLayout(), questionLbl, localBtn, serverIpLbl, serverIpEnt, serverPortLbl, serverPortEnt, serverBtn, infoLbl)
}

func localSetup() {
	configFile := "conf.csv"
	file, err := os.Create(configFile)

	if err != nil {
		panic(err)
	} else {
		writer := bufio.NewWriter(file)

		writer.WriteString("mode,local\n")
		writer.WriteString("ip,198.51.100.1\n")
		writer.WriteString("port,7529\n")

		writer.Flush()
	}

	os.Exit(1)
}

func serverSetup(ip string, port string) {
	if port == "" {
		port = "7529"
	}

	configFile := "conf.csv"
	file, err := os.Create(configFile)

	if err != nil {
		panic(err)
	} else {
		writer := bufio.NewWriter(file)

		writer.WriteString("mode,sync\n")
		writer.WriteString("ip," + ip + "\n")
		writer.WriteString("port," + port + "\n")

		writer.Flush()

	}
	os.Exit(1)
}
