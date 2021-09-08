package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	staticRoute := "/static/"
	staticFileServer := http.FileServer(http.Dir("./static"))
	http.Handle(staticRoute, http.StripPrefix(staticRoute, staticFileServer))

	http.HandleFunc("/", handleRoot)

	log.Printf("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	fmt.Fprintf(w, "%s", marshallToString(response))
}
