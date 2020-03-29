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
	"strconv"

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
		log.Println("upload")
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)[1:]
		onDisk, err := tempfile("", ext)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer onDisk.Close()
		_, err = io.Copy(onDisk, file)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)
	})

	mux.HandleFunc("/modify/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle modify")
		f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		ext := filepath.Ext(f.Name())[1:]
		modeStr := r.FormValue("mode")
		numStr := r.FormValue("n")

		// render num mode choices
		if modeStr == "" {
			renderModeChoices(w, r, f, ext)
			return
		}
		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// render num shapes choices
		if numStr == "" {
			renderNumShapeChoices(w, r, f, ext, primitive.Mode(mode))
			return
		}
		_, err = strconv.Atoi(numStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/img/"+filepath.Base(f.Name()), http.StatusFound)
	})

	// static image server
	fs := http.FileServer(http.Dir("./img"))
	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	port := "3000"
	log.Printf("Listening at http://127.0.0.1:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	log.Println("renderModeChoices")
	opts := []genOpts{
		{N: 50, M: primitive.ModeCircle},
		{N: 50, M: primitive.ModeTriangle},
		{N: 50, M: primitive.ModePolygon},
		{N: 50, M: primitive.ModeCombo},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}">
					<img style="width: 24%;" src="/img/{{.Name}}" />
				</a>
			{{end}}
		</body><html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name string
		Mode primitive.Mode
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name: filepath.Base(img),
			Mode: opts[i].M,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func renderNumShapeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	log.Println("renderNumShapeChoices")
	opts := []genOpts{
		{N: 30, M: mode},
		{N: 40, M: mode},
		{N: 50, M: mode},
		{N: 60, M: mode},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumShapes}}">
					<img style="width: 24%;" src="/img/{{.Name}}" />
				</a>
			{{end}}
		</body><html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name      string
		Mode      primitive.Mode
		NumShapes int
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name:      filepath.Base(img),
			Mode:      opts[i].M,
			NumShapes: opts[i].N,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type genOpts struct {
	N int
	M primitive.Mode
}

func genImages(rs io.ReadSeeker, ext string, opts ...genOpts) ([]string, error) {
	var ret []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		img, err := genImage(rs, ext, opt.N, opt.M)
		if err != nil {
			return ret, err
		}
		ret = append(ret, img)
	}
	return ret, nil
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
