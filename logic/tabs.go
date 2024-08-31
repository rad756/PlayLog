package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"fyne.io/fyne/v2/storage"
)

type Tab struct {
	Mode     string // Alpha or Beta
	Name     string // Manditory
	Kind     string // Optional - needed for Alpha
	Count    string // Optional - needed for Beta
	SubCount string // Optional - needed for Beta
}

func GetTabs(MyApp *MyApp) {
	name := "Tabs.json"
	if PathExists(name, MyApp) {
		path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

		file, _ := storage.LoadResourceFromURI(path)

		json.Unmarshal(file.Content(), &MyApp.Tabs)
	}
}

func DownloadTabs(MyApp *MyApp) error {
	fileName := "Tabs.json"
	uri := fmt.Sprintf("http://%s:%s/%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"), fileName)

	// if ip == "" {
	// 	ui.ShowError("No server IP configured")
	// }

	// if IsServerAccessible(uri) == false {
	// 	ui.ShowError("Server with " + ip + " inaccessible")
	// }

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if !IsTabsFile(data) {
		return errors.New("not tab file")
	}
	// if err != nil {
	// 	panic(err)
	// }
	// filePath := filepath.Join("files", fileName)

	// if err := os.WriteFile(filePath, data, 0644); err != nil {
	// 	panic(err)
	// }

	path, _ := storage.Child(MyApp.App.Storage().RootURI(), fileName)

	file, _ := storage.Writer(path)

	file.Write(data)

	return nil
}

func CreateTabsFile(MyApp *MyApp) {
	name := "Tabs.json"
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

	file, _ := storage.Writer(path)

	mar, _ := json.Marshal(MyApp.Tabs)

	file.Write(mar)

	if MyApp.App.Preferences().String("StorageMode") == "Sync" {
		Upload(mar, name, MyApp)
	}
}

func MoveTabUp(id int, MyApp *MyApp) {
	tabs := MyApp.Tabs

	swapper := reflect.Swapper(tabs)
	swapper(id, id-1)

	MyApp.Tabs = tabs
}

func MoveTabDown(id int, MyApp *MyApp) {
	tabs := MyApp.Tabs

	swapper := reflect.Swapper(tabs)
	swapper(id+1, id)

	MyApp.Tabs = tabs
}

func DeleteTab(id int, MyApp *MyApp) {
	if id < 0 {
		return
	}

	if MyApp.Tabs[id].Mode == "Alpha" {
		deleteTabFiles(MyApp.Tabs[id].Name, MyApp)
		deleteTabFiles(MyApp.Tabs[id].Name+"-"+MyApp.Tabs[id].Kind, MyApp)
	} else if MyApp.Tabs[id].Mode == "Beta" {
		deleteTabFiles(MyApp.Tabs[id].Name, MyApp)
	}
	MyApp.Tabs = append(MyApp.Tabs[:id], MyApp.Tabs[id+1:]...)
	CreateTabsFile(MyApp)
}

func deleteTabFiles(name string, MyApp *MyApp) {
	name = name + ".json"
	if PathExists(name, MyApp) {
		path, _ := storage.Child(MyApp.App.Storage().RootURI(), name)

		_ = storage.Delete(path)
	}
}

func IsTabsFile(t interface{}) bool {
	switch t.(type) {
	case Tab:
		return true
	default:
		return false
	}
}
