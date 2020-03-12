package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/jackytck/gophercises/ex16/twitter"
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

	client, err := twitter.New(key, secret)
	if err != nil {
		panic(err)
	}

	newUsernames, err := client.Retweeters(tweetID)
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
