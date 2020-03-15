package main

import (
	"log"

	"github.com/jackytck/gophercises/ex17/vault"
)

func main() {
	v := vault.Memory("fake-key")
	err := v.Set("demo_key", "some secret value")
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", v)
	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}
	log.Println(plain)
}
