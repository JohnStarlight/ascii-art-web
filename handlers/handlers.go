package handlers

import (
	"html/template"
	"net/http"

	"ascii-art-web/ascii"
)

type PageData struct {
	Text   string
	Banner string
	Result string
	Error  string
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func AsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	data := PageData{
		Text:   text,
		Banner: banner,
	}

	if text == "" {
		data.Error = "Please enter some text"
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	validBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	if !validBanners[banner] {
		data.Error = "Invalid banner style"
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	result, err := ascii.Generate(text, banner)
	if err != nil {
		data.Error = err.Error()
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	data.Result = result
	renderTemplate(w, data, http.StatusOK)
}

func renderTemplate(w http.ResponseWriter, data PageData, status int) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	tmpl.Execute(w, data)
}
