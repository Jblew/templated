package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type ServeConfig struct {
	Pages []PageConfig `json:"pages"`
}

type PageConfig struct {
	URL      string `json:"url"`
	Template string `json:"template"`
}

type TemplateData struct {
	Params map[string]string
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	config := loadConfig()
	staticRoute := "/static/"
	staticFileServer := http.FileServer(http.Dir("./static"))
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
	var config ServeConfig
	json.Unmarshal(bytes, &config)
	return config
}
