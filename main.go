package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	staticRoute := "/static/"
	staticFileServer := http.FileServer(http.Dir("./static"))
	http.Handle(staticRoute, http.StripPrefix(staticRoute, staticFileServer))

	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/role/{name}", handleRole)
	http.Handle("/", r)

	log.Printf("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func handleRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "name: %v\n", vars["name"])
}
