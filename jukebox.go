package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Album struct {
	Title, Artist, Folder string
}

var albums []Album

var html *template.Template

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	if r.FormValue("u") != "" {
		ioutil.WriteFile("/tmp/mplayer", []byte("loadlist "+r.FormValue("u")+"\n"), 0644)
	} else if r.FormValue("f") != "" {
		ioutil.WriteFile("/tmp/mplayer", []byte("loadfile '"+r.FormValue("f")+"'\n"), 0644)
	} else if r.FormValue("d") != "" {
		folder := r.FormValue("d")
		cmd := exec.Command("find", folder, "-type", "f")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		playlist, _ := ioutil.TempFile("", "jukebox")
		ioutil.WriteFile(playlist.Name(), []byte(out.String()), 0644)
		ioutil.WriteFile("/tmp/mplayer", []byte("loadlist '"+playlist.Name()+"'\n"), 0644)
	} else if r.FormValue("c") != "" {
		ioutil.WriteFile("/tmp/mplayer", []byte(r.FormValue("c")+"\n"), 0644)
	}

	html.Execute(w, albums)
}

func buildTemplates() {
	const _html = `<!DOCTYPE html>
<head><title>jukebox</title><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<ul>
  <li><a href="?c=stop">Stop</a>
  <li><a href="?u=http://www.bbc.co.uk/radio/listen/live/r4.asx">Radio 4</a>
  <li><a href="?u=http://www.bbc.co.uk/fivelive/live/live_int.asx">Radio 5 live</a>
  <li><a href="?u=http://somafm.com/startstream=groovesalad.pls">Groove Salad</a>
  {{range .}}<li><a href="?d={{.Folder}}">{{.Artist}} - {{.Title}}</a>{{end}}
</ul>`
	html = template.Must(template.New("html").Parse(_html))
}

func findAlbums(root string) {
	cmd := exec.Command("find", root, "-mindepth", "2", "-maxdepth", "2", "-type", "d")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("findAlbums: ", err)
	}
	lines := strings.Split(string(output), "\n")

	albums = make([]Album, len(lines)-1)
	for i, line := range lines {
		if len(line) > 0 {
			parts := strings.Split(line, "/")
			albums[i] = Album{
				parts[len(parts)-1],
				parts[len(parts)-2],
				line,
			}
		}
	}
}

func main() {
	port := flag.Int("port", 80, "port")
	root := flag.String("root", "", "root")
	flag.Parse()
	if *root == "" {
		panic("root required")
	}

	buildTemplates()
	findAlbums(*root)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
