package main

import (
	"bufio"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func loadSetupUI() fyne.CanvasObject {
	questionLbl := widget.NewLabel("Select which mode app will run in, Local or Server Sync")
	localBtn := widget.NewButton("Local Setup", func() { localSetup() })

	serverIpLbl := widget.NewLabel("Enter Server IP Address")
	serverIpEnt := widget.NewEntry()
	serverIpEnt.SetPlaceHolder("Enter Server IP")
	serverPortLbl := widget.NewLabel("Enter Server Port below")
	serverPortEnt := widget.NewEntry()
	serverPortEnt.SetPlaceHolder("Default is 7529")
	serverBtn := widget.NewButton("Server Sync", func() {
		if net.ParseIP(serverIpEnt.Text) != nil { //checks for valid ip
			serverSetup(serverIpEnt.Text, serverPortEnt.Text)
		}
	})

	return container.New(layout.NewVBoxLayout(), questionLbl, localBtn, serverIpLbl, serverIpEnt, serverPortLbl, serverPortEnt, serverBtn)
}

func localSetup() {
	configFile := "conf.csv"
	file, err := os.Create(configFile)

	if err != nil {
		panic(err)
	} else {
		writer := bufio.NewWriter(file)

		writer.WriteString("mode,local\n")
		writer.WriteString("ip,\n")
		writer.WriteString("port,\n")
		if serverDownMode {
			writer.WriteString("serverDownMode,1" + "\n")
		} else {
			writer.WriteString("serverDownMode,0" + "\n")
		}

		writer.Flush()
	}

	mainWin.SetContent(loadMainMenuUI())
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
		if serverDownMode {
			writer.WriteString("serverDownMode,1" + "\n")
		} else {
			writer.WriteString("serverDownMode,0" + "\n")
		}

		writer.Flush()

	}
	mainWin.SetContent(loadMainMenuUI())
}

func writeConfig() {
	configFile := "conf.csv"
	file, err := os.Create(configFile)

	if err != nil {
		panic(err)
	} else {
		writer := bufio.NewWriter(file)

		if serverMode {
			writer.WriteString("mode,sync\n")
		} else {
			writer.WriteString("mode,local\n")
		}
		writer.WriteString("ip," + serverIP + "\n")
		writer.WriteString("port," + serverPort + "\n")
		if serverDownMode {
			writer.WriteString("serverDownMode,1" + "\n")
		} else {
			writer.WriteString("serverDownMode,0" + "\n")
		}

		writer.Flush()

	}
}
