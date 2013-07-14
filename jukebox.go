package main

import (
	"bytes"
	"flag"
	"html/template"
	"io"
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

var mplayer io.Writer

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "" {
		html.Execute(w, albums)
		return
	}

	if r.FormValue("u") != "" {
		mplayer.Write([]byte("loadlist '" + r.FormValue("u") + "'\n"))
	} else if r.FormValue("f") != "" {
		mplayer.Write([]byte("loadfile '" + r.FormValue("f") + "'\n"))
	} else if r.FormValue("d") != "" {
		folder := r.FormValue("d")
		cmd := exec.Command("find", folder, "-type", "f")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		playlist, _ := ioutil.TempFile("", "jukebox")
		ioutil.WriteFile(playlist.Name(), []byte(out.String()), 0644)
		mplayer.Write([]byte("loadlist '" + playlist.Name() + "'\n"))
	} else if r.FormValue("c") != "" {
		mplayer.Write([]byte(r.FormValue("c") + "\n"))
	}
	http.Redirect(w, r, "", http.StatusFound)
}

func buildTemplates() {
	const _html = `<!DOCTYPE html>
<head><title>jukebox</title><meta name="viewport" content="width=device-width, initial-scale=1.0"></head>
<ul>
  <li><a href="?c=stop">[stop]</a>
  <li><a href="?u=http://www.bbc.co.uk/radio/listen/live/r4.asx">Radio 4</a>
  <li><a href="?u=http://www.bbc.co.uk/fivelive/live/live_int.asx">Radio 5 live</a>
  <li><a href="?u=http://somafm.com/startstream=groovesalad.pls">Groove Salad</a>
  <li><form><input name="f" placeholder="URL"></form>
  {{range .}}<li><a href="?d={{.Folder}}">{{.Artist}} - {{.Title}}</a>{{end}}
</ul>`
	html = template.Must(template.New("html").Parse(_html))
}

func findAlbums(root string) {
	cmd := exec.Command("find", root, "-mindepth", "2", "-maxdepth", "2", "-type", "d")
	if output, err := cmd.Output(); err != nil {
		log.Fatal("find: ", err)
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

func startMPlayer() {
	cmd := exec.Command("mplayer", "-slave", "-really-quiet", "-cache", "32", "-idle")
	mplayer, _ = cmd.StdinPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.Int("port", 80, "port")
	root := flag.String("root", "", "root")
	flag.Parse()
	if len(*root) == 0 {
		panic("root required")
	}

	buildTemplates()
	findAlbums(*root)
	startMPlayer()

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), nil); err != nil {
		log.Fatal("http: ", err)
	}
}
