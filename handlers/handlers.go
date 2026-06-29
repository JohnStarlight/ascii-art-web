package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"ascii-art-web/ascii"
)

// PageData είναι το struct που περνάμε στο HTML template
// Χρησιμοποιείται για να εμφανιστούν τα δεδομένα του χρήστη, το αποτέλεσμα και τα errors
type PageData struct {
	Text    string        // κείμενο που έγραψε ο χρήστης
	Banner  string        // banner style που διάλεξε (standard/shadow/thinkertoy)
	Color   string        // χρώμα για το ASCII art (π.χ. red, #ff0000, rgb(255,0,0))
	Letters string        // substring που θα χρωματιστεί — αν κενό, χρωματίζεται όλο
	Result  template.HTML // ASCII art αποτέλεσμα — HTML ώστε να μην escape-αρουν τα <span> tags
	Error   string        // μήνυμα λάθους αν κάτι πάει στραβά
}

// Home χειρίζεται το GET "/" — επιστρέφει την αρχική σελίδα με τη φόρμα
func Home(w http.ResponseWriter, r *http.Request) {
	// Το HandleFunc("/") πιάνει ΟΛΑ τα άγνωστα paths, οπότε ελέγχουμε εδώ για 404
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	// Η αρχική σελίδα δέχεται μόνο GET
	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	// nil γιατί στην αρχική σελίδα δεν έχουμε δεδομένα να περάσουμε στο template
	tmpl.Execute(w, nil)
}

// AsciiArt χειρίζεται το POST "/ascii-art" — παίρνει το κείμενο και το banner,
// καλεί την ascii.Generate() και εμφανίζει το αποτέλεσμα
func AsciiArt(w http.ResponseWriter, r *http.Request) {
	// Αυτό το route δέχεται μόνο POST
	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Διαβάζουμε τα δεδομένα από τη φόρμα
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	color := r.FormValue("color")
	letters := r.FormValue("letters")

	// Αρχικοποιούμε το PageData με τα δεδομένα του χρήστη
	// ώστε να παραμένουν στη φόρμα μετά το submit
	data := PageData{
		Text:    text,
		Banner:  banner,
		Color:   color,
		Letters: letters,
	}

	// Validation: κενό κείμενο → 400
	if text == "" {
		data.Error = "Please enter some text"
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	// Validation: έλεγχος ότι το banner είναι ένα από τα 3 έγκυρα styles
	validBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	if !validBanners[banner] {
		data.Error = "Invalid banner style"
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	// Κλήση της ascii.Generate() από το package του Ατόμου 2
	result, err := ascii.Generate(text, banner, color, letters)
	if err != nil {
		data.Error = err.Error()
		renderTemplate(w, data, http.StatusBadRequest)
		return
	}

	// template.HTML λέει στο Go template "μην κάνεις escape αυτό το string"
	// χρειάζεται γιατί το αποτέλεσμα περιέχει <span> tags για τα χρώματα
	data.Result = template.HTML(result)
	renderTemplate(w, data, http.StatusOK)
}

// renderTemplate φορτώνει και εκτελεί το HTML template με τα δεδομένα και τον HTTP status code
func renderTemplate(w http.ResponseWriter, data PageData, status int) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Γράφουμε πρώτα σε buffer — αν αποτύχει το Execute αφού έχουμε ήδη
	// γράψει στο w, δεν μπορούμε να αλλάξουμε τον status code
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
