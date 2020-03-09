package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func main() {
	var (
		keyFile   string
		usersFile string
		tweetID   string
	)
	flag.StringVar(&keyFile, "key", ".keys.json", "The file where you store your key and secrets.")
	flag.StringVar(&usersFile, "users", "users.csv", "The file where users who have retweeted the tweet are stored. This will be created if it does not exists.")
	flag.StringVar(&tweetID, "tweet", "", "The ID of the tweet you wish to find retweeters of.")
	flag.Parse()

	key, secret, err := keys(keyFile)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", token)

	var conf oauth2.Config
	tclient := conf.Client(context.Background(), &token)

	usernames, err := retweeters(tclient, "1235185319257006080")
	if err != nil {
		panic(err)
	}
	log.Println(usernames)
}

func keys(keyFile string) (key, secret string, err error) {
	var creds struct {
		Key    string `json:"api_key"`
		Secret string `json:"api_secret"`
	}
	f, err := os.Open(keyFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&creds)
	// fmt.Printf("%+v\n", creds)

	return creds.Key, creds.Secret, nil
}

func retweeters(client *http.Client, tweetID string) ([]string, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", tweetID)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var retweets []retweet
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&retweets)
	if err != nil {
		return nil, err
	}

	usernames := make([]string, 0, len(retweets))
	for _, retweet := range retweets {
		usernames = append(usernames, retweet.User.ScreenName)
	}
	return usernames, nil
}

type retweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}
