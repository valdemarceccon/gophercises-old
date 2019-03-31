package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)

		out, err := primitive.Transform(file, ext, 33)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch ext {
		case ".jpg":
			fallthrough
		case ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		default:
			http.Error(w, fmt.Sprintf("Invalid image type: %s", ext), http.StatusBadRequest)
			return
		}

		io.Copy(w, out)
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
