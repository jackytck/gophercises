package main

import (
	"fmt"
	"regexp"

	phonedb "github.com/jackytck/gophercises/ex8/db"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jacky"
	password = "natnat"
	dbname   = "gophercises_phone"
)

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Seed())

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				// delete this number
				must(db.DeletePhone(p.ID))
			} else {
				// update this number
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No change required")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
