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
	"strconv"
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

	mux.HandleFunc("/modify/", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(f.Name())[1:]

		modeStr := r.FormValue("mode")

		if modeStr == "" {
			renderModeChoices(w, r, f, ext)
			return
		}

		mode, err := strconv.Atoi(modeStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nStr := r.FormValue("n")

		if nStr == "" {
			renderNumShapeChoices(w, r, f, ext, primitive.Mode(mode))
			return
		}

		_, err = strconv.Atoi(nStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/img/"+filepath.Base(f.Name()), http.StatusFound)

	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]
		onDisk, err := tempfile("", ext)
		defer onDisk.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(onDisk, file)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)
	})

	fs := http.FileServer(http.Dir("./img/"))

	mux.Handle("/img/", http.StripPrefix("/img/", fs))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

type genOpts struct {
	N int
	M primitive.Mode
}

func renderNumShapeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	opts := []genOpts{
		{N: 50, M: mode},
		{N: 100, M: mode},
		{N: 150, M: mode},
		{N: 200, M: mode},
	}

	imgs, err := genImages(rs, ext, opts...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	html := `<html><body>
	{{range .}}
		<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumShapes}}">
			<img style="width:20%;" src="/img/{{.Name}}">
		</a>
	{{end}}
	</body></html>`

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

	tpl.Execute(w, data)
}

func genImages(rs io.ReadSeeker, ext string, opts ...genOpts) ([]string, error) {
	var ret []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		f, err := genImage(rs, ext, opt.N, opt.M)

		if err != nil {
			return nil, err
		}

		ret = append(ret, f)
	}

	return ret, nil
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	opts := []genOpts{
		{N: 10, M: primitive.ModeCircle},
		{N: 10, M: primitive.ModeBeziers},
		{N: 10, M: primitive.ModePolygon},
		{N: 10, M: primitive.ModeCombo},
	}

	imgs, err := genImages(rs, ext, opts...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	html := `<html><body>
	{{range .}}
		<a href="/modify/{{.Name}}?mode={{.Mode}}">
			<img style="width:20%;" src="/img/{{.Name}}">
		</a>
	{{end}}
	</body></html>`

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

	tpl.Execute(w, data)
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
