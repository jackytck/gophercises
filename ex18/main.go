package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		ext := filepath.Ext(header.Filename)[1:]
		out, err := primitive.Transform(file, ext, 50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch ext {
		case "jpg":
			fallthrough
		case "jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case "png":
			w.Header().Set("Content-Type", "image/png")
		default:
			http.Error(w, fmt.Sprintf("invalid image type %s", ext), http.StatusBadRequest)
		}
		io.Copy(w, out)
	})
	port := "3000"
	log.Printf("Listening at http://127.0.0.1:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
