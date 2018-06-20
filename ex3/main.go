package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web application on")
	file := flag.String("file", "gopher.json", "The JSON file with the CYOA story")
	tpl := flag.String("template", "story.gohtml", "The template html")
	flag.Parse()

	stories := parseStories2(*file)
	handler := newHandler(stories, *tpl)

	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}

func parseStories(filePath string) story {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var stories story
	json.Unmarshal(data, &stories)
	return stories
}

func parseStories2(filePath string) story {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	d := json.NewDecoder(f)
	var stories story
	if err := d.Decode(&stories); err != nil {
		panic(err)
	}
	return stories
}

type story map[string]chapter

type chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func newHandler(stories story, tplPath string) http.Handler {
	if _, ok := stories["intro"]; !ok {
		panic("intro arc is not found!")
	}
	tpl := template.Must(template.ParseFiles(tplPath))

	mux := http.NewServeMux()
	mux.HandleFunc("/", renderStory(tpl, stories["intro"]))

	for k, v := range stories {
		mux.HandleFunc("/"+k, renderStory(tpl, v))
	}

	return mux
}

func renderStory(tp *template.Template, s chapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tp.Execute(w, s)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
	}
}
