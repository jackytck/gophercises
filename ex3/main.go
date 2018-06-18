package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	stories := parseStories("./gopher.json")
	if _, ok := stories["intro"]; !ok {
		panic("intro arc is not found!")
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, stories["intro"])
	})

	http.ListenAndServe(":3000", handler)
}

func parseStories(filePath string) map[string]story {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	stories := make(map[string]story)
	json.Unmarshal(data, &stories)

	return stories
}

type story struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}
