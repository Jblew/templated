package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/mux"
)

type ServeConfig struct {
	StaticDir    string       `json:"staticDir"`
	TemplatesDir string       `json:"templatesDir"`
	Pages        []PageConfig `json:"pages"`
}

type PageConfig struct {
	URL      string `json:"url"`
	Template string `json:"template"`
}

type TemplateData struct {
	Params  map[string]string
	Headers map[string][]string
}

var templates *template.Template

func main() {
	config := loadConfig()

	templatesGlob := fmt.Sprintf("%s/*", config.TemplatesDir)
	templates = template.Must(template.New("base").Funcs(sprig.FuncMap()).Funcs(localFuncMap()).ParseGlob(templatesGlob))

	staticRoute := "/static/"
	staticFileServer := http.FileServer(http.Dir(config.StaticDir))
	http.Handle(staticRoute, http.StripPrefix(staticRoute, staticFileServer))

	r := buildMux(config)
	http.Handle("/", r)

	log.Printf("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func buildMux(config ServeConfig) *mux.Router {
	r := mux.NewRouter()
	for _, page := range config.Pages {
		r.HandleFunc(page.URL, makePageHandler(page.Template))
	}
	return r
}

func makePageHandler(templateName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Printf("Headers: %+v", r.Header)
		data := TemplateData{Params: params, Headers: r.Header}
		err := templates.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}

func loadConfig() ServeConfig {
	bytes, _ := ioutil.ReadFile("serve.json")
	config := ServeConfig{
		StaticDir:    "static",
		TemplatesDir: "templates",
	}
	json.Unmarshal(bytes, &config)
	return config
}

func localFuncMap() map[string]interface{} {
	funcMap := make(map[string]interface{})
	funcMap["fetchJSON"] = func(arg1 string, headers map[string][]string) (map[string]interface{}, error) {
		log.Printf("Headers from context: %+v", headers)
		u, _ := url.ParseRequestURI(arg1)
		if u.Scheme == "http" || u.Scheme == "https" {
			return fetchJSONFromURL(u.String())
		} else if u.Scheme == "file" {
			return readJSONFile(u.Path)
		} else {
			return make(map[string]interface{}), fmt.Errorf("Unsupported url scheme \"%s\" in URL: \"%s\"", u.Scheme, arg1)
		}
	}
	return funcMap
}
