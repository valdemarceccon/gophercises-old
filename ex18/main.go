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
	"text/template"

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

		a, err := genImage(file, ext, 33, primitive.ModeCircle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		file.Seek(0, 0)

		b, err := genImage(file, ext, 33, primitive.ModeEllipse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		file.Seek(0, 0)
		c, err := genImage(file, ext, 33, primitive.ModePolygon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		file.Seek(0, 0)
		d, err := genImage(file, ext, 33, primitive.ModeCombo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		file.Seek(0, 0)

		html := `<html><body>
		{{range .}}
			<img src="/{{.}}">
		{{end}}
		</body></html>`

		tpl := template.Must(template.New("").Parse(html))
		images := []string{a, b, c, d}

		tpl.Execute(w, images)
	})

	fs := http.FileServer(http.Dir("./img/"))

	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))

	if err != nil {
		return "", err
	}

	outFile, err := tempfile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary file")
	}

	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))

}
