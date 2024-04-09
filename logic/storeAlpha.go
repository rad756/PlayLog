package logic

import (
	"encoding/json"

	"fyne.io/fyne/v2/storage"
)

func ReadAlphaSlice(name string, MyApp MyApp) *AlphaSlice {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.LoadResourceFromURI(path)

	var as *AlphaSlice

	json.Unmarshal(file.Content(), &as)

	return as
}

func SaveAlphaSlice(name string, MyApp MyApp, as AlphaSlice) {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.Writer(path)

	mar, _ := json.Marshal(as)

	file.Write(mar)

	if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		Upload(mar, name, MyApp)
	}
}

func ReadAlphaKind(name string, MyApp MyApp) *Kind {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.LoadResourceFromURI(path)

	var k *Kind

	json.Unmarshal(file.Content(), &k)

	return k
}

func SaveAlphaKind(name string, MyApp MyApp, k Kind) {
	name = name + ".json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.Writer(path)

	mar, _ := json.Marshal(k)

	file.Write(mar)

	if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		Upload(mar, name, MyApp)
	}
}
