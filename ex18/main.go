package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackytck/gophercises/ex18/primitive"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="image">
				<button type="submit">Upload Image</button>
			</form>
			</body></html>
		`
		fmt.Fprint(w, html)
	})

	// handle image upload
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)[1:]
		a, err := genImage(file, ext, 50, primitive.ModeCircle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		b, err := genImage(file, ext, 50, primitive.ModeTriangle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		c, err := genImage(file, ext, 50, primitive.ModePolygon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		d, err := genImage(file, ext, 50, primitive.ModeCombo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html := `<html><body>
			{{range .}}
				<img src="/{{.}}" />
			{{end}}
		</body><html>`
		tpl := template.Must(template.New("").Parse(html))
		tpl.Execute(w, []string{a, b, c, d})
	})

	// static image server
	fs := http.FileServer(http.Dir("./img"))
	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	port := "3000"
	log.Printf("Listening at http://127.0.0.1:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	// save image file to /img
	outFile, err := tempfile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func tempfile(prefix, ext string) (*os.File, error) {
	f, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary file")
	}
	defer f.Close()
	defer os.Remove(f.Name())
	return os.Create(fmt.Sprintf("%s.%s", f.Name(), ext))

}
