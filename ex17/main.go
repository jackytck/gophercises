package main

import (
	"log"

	"github.com/jackytck/gophercises/ex17/vault"
)

func main() {
	v := vault.File("fake-key", ".secrets")
	err := v.Set("demo_key1", "some secret value 1")
	if err != nil {
		panic(err)
	}
	err = v.Set("demo_key2", "some secret value 2")
	if err != nil {
		panic(err)
	}
	err = v.Set("demo_key3", "some secret value 3")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo_key1")
	if err != nil {
		panic(err)
	}
	log.Println(plain)
	plain, err = v.Get("demo_key2")
	if err != nil {
		panic(err)
	}
	log.Println(plain)
	plain, err = v.Get("demo_key3")
	if err != nil {
		panic(err)
	}
	log.Println(plain)
}
