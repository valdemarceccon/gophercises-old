package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/valdemarceccon/gophercises/ex18/primitive"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="image">
				<button type="submit">Upload image</button>
			</form>
		</body></html>`

		fmt.Fprintf(w, html)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]

		out, err := primitive.Transform(file, ext, 100, primitive.WithMode(primitive.ModeCircle))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		outFile, err := tempfile("", ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
		}
		defer outFile.Close()
		io.Copy(outFile, out)

		redirURL := fmt.Sprintf("/%s", outFile.Name())

		http.Redirect(w, r, redirURL, http.StatusFound)
	})

	fs := http.FileServer(http.Dir("./img/"))

	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary file")
	}

	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))

}
