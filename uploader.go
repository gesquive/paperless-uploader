package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Uploader struct {
	baseUrl     string
	apiUrl      string
	token       string
	tokenHeader string
}

func NewUploader(baseUrl string, token string) *Uploader {
	u := new(Uploader)
	u.baseUrl = baseUrl
	u.token = token

	u.apiUrl = fmt.Sprintf("%s/api/documents/post_document/", u.baseUrl)

	u.tokenHeader = fmt.Sprintf("Token %s", token)

	return u
}

func (u Uploader) uploadAll(targetPath string, deleteAfterUpload bool) bool {
	allSuccess := true
	fileStats, err := os.Stat(targetPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Errorf("upload path does not exist, path=%s", targetPath)
		return false
	}

	// if uploadPath is a directory, upload everything in there
	if fileStats.IsDir() {
		d, err := os.Open(targetPath)
		if err != nil {
			log.Errorf("could not open: %s", err)
			allSuccess = false
		}
		fileNames, err := d.Readdirnames(0)
		if err != nil {
			log.Errorf("could not list: %s", err)
			allSuccess = false
		}
		for _, fileName := range fileNames {
			u.uploadAll(path.Join(targetPath, fileName), deleteAfterUpload)
		}
	} else if fileStats.Mode().IsRegular() {
		if !u.uploadFile(targetPath) {
			allSuccess = false
		} else if deleteAfterUpload {
			if !u.deleteFile(targetPath) {
				allSuccess = false
			}
		}
	}
	return allSuccess
}

func (u Uploader) uploadFile(targetPath string) bool {
	log.Infof("action: upload, path: '%s'", targetPath)
	file, _ := os.Open(targetPath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("document", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", u.apiUrl, body)
	r.Header.Add("Authorization", u.tokenHeader)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		log.Errorf("upload error: %s", err)
		return false
	}
	log.Debugf("upload response code=%s", response.Status)
	if response.StatusCode >= 400 {
		log.Errorf("upload error response code=%d", response.StatusCode)
		return false
	}

	return true
}

func (u Uploader) deleteFile(targetPath string) bool {
	log.Infof("action: delete, path: '%s'", targetPath)

	if err := os.Remove(targetPath); err != nil {
		log.Errorf("delete fail: %s", err)
		return false
	}
	return true
}
