# Ενημέρωση από Άτομο 1 — HTTP Server & Routing

## Τι έγινε

Ο server είναι έτοιμος και τρέχει στο port **8080**.

### Αρχεία που δημιουργήθηκαν

- `main.go` — εκκίνηση server, ορισμός routes
- `handlers/handlers.go` — λογική για GET `/` και POST `/ascii-art`
- `ascii/ascii.go` — **STUB** (προσωρινό, αντικαθίσταται από Άτομο 2)
- `templates/index.html` — **STUB** (προσωρινό, αντικαθίσταται από Άτομο 3)

### Τι δουλεύει ήδη

- GET `/` → επιστρέφει τη φόρμα με `200 OK`
- POST `/ascii-art` → παίρνει `text` + `banner` + `color` + `letters`, καλεί `ascii.Generate()`, εμφανίζει αποτέλεσμα
- Κενό input → `400 Bad Request`
- Άγνωστο path → `404 Not Found`
- Server error → `500 Internal Server Error`

---

## Για το Άτομο 2 — ASCII Logic

### Τι περιμένουμε

Να αντικαταστήσεις το αρχείο `ascii/ascii.go` με τον πραγματικό κώδικα.

### Σημείο αντικατάστασης

Το αρχείο `ascii/ascii.go` περιέχει αυτή τη στιγμή:

```go
package ascii

func Generate(text, banner, color, letters string) (string, error) {
    return "ASCII art coming soon...", nil
}
```

Η συνάρτηση **πρέπει να έχει ακριβώς αυτή την υπογραφή:**

```go
func Generate(text, banner, color, letters string) (string, error)
```

### Κανόνες που ΠΡΕΠΕΙ να ακολουθήσεις

- Το package πρέπει να ονομάζεται `ascii`
- Η συνάρτηση πρέπει να ονομάζεται `Generate` (κεφαλαίο G)
- Παράμετροι: `text string`, `banner string`, `color string`, `letters string` — με αυτή τη σειρά
- Επιστροφή: `(string, error)` — με αυτή τη σειρά
- Σε μη έγκυρο χαρακτήρα → επέστρεψε `error` (όχι panic)
- Σε μη έγκυρο banner → επέστρεψε `error`
- Αν `color` είναι κενό → παράγε ASCII art χωρίς χρωματισμό
- Αν `color` έχει τιμή αλλά `letters` είναι κενό → χρωμάτισε όλο το output
- Αν και τα δύο έχουν τιμή → χρωμάτισε μόνο τις εμφανίσεις του `letters` substring
- Το output πρέπει να περιέχει `<span style="color: X">` tags — όχι ANSI codes
- Πρέπει να δημιουργήσεις τον φάκελο `banners/` με τα αρχεία: `standard.txt`, `shadow.txt`, `thinkertoy.txt`

---

## Για το Άτομο 3 — Frontend

### Τι περιμένουμε

Να αντικαταστήσεις το αρχείο `templates/index.html` με το πλήρες styled template.

### Σημείο αντικατάστασης

Το αρχείο `templates/index.html` είναι αυτή τη στιγμή ένα minimal placeholder.

### Κανόνες που ΠΡΕΠΕΙ να ακολουθήσεις

Το template δέχεται ένα Go struct με αυτά τα πεδία:

```go
type PageData struct {
    Text    string        // το κείμενο που έγραψε ο χρήστης
    Banner  string        // το banner style που διάλεξε
    Color   string        // χρώμα (π.χ. red, #ff0000, rgb(255,0,0))
    Letters string        // substring που θα χρωματιστεί — αν κενό, χρωματίζεται όλο
    Result  template.HTML // το ASCII art αποτέλεσμα — περιέχει <span> tags, ΜΗΝ το escape-άρεις
    Error   string        // μήνυμα λάθους (αν υπάρχει)
}
```

Στο HTML χρησιμοποιείς τα πεδία ως εξής:

| Πεδίο | Χρήση στο template |
|-------|-------------------|
| `{{.Text}}` | value του textarea (για να παραμένει μετά το submit) |
| `{{.Banner}}` | για να παραμένει επιλεγμένο το σωστό option |
| `{{.Color}}` | value του color input |
| `{{.Letters}}` | value του letters input |
| `{{.Result}}` | μέσα σε `<pre>` tag — **χρησιμοποίησε `{{.Result}}` όχι `{{html .Result}}`** |
| `{{.Error}}` | εμφάνιση μηνύματος λάθους |

**Απαραίτητα στοιχεία της φόρμας:**
- `<form method="POST" action="/ascii-art">` — ακριβώς έτσι
- `<textarea name="text">` — το name πρέπει να είναι `text`
- `<select name="banner">` — options: `standard`, `shadow`, `thinkertoy` (lowercase, ακριβώς έτσι)
- `<select name="color">` — dropdown με συγκεκριμένες τιμές (δες παρακάτω)
- `<input type="text" name="letters">` — για το substring (προαιρετικό)
- `<pre>` tag για το αποτέλεσμα — απαραίτητο για το whitespace

**Επιτρεπτές τιμές για το `color` select (ακριβώς αυτά τα values):**
`black` (default), `red`, `yellow`, `blue`, `green`, `purple`, `orange`, `gray`, `pink`, `lightblue`, `lightgreen`

Το selected option πρέπει να παραμένει μετά το submit — χρησιμοποίησε `{{if eq .Color "value"}}selected{{end}}` σε κάθε option.

Ίδια λογική για το banner select — χρησιμοποίησε `{{if eq .Banner "value"}}selected{{end}}` σε κάθε option.

**Χωρίς εξωτερικές εξαρτήσεις** — no CDN, no JS frameworks.

---

## Οδηγίες Git / Gitea για όλους

**Repo URL:** `https://platform.zone01.gr/git/ivogiake/ascii-art-web.git`

### Πώς να ξεκινήσετε (clone)

```bash
git clone https://platform.zone01.gr/git/ivogiake/ascii-art-web.git
cd ascii-art-web
```

---

### Άτομο 2 — Git workflow

**1. Δημιουργία branch:**
```bash
git checkout -b feature/ascii-logic
```

**2. Δουλεύεις στο αρχείο:**
- `ascii/ascii.go` — αντικαθιστάς το stub με τον πραγματικό κώδικα
- `banners/standard.txt`, `banners/shadow.txt`, `banners/thinkertoy.txt` — προσθέτεις τα banner αρχεία

**3. Όταν τελειώσεις κάποιο κομμάτι, commit:**
```bash
git add ascii/ascii.go
git commit -m "feat: implement ascii Generate function"

git add banners/
git commit -m "feat: add banner files (standard, shadow, thinkertoy)"
```

**4. Push στο Gitea:**
```bash
git push origin feature/ascii-logic
```

**5. Όταν είσαι έτοιμος**, κάνε merge στο `main`:
```bash
git checkout main
git merge feature/ascii-logic
git push origin main
```

Ειδοποίησε το Άτομο 1 ότι έκανες merge ώστε να τραβήξει τις αλλαγές.

**⚠️ Προσοχή:**
- Μην αγγίξεις τα `main.go`, `handlers/handlers.go`, `templates/`
- Μην αλλάξεις την υπογραφή της `Generate()` χωρίς να ενημερώσεις το Άτομο 1

---

### Άτομο 3 — Git workflow

**1. Δημιουργία branch:**
```bash
git checkout -b feature/frontend
```

**2. Δουλεύεις στο αρχείο:**
- `templates/index.html` — αντικαθιστάς το placeholder με το πλήρες styled template

**3. Όταν τελειώσεις, commit:**
```bash
git add templates/index.html
git commit -m "feat: implement styled HTML template"
```

**4. Push στο Gitea:**
```bash
git push origin feature/frontend
```

**5. Όταν είσαι έτοιμος**, κάνε merge στο `main`:
```bash
git checkout main
git merge feature/frontend
git push origin main
```

Ειδοποίησε το Άτομο 1 ότι έκανες merge ώστε να τραβήξει τις αλλαγές.

**⚠️ Προσοχή:**
- Μην αγγίξεις τα `main.go`, `handlers/handlers.go`, `ascii/`
- Τα `name` attributes της φόρμας πρέπει να είναι ακριβώς: `text`, `banner`, `color`, `letters`
- Τα Go template tags (`{{.Text}}`, `{{.Result}}` κλπ.) πρέπει να είναι ακριβώς έτσι

---

### Άτομο 1 — Sync workflow

Όταν τα Άτομα 2 και 3 κάνουν merge στο `main`, τράβα τις αλλαγές τοπικά:

**1. Τράβα τις αλλαγές:**
```bash
git checkout main
git pull origin main
```

**2. Έλεγξε** ότι το `ascii/ascii.go` αντικαταστάθηκε από το Άτομο 2 και το `templates/index.html` από το Άτομο 3.

**3. Τρέξε τελικό test:**
```bash
go run .
go vet ./...
```
