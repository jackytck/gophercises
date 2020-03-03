package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var creds struct {
		Key    string `json:"api_key"`
		Secret string `json:"api_secret"`
	}
	f, err := os.Open(".keys.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&creds)
	fmt.Printf("%+v\n", creds)
}
