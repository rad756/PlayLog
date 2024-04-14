package logic

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/storage"
)

func PathExists(s string, MyApp MyApp) bool {
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), s)
	exists, _ := storage.Exists(path)

	return exists
}

func IsNum(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	} else {
		return false
	}
}

func IsInSyncModeAndServerInaccessible(MyApp MyApp) bool {
	if MyApp.App.Preferences().String("StorageMode") == "Sync" && !IsServerAccessible(fmt.Sprintf("http://%s:%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"))) {
		return true
	} else {
		return false
	}
}

func ContainsComma(s string) bool {
	return strings.Contains(s, ",")
}
