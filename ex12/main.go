package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := "birthday_001.txt"
	// -> Birthday - 1 of 4.txt
	newName, err := match(filename, 4)
	if err != nil {
		fmt.Println("no match")
		os.Exit(1)
	}
	fmt.Println(newName)
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(fileName string, total int) (string, error) {
	// "birthday", "001", "txt"
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	// -> Birthday - 1 of 4.txt
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
