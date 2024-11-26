// Copyright 2024 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/lf-edge/ekuiper/internal/conf"
)

func zipFile(targetFilePath string, tmpFilename string) (string, string, error) {
	tmpdir := os.TempDir()
	tmpZipFile, err := os.Create(filepath.Join(tmpdir, tmpFilename))
	if err != nil {
		return "", "", err
	}
	defer tmpZipFile.Close()
	zipWriter := zip.NewWriter(tmpZipFile)
	defer zipWriter.Close()
	fileToZip, err := os.Open(targetFilePath)
	if err != nil {
		return "", "", err
	}
	defer fileToZip.Close()
	fileInfo, err := fileToZip.Stat()
	if err != nil {
		return "", "", err
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return "", "", err
	}
	header.Name = filepath.Base(targetFilePath)
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return "", "", err
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return "", "", err
	}
	return tmpdir, tmpZipFile.Name(), nil
}

func logsDownloadHandler(w http.ResponseWriter, r *http.Request) {
	logPath, err := conf.GetLogLoc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	name := vars["filename"]
	fp := filepath.Join(logPath, name)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	tmpFilename := fmt.Sprintf("%s.zip", name)
	_, zipFilePath, err := zipFile(fp, tmpFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		os.Remove(zipFilePath)
	}()
	downloadHandler(zipFilePath, w, r)
}

func downloadHandler(targetFilePath string, w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	file, err := os.Open(targetFilePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(targetFilePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	logPath, err := conf.GetLogLoc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	files, err := listDirectory(logPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse(files, w, logger)
}

func listDirectory(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}
	return files, nil
}
