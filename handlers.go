package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
    "mime"
	"github.com/golang/glog"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, template.HTML(fmt.Sprintf(`http://%s/`, r.Host)))
}

func upload(w http.ResponseWriter, r *http.Request) {
	glog.Info("Upload request recieved")

	var uuid string = GenerateUUID()
	var filepath string = fmt.Sprintf("./storage/%s/", uuid)

	file, header, err := r.FormFile("file")
	defer func() {
		file.Close()
		glog.Infof(`File "%s" closed.`, header.Filename)
	}()
	if err != nil {
		glog.Errorf("Error retrieving file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request. Error retrieving file.")
		return
	}
	_, err = os.Stat(filepath)
	for !os.IsNotExist(err) {
		uuid = GenerateUUID()
		filepath := fmt.Sprintf("./storage/%s/", uuid)
		_, err = os.Stat(filepath)
	}

	if err := os.MkdirAll(filepath, 0777); err != nil {
		glog.Error("Error saving file on server...")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No storage available.")
		return
	}

	f, err := os.OpenFile(path.Join(filepath, header.Filename), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		glog.Errorf("Error creating file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating file.")
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		glog.Errorf("Error writing file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInsufficientStorage)
		fmt.Fprintf(w, "Insufficient Storage. Error storing file.")
		return
	}

	fmt.Fprintf(w, "OK, Successfully Uploaded\n http://%s/%s\n", r.Host, uuid)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	glog.Info("Retrieve request received")
	uuid := strings.Replace(r.URL.Path[1:], "/", "", -1)
	dirPath := fmt.Sprintf("./storage/%s", uuid)

	files, err := ioutil.ReadDir(dirPath)
	if err != nil || len(files) == 0 {
		glog.Errorf(`File not found: %s`, dirPath)
		http.Error(w, "File Not Found.", http.StatusNotFound)
		return
	}

	filename := files[0].Name()
	filePath := path.Join(dirPath, filename)

	ext := path.Ext(filename)
	mtype := mime.TypeByExtension(ext)
	if mtype == "" {
		mtype = "application/octet-stream"
	}

	if strings.HasPrefix(mtype, "text/") ||
		strings.HasPrefix(mtype, "image/") ||
		mtype == "application/pdf" {
		w.Header().Set("Content-Disposition", `inline; filename="`+filename+`"`)
	} else {
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	}

	w.Header().Set("Content-Type", mtype)
	http.ServeFile(w, r, filePath)
}
