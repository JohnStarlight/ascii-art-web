# PRD — ascii-art-web (Zone01)

## Σκοπός

Δημιουργία web εφαρμογής σε Go που δέχεται κείμενο από τον χρήστη μέσω φόρμας και επιστρέφει το αντίστοιχο ASCII art, χρησιμοποιώντας banner αρχεία (standard, shadow, thinkertoy). Υποστηρίζει χρωματισμό του output μέσω dropdown επιλογής χρώματος και προαιρετικού substring. Χρησιμοποιείται **μόνο** η standard library της Go.

---

## Τεχνικές Απαιτήσεις (Zone01)

- Γλώσσα: **Go**
- Packages: μόνο standard library
- Port: **8080**
- HTTP methods: GET (αρχική σελίδα), POST (αποστολή φόρμας)
- Banner styles: `standard`, `shadow`, `thinkertoy`
- HTTP status codes:
  - `200 OK` — επιτυχής απόκριση
  - `400 Bad Request` — λανθασμένο input
  - `404 Not Found` — άγνωστο path
  - `500 Internal Server Error` — σφάλμα server
- Το πρόγραμμα **δεν πρέπει να κρασάρει** σε καμία περίπτωση

---

## Δομή Project

```
ascii-art-web/
├── main.go
├── handlers/
│   └── handlers.go
├── ascii/
│   └── ascii.go
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
└── templates/
    └── index.html
```

---

## Χωρισμός Εργασίας

---

### Άτομο 1 — HTTP Server & Routing

**Αρχεία:** `main.go`, `handlers/handlers.go`

#### Καθήκοντα

1. **`main.go`** — εκκίνηση server
   - `http.ListenAndServe(":8080", nil)`
   - Καταχώρηση routes (`/`, `/ascii-art`)
   - Log μηνύματα εκκίνησης

2. **GET `/`** — αρχική σελίδα
   - Φόρτωση και render του `templates/index.html`
   - Επιστροφή `200 OK`
   - Σε άγνωστο path → `404`

3. **POST `/ascii-art`** — λήψη φόρμας
   - Parse των form values (`text`, `banner`, `color`, `letters`)
   - Validation input (κενό text, μη έγκυρο banner, non-ASCII χαρακτήρες)
   - Κλήση `ascii.Generate(text, banner, color, letters)`
   - Render αποτελέσματος στο template
   - Σωστά HTTP status codes σε κάθε περίπτωση

4. **Error handling**
   - Custom σελίδες/μηνύματα για 400, 404, 500
   - Χρήση `http.Error()` ή template rendering

#### Παραδοτέα

- [ ] Server εκκινεί στο port 8080
- [ ] GET `/` επιστρέφει HTML φόρμα με `200`
- [ ] POST `/ascii-art` επιστρέφει αποτέλεσμα με `200`
- [ ] Λανθασμένο input → `400`
- [ ] Άγνωστο path → `404`
- [ ] Server error → `500`

---

### Άτομο 2 — ASCII Art Logic

**Αρχεία:** `ascii/ascii.go`, `banners/*.txt`

#### Καθήκοντα

1. **Κατέβασμα banner αρχείων**
   - `standard.txt`, `shadow.txt`, `thinkertoy.txt`
   - Πηγή: Zone01 repository / project subject
   - Επαλήθευση integrity (SHA256 checksums από το subject)

2. **Parsing banner αρχείου**
   - Κάθε banner αρχείο περιέχει ASCII χαρακτήρες από space (32) έως tilde (126)
   - Κάθε χαρακτήρας έχει ύψος **8 γραμμές**
   - Χαρακτήρες χωρίζονται με κενή γραμμή
   - Υλοποίηση `LoadBanner(filename string) (map[rune][]string, error)`

3. **Γεννήτρια ASCII art**
   - Υλοποίηση `Generate(text, banner, color, letters string) (string, error)`
   - Υποστήριξη `\n` (newline) στο input κείμενο
   - Για κάθε γραμμή κειμένου: συνένωση των 8 γραμμών κάθε χαρακτήρα οριζόντια
   - Επιστροφή error για non-printable ή non-ASCII χαρακτήρες
   - Αν `color` κενό → output χωρίς χρωματισμό
   - Αν `color` έχει τιμή, `letters` κενό → χρωματίζεται όλο το output
   - Αν και τα δύο έχουν τιμή → χρωματίζονται μόνο οι εμφανίσεις του `letters` substring
   - Χρωματισμός μέσω `<span style="color: X">` tags — **όχι** ANSI codes

4. **Unit tests**
   - Test για κάθε banner style
   - Test για `\n` στο input
   - Test για edge cases: κενό string, μόνο newlines, ειδικοί χαρακτήρες

#### Παραδοτέα

- [ ] `LoadBanner` φορτώνει σωστά κάθε banner
- [ ] `Generate` παράγει σωστό ASCII art για απλό κείμενο
- [ ] `Generate` χειρίζεται `\n` στο input
- [ ] `Generate` επιστρέφει error για μη έγκυρο input
- [ ] `Generate` χρωματίζει όλο το output όταν δίνεται μόνο χρώμα
- [ ] `Generate` χρωματίζει μόνο το substring όταν δίνονται και τα δύο
- [ ] Τουλάχιστον 5 unit tests που περνάνε

---

### Άτομο 3 — Frontend (Templates & Styling)

**Αρχεία:** `templates/index.html`

#### Καθήκοντα

1. **HTML φόρμα**
   - `<textarea>` για εισαγωγή κειμένου (name=`text`)
   - `<select>` για επιλογή banner style (name=`banner`) με options: standard, shadow, thinkertoy
   - `<select>` για επιλογή χρώματος (name=`color`) με options: black, red, yellow, blue, green, purple, orange, gray, pink, lightblue, lightgreen
   - `<input type="text">` για substring χρωματισμού (name=`letters`) — προαιρετικό
   - `<button type="submit">` για αποστολή
   - Method: `POST`, action: `/ascii-art`

2. **Εμφάνιση αποτελέσματος**
   - Εμφάνιση ASCII art output μέσα σε `<pre>` tag (για διατήρηση whitespace)
   - Το template δέχεται Go struct με πεδία: `Text`, `Banner`, `Color`, `Letters`, `Result`, `Error`
   - Αν υπάρχει `Error`: εμφάνιση μηνύματος σφάλματος
   - Αν υπάρχει `Result`: εμφάνιση ASCII art (περιέχει `<span>` tags — ΜΗΝ κάνεις escape)

3. **CSS Styling**
   - Καθαρό, minimal design (inline ή `<style>` tag — χωρίς εξωτερικά CDN)
   - Monospace font για το `<pre>` output
   - Responsive layout (κεντραρισμένη φόρμα)
   - Διαφορετικό styling για error state vs success state

4. **UX βελτιώσεις**
   - Το κείμενο που έγραψε ο χρήστης παραμένει στο textarea μετά το submit
   - Το επιλεγμένο banner style παραμένει επιλεγμένο
   - Το επιλεγμένο χρώμα παραμένει επιλεγμένο
   - Το letters input διατηρεί την τιμή του μετά το submit
   - Placeholder text στο textarea

#### Παραδοτέα

- [ ] Φόρμα λειτουργεί (POST στέλνεται σωστά)
- [ ] ASCII art εμφανίζεται σε `<pre>` με monospace font
- [ ] Χρωματισμένο ASCII art εμφανίζεται σωστά (τα `<span>` tags δεν escape-άρονται)
- [ ] Error state εμφανίζεται καθαρά
- [ ] Επιλογές φόρμας διατηρούνται μετά το submit (text, banner, color, letters)
- [ ] Χωρίς εξωτερικές εξαρτήσεις (no CDN, no JS frameworks)

---

## Go Struct για Templates (κοινό — ορίζει Άτομο 1)

```go
// handlers/handlers.go
type PageData struct {
    Text    string
    Banner  string
    Color   string        // χρώμα επιλογής (π.χ. red, blue, lightgreen)
    Letters string        // substring που θα χρωματιστεί — αν κενό, χρωματίζεται όλο
    Result  template.HTML // HTML output — περιέχει <span> tags για χρώματα
    Error   string
}
```

---

## Interface μεταξύ ατόμων

| Από | Προς | Interface |
|-----|------|-----------|
| Άτομο 1 | Άτομο 2 | `ascii.Generate(text, banner, color, letters string) (string, error)` |
| Άτομο 1 | Άτομο 3 | `PageData` struct + `templates/index.html` |
| Άτομο 2 | Άτομο 1 | Επιστρέφει `(string, error)` — output με `<span>` tags για χρώματα |

---

## Σειρά Εργασίας (Suggested Timeline)

```
Μέρα 1:  Άτομο 2 → κατεβάζει banners + υλοποιεί LoadBanner
          Άτομο 1 → στήνει skeleton server (main.go + routes)
          Άτομο 3 → φτιάχνει HTML φόρμα (στατική, χωρίς data)

Μέρα 2:  Άτομο 2 → υλοποιεί Generate + tests
          Άτομο 1 → POST handler + integration με ascii.Generate
          Άτομο 3 → Go template rendering + CSS

Μέρα 3:  Ολοκλήρωση → error handling, edge cases, testing ολόκληρης εφαρμογής
```

---

## Έλεγχος Ολοκλήρωσης (Definition of Done)

- [ ] `go run .` εκκινεί server χωρίς errors
- [ ] Όλα τα banner styles παράγουν σωστό output
- [ ] Σωστά HTTP status codes σε όλες τις περιπτώσεις
- [ ] `go vet ./...` — χωρίς warnings
- [ ] Χωρίς χρήση external packages (`go.mod` έχει μόνο `module` + `go` directive)
- [ ] Δοκιμή με: απλό κείμενο, κείμενο με `\n`, κενό input, μη-ASCII χαρακτήρες
- [ ] Δοκιμή χρωματισμού: χρώμα χωρίς letters, χρώμα με letters, χωρίς χρώμα
