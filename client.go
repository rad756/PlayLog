package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func upload(filePath string, ip string, port string) {
	// Create a buffer to store the request body
	var buf bytes.Buffer
	url := "http://" + ip + ":" + port + "/upload"

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	// Add a file to the request
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

func download(fileName string, ip string, port string) {
	uri := "http://" + ip + ":" + port + "/" + fileName

	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join("files", fileName)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		panic(err)
	}
}
