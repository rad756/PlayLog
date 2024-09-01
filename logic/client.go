package logic

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func Upload(mar []byte, filePath string, MyApp *MyApp) {
	// Create a buffer to store the request body
	var buf bytes.Buffer
	ip := MyApp.App.Preferences().String("IP")
	port := MyApp.App.Preferences().String("Port")

	url := "http://" + ip + ":" + port + "/upload"

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	file := bytes.NewReader(mar)

	// Create a new form field
	fw, err := w.CreateFormFile("file", filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		log.Fatal(err)
	}

	// Close the multipart writer to finalize the request
	w.Close()

	// Send the request
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func Download(fileName string, MyApp *MyApp) {
	uri := fmt.Sprintf("http://%s:%s/%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"), fileName)

	// if ip == "" {
	// 	ui.ShowError("No server IP configured")
	// }

	// if IsServerAccessible(uri) == false {
	// 	ui.ShowError("Server with " + ip + " inaccessible")
	// }

	resp, _ := http.Get(uri)
	// if err != nil {
	// 	ui.ShowError("Cannot find server with IP: " + ip)
	// }

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
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
}

// Add while this executes, a loading screen
func IsServerAccessible(uri string) bool {
	_, err := http.Get(uri)

	if err != nil {
		return false
	}

	return true
}

func DownloadToMemory(fileName string, MyApp *MyApp) []byte {
	uri := fmt.Sprintf("http://%s:%s/%s", MyApp.App.Preferences().String("IP"), MyApp.App.Preferences().String("Port"), fileName)

	// if ip == "" {
	// 	ui.ShowError("No server IP configured")
	// }

	// if !IsServerAccessible(uri) {
	// 	ui.ShowError("Server with " + ip + " inaccessible")
	// }

	resp, err := http.Get(uri)
	// if err != nil {
	// 	ui.ShowError("Cannot find server with IP: " + ip)
	// }

	defer resp.Body.Close()

	//var buf bytes.Buffer

	body := resp.Body

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, body)

	if err != nil {
		return []byte{}
	}

	return buf.Bytes()
}

func LocalFileToMemory(fileName string, MyApp *MyApp) []byte {
	path, _ := storage.Child(MyApp.App.Storage().RootURI(), fileName)

	file, _ := storage.LoadResourceFromURI(path)

	return file.Content()
}

func DownloadFromServer(MyApp *MyApp) {
	files := []string{}

	for _, v := range MyApp.Tabs {
		if v.Mode == "Alpha" {
			files = append(files, v.Name+".json")
			files = append(files, v.Name+"-"+v.Kind+".json")
		} else if v.Mode == "Beta" {
			files = append(files, v.Name+".json")
		}

	}

	files = append(files, "Tabs.json")

	for _, v := range files {
		Download(v, MyApp)
	}
}

func UploadToServer(MyApp *MyApp) {
	files := []string{}

	for _, v := range MyApp.Tabs {
		if v.Mode == "Alpha" {
			files = append(files, v.Name+".json")
			files = append(files, v.Name+"-"+v.Kind+".json")
		} else if v.Mode == "Beta" {
			files = append(files, v.Name+".json")
		}

	}

	files = append(files, "Tabs.json")

	for _, v := range files {
		Upload(LocalFileToMemory(v, MyApp), v, MyApp)
	}
}

func ChangeServer(MyApp *MyApp, err error, popup *widget.PopUp, ip string, port string, callback func(*MyApp)) {
	if !popup.Hidden {
		popup.Hide()
	}

	if err != nil {
		dialog.ShowError(fmt.Errorf("Cannot connect to %s:%s", ip, port), MyApp.Win)
	} else {
		MyApp.App.Preferences().SetString("IP", ip)
		MyApp.App.Preferences().SetString("Port", port)

		if FileConflictCheck(MyApp) {
			callback(MyApp)
		} else {
			MyApp.App.Preferences().SetString("StorageMode", "Sync")
		}
	}
}
