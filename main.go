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
	Params map[string]string
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

func makePageHandler(templateName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		data := TemplateData{Params: params}
		err := templates.ExecuteTemplate(w, templateName, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
}

func buildMux(config ServeConfig) *mux.Router {
	r := mux.NewRouter()
	for _, page := range config.Pages {
		r.HandleFunc(page.URL, makePageHandler(page.Template))
	}
	return r
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
	funcMap["fetchJSON"] = func(arg1 string) (string, error) {
		u, _ := url.ParseRequestURI(arg1)
		if u.Scheme == "http" || u.Scheme == "https" {
			return fmt.Sprintf("[Would fetch remote URL from \"%s\"]", u.String()), nil
		} else if u.Scheme == "file" {
			return fmt.Sprintf("[Would read file \"%s\"]", u.Path), nil
		} else {
			return "", fmt.Errorf("Unsupported url scheme \"%s\" in URL: \"%s\"", u.Scheme, arg1)
		}
	}
	return funcMap
}
