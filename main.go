package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

var port *uint64
var tmpl *template.Template
var host *string

func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.Header.Get("Content-type"), "multipart/form-data"):
		upload(w, r)
	case uuidMatch.MatchString(r.URL.Path):
		getFile(w, r)
	default:
		home(w, r)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	tmpl = template.Must(template.ParseFiles("./templates/index.html"))
	host = flag.String("h", "127.0.0.1", "Address to serve on")
	port = flag.Uint64("p", 8000, "port")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./0xlibre -p=8080 -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	glog.Flush()
	http.HandleFunc("/", router)
	http.ListenAndServe(fmt.Sprintf("%s:%d",*host,*port), nil)
}
