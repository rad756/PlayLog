package logic

import (
	"reflect"

	"fyne.io/fyne/v2"
)

type MyApp struct {
	App    fyne.App
	Win    fyne.Window
	Tabs   []Tab
	Mobile bool
}

func Ini(MyApp *MyApp) {
	if MyApp.App.Preferences().BoolWithFallback("FirstRun", true) {
		MyApp.App.Preferences().SetBool("FirstRun", true)
	} else {
		MyApp.App.Preferences().SetBool("FirstRun", false)
	}

	if MyApp.App.Preferences().Bool("FirstRun") {
		CreateTabsFile(MyApp)
		MyApp.App.Preferences().SetFloat("GlobalOffset", 0.6)
	}

	if MyApp.App.Driver().Device().IsMobile() {
		MyApp.Mobile = true
	} else {
		MyApp.Mobile = false
	}

	GetTabs(MyApp)
}

func ServerSetup(ip string, port string, mode string, MyApp *MyApp) {
	MyApp.App.Preferences().SetString("IP", ip)
	MyApp.App.Preferences().SetString("Port", port)
	MyApp.App.Preferences().SetString("StorageMode", mode)
}

// Checks if local and server files are different, returns true if conflict
func FileConflictCheck(MyApp *MyApp) bool {
	files := []string{}
	if TabConflictCheck(MyApp) {
		return true
	}

	Download("Tabs.json", MyApp)

	for _, v := range MyApp.Tabs {
		if v.Mode == "Alpha" {
			files = append(files, v.Name+".json")
			files = append(files, v.Name+"-"+v.Kind+".json")
		} else if v.Mode == "Beta" {
			files = append(files, v.Name+".json")
		}

	}

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

// Checks if local and server tab are different, returns true if conflict
func TabConflictCheck(MyApp *MyApp) bool {
	serverTab := DownloadToMemory("Tabs.json", MyApp)
	localTab := LocalFileToMemory("Tabs.json", MyApp)

	if reflect.DeepEqual(serverTab, localTab) {
		return false
	} else {
		return true
	}
}
