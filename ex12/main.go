package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", true, "whether or not this should be a real or dry run")
	flag.Parse()

	walkDir := "./sample"
	toRename := make(map[string][]file)
	var toRenameRe []string

	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if r, err := match(info.Name()); err == nil {
			key := filepath.Join(filepath.Dir(path), r.base, r.ext)
			toRename[key] = append(toRename[key], file{
				name: info.Name(),
				path: path,
			})
		}

		toRenameRe = append(toRenameRe, path)
		return nil
	})

	for _, files := range toRename {
		n := len(files)
		for i, f := range files {
			origPath := f.path
			r, _ := match(f.name)
			newName := fmt.Sprintf("%s - %d of %d.%s", r.base, i+1, n, r.ext)
			newPath := filepath.Join(filepath.Dir(origPath), newName)
			fmt.Printf("mv %s => %s\n", origPath, newPath)
			if !dry {
				err := os.Rename(origPath, newPath)
				if err != nil {
					fmt.Println("Error renaming:", origPath, err.Error())
				}
			}
		}
	}

	for _, f := range toRenameRe {
		oldName := filepath.Base(f)
		if newName, err := matchRe(oldName); err == nil {
			fmt.Println(oldName, "=>", newName)
		}
	}
}

type file struct {
	name string
	path string
}

type matchResult struct {
	base  string
	index int
	ext   string
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
func match(filename string) (*matchResult, error) {
	// "birthday", "001", "txt"
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s didn't match our pattern", filename)
	}
	// -> Birthday - 1 of 4.txt
	res := matchResult{
		base:  strings.Title(name),
		index: number,
		ext:   ext,
	}
	return &res, nil
}

func matchRe(filename string) (string, error) {
	if !re.MatchString(filename) {
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}
	return re.ReplaceAllString(filename, replaceString), nil
}
