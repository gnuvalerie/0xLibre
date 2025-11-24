package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/glog"
)

var adminPassword string
var adminTmpl *template.Template

func initAdmin(pass string) {
	adminPassword = pass
	adminTmpl = template.Must(template.ParseFiles("./templates/admin.html"))
}

func checkAuth(r *http.Request) bool {
	if adminPassword == "" {
		return true
	}
	cookie, err := r.Cookie("admin_auth")
	if err != nil {
		return false
	}
	return cookie.Value == adminPassword
}

func adminPanel(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	pass := r.URL.Query().Get("password")
	uuid := r.URL.Query().Get("uuid")

	if action == "login" && pass != "" {
		hash := sha512.Sum512([]byte(pass))
		hashStr := hex.EncodeToString(hash[:])
		
		if hashStr == adminPassword {
			http.SetCookie(w, &http.Cookie{
				Name:     "admin_auth",
				Value:    adminPassword,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   86400,
			})
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		fmt.Fprintf(w, "wrong password")
		return
	}

	if !checkAuth(r) {
		fmt.Fprintf(w, `<form action="/admin" method="get">
			<input type="hidden" name="action" value="login">
			<input type="password" name="password">
			<button>login</button>
		</form>`)
		return
	}

	if action == "delete" && uuid != "" {
		dirPath := filepath.Join("./storage", uuid)
		if err := os.RemoveAll(dirPath); err != nil {
			glog.Errorf("delete failed: %s", err.Error())
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	files := []struct {
		UUID     string
		Filename string
		Size     int64
		Time     time.Time
	}{}
	
	entries, _ := ioutil.ReadDir("./storage")
	for _, entry := range entries {
		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(filepath.Join("./storage", entry.Name()))
			if err == nil && len(dirFiles) > 0 {
				files = append(files, struct {
					UUID     string
					Filename string
					Size     int64
					Time     time.Time
				}{
					UUID:     entry.Name(),
					Filename: dirFiles[0].Name(),
					Size:     dirFiles[0].Size(),
					Time:     dirFiles[0].ModTime(),
				})
			}
		}
	}
	
	adminTmpl.Execute(w, files)
}
