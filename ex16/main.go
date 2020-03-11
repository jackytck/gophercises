package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func main() {
	var (
		keyFile    string
		usersFile  string
		tweetID    string
		numWinners int
	)
	flag.StringVar(&keyFile, "key", ".keys.json", "The file where you store your key and secrets.")
	flag.StringVar(&usersFile, "users", "users.csv", "The file where users who have retweeted the tweet are stored. This will be created if it does not exists.")
	flag.StringVar(&tweetID, "tweet", "1235185319257006080", "The ID of the tweet you wish to find retweeters of.")
	flag.IntVar(&numWinners, "winners", 0, "The number of winners to pick for the contest.")
	flag.Parse()

	key, secret, err := keys(keyFile)
	if err != nil {
		panic(err)
	}

	client, err := twitterClient(key, secret)
	if err != nil {
		panic(err)
	}

	newUsernames, err := retweeters(client, tweetID)
	if err != nil {
		panic(err)
	}

	existingUsernames := existing(usersFile)
	allUsernames := merge(existingUsernames, newUsernames)
	err = writeUsers(usersFile, allUsernames)
	if err != nil {
		panic(err)
	}

	existingUsernames = existing(usersFile)
	if numWinners <= 0 || len(existingUsernames) == 0 {
		return
	}
	winners := pickWinners(existingUsernames, numWinners)
	log.Println("The winners are:")
	for i, username := range winners {
		log.Printf(" %d: %s\n", i+1, username)
	}
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

func twitterClient(key, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("%+v\n", token)

	var conf oauth2.Config
	tclient := conf.Client(context.Background(), &token)
	return tclient, nil
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

func existing(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return []string{}
	}
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}
	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, s := range a {
		uniq[s] = struct{}{}
	}
	for _, s := range b {
		uniq[s] = struct{}{}
	}
	ret := make([]string, 0, len(uniq))
	for s := range uniq {
		ret = append(ret, s)
	}
	sort.Strings(ret)
	return ret
}

func writeUsers(usersFile string, users []string) error {
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, u := range users {
		if err := w.Write([]string{u}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func pickWinners(users []string, numWinners int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(users))
	winners := perm[:numWinners]
	ret := make([]string, 0, numWinners)
	for _, idx := range winners {
		ret = append(ret, users[idx])
	}
	return ret
}

type retweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}
