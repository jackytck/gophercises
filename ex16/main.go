package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
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
	// fmt.Printf("%+v\n", creds)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(creds.Key, creds.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", token)

	var conf oauth2.Config
	tclient := conf.Client(context.Background(), &token)
	res2, err := tclient.Get("https://api.twitter.com/1.1/statuses/retweets/1235185319257006080.json")
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()
	io.Copy(os.Stdout, res2.Body)
}
