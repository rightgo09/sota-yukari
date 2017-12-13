package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"
	"time"

	"github.com/rightgo09/sota-yukari/docomo"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	apiKey := os.Getenv("DOCOMO_APIKEY")
	if apiKey == "" {
		panic("DOCOMO_APIKEY must be set")
	}

	client := docomo.NewClient(apiKey)

	go sotaSpeak(client)

	http.Handle("/", &templateHandler{filename: "main.html"})
	http.HandleFunc("/say", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		word := r.Form["w"][0]
		client.Q <- word
		w.Write([]byte("ok"))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServeTLS(":12345", "ssl/myself.crt", "ssl/myself.key", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func sotaSpeak(client *docomo.Client) {
	var err error
	for {
		w := <-client.Q
		log.Println("sotaSpeak(" + w + ")")
		if runtime.GOOS == "linux" {
			yukari := docomo.Yukari(w)
			now := time.Now().Format("20060102150405")
			rawName := now + ".raw"
			err = client.Synthesize(yukari, rawName)
			if err != nil {
				log.Println(err)
				continue
			}
			cmd := exec.Command(
				"aplay", "-q",
				"-t", "raw",
				"-r", "16k",
				"-c", "1",
				"-f", "S16_BE",
				rawName)
			err = cmd.Run()
			if err != nil {
				log.Println(err)
			}
		} else if runtime.GOOS == "darwin" {
			cmd := exec.Command("say", "-v", "Kyoko", w)
			err = cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
