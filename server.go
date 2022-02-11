package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"webpractice/banner"
)

var T = template.Must(template.ParseGlob("static/*.html"))

func AsciiArtOutput(input, ban string) string {
	return banner.PrintAsciiArt(input, "banner/"+ban+".txt")
}
func process(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		input := ""
		banner := "standard"
		output := AsciiArtOutput(input, banner)
		if err := T.Execute(w, output); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		input := r.FormValue("input-name")
		banner := r.FormValue("banner")
		if len(input) == 0 && len(banner) == 0 {
			input = ""
			banner = "shadow"
		}
		output := ""
		if strings.Contains(input, "\r\n") {
			output = strings.Replace(input, "\r\n", "\\n", -1)
		} else {
			output = input
		}
		output = AsciiArtOutput(output, banner)

		if err := T.Execute(w, output); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
func main() {
	http.HandleFunc("/", process)
	fmt.Printf("Starting server at port 2000\n")
	log.Fatal(http.ListenAndServe(":2000", nil))
}
