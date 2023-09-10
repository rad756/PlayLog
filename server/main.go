package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var serverPort = "7529"                   // 7529 default port
var maxUploadSize int64 = 2 * 1024 * 1024 // 2MB default max

func main() {
	ini()
	// serves files
	fileHandler := http.FileServer(http.Dir("files"))
	http.Handle("/", fileHandler)
	http.HandleFunc("/upload", upload())

	http.ListenAndServe(":"+serverPort, nil)

}

func upload() http.HandlerFunc {
	filePath := "files/"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		// parse and validate file and post parameters
		file, fileHeader, err := r.FormFile("file")

		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		// Get and print out file size
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		// validate file size
		if fileSize > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// write file
		newFile, err := os.Create(filePath + fileHeader.Filename)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("SUCCESS - use /files/%v to access the file", filePath+fileHeader.Filename)))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func ini() {
	//checks if dir files exists, if not creates it
	if _, err := os.Stat("files"); os.IsNotExist(err) {
		os.Mkdir("files", 0777)
		files := []string{"game.csv", "game-type.csv", "movie.csv", "movie-type.csv", "show.csv"}

		for _, v := range files {
			os.Create(filepath.Join("files", v))
		}
	}

	args := os.Args[1:] //Cuts out program path

	if len(args) == 1 {
		serverPort = args[0]
	} else if len(args) == 2 {
		serverPort = args[0]
		maxUploadSize, _ = strconv.ParseInt(args[1], 10, 64)
	}
}
