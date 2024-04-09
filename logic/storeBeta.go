package logic

import (
	"encoding/json"

	"fyne.io/fyne/v2/storage"
)

func ReadBetaSlice(name string, MyApp MyApp) *BetaSlice {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.LoadResourceFromURI(path)

	var bs *BetaSlice

	json.Unmarshal(file.Content(), &bs)

	return bs
}

func SaveBetaSlice(name string, MyApp MyApp, bs BetaSlice) {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.Writer(path)

	mar, _ := json.Marshal(bs)

	file.Write(mar)

	if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		Upload(mar, name, MyApp)
	}
}
