package ownstory

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Story map[string]Arc

func init() {
	tpl = template.Must(template.New("").Parse(storyTemplate))
}

var tpl *template.Template

const storyTemplate string = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose your own adventure</title>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
      <p>{{.}}</p>
      {{end}}
      {{if .Options}}
      <ul>
      {{range .Options}}
        <li><a href="/{{.Arc}}">{{.Text}}</a></li>
      {{end}}
      </ul>
      {{else}}
      <h3>The end.</h3>
      {{end}}
    </section>
  </body>
</html>`

type HandlerOpts func(*handler)

func WithTemplate(t *template.Template) HandlerOpts {
  return func (h *handler) {
    h.template = t
  }
}

func WithCustomPath(p func(*http.Request) string) HandlerOpts {
  return func(h *handler) {
    h.pathFunc = p
  }
}

func NewHandler(s Story, opts ...HandlerOpts) http.Handler {
  h :=  handler{s, tpl, DefaultPathFunc}
  for _, opt := range opts {
    opt(&h)
  }
	return h
}

type handler struct {
	story Story
  template *template.Template
  pathFunc func(*http.Request) string
}

func DefaultPathFunc(r *http.Request) string {
  path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  path := h.pathFunc(r)

	if arc, ok := h.story[path]; ok {
		err := h.template.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Chapter not found.", http.StatusNotFound)
	}
}

func LoadStoryJSON(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
