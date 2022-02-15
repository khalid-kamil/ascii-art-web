package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"webpractice/banner"
)

var T = template.Must(template.ParseGlob("static/*.html")) // calling the html file

func AsciiArt(input, ban string) (string, error) { // this function accept two argument, the "input" and  the type of the "banner"
	if input == "" || (ban == "" || !(ban == "standard" || ban == "shadow" || ban == "thinkertoy")) {
		return "", errors.New("invalid input") // return if theirs an error
	}
	return banner.PrintAsciiArt(input, "banner/"+ban+".txt"), nil // return if theirs no error
}

func process(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/ascii-art" {
		http.Error(w, "404 Status not found", http.StatusNotFound)
		return
	}
	// http.Handle("/", http.FileServer(http.Dir("css/")))
	switch r.Method {
	case "GET":
		input := ""
		banner := "standard"
		output, _ := AsciiArt(input, banner)
		if err := T.Execute(w, output); err != nil { // Execute the AsciiArt to prevent printing {{.}}
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		input := r.FormValue("input-name") // getting the input in <textarea>
		banner := r.FormValue("banner")    // getting the banner value in the radio button
		output := ""
		if strings.Contains(input, "\r\n") { // cathing the "return or enter" value
			output = strings.Replace(input, "\r\n", "\\n", -1) // and replace it with newline
		} else {
			output = input
		}
		output, err := AsciiArt(output, banner) //call the AsciiArt() to convert the input into Ascii Art
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		if err := T.Execute(w, output); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", process)
	http.HandleFunc("/ascii-art", process)
	fmt.Printf("Starting server at port 5500\n")
	log.Fatal(http.ListenAndServe(":5500", nil))
}
