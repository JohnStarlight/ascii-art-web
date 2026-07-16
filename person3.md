# Update from Person 3 — Frontend UI & Error Pages

## What was done

- Built the HTML structure for `templates/index.html` with Flexbox layout
- Added JavaScript to dynamically change the preview background color based on selected color
- Created 4 error pages: `templates/400.html`, `templates/404.html`, `templates/405.html`, `templates/500.html`
- Implemented `statusTemplate()` function in `handlers/handlers.go` that serves the correct error page based on HTTP status code
- 🌸 Easter egg (no further comments)

## Files created / modified

- `templates/index.html` — main page with form and preview area
- `templates/400.html` — Bad Request page
- `templates/404.html` — Not Found page
- `templates/405.html` — Method Not Allowed page
- `templates/500.html` — Internal Server Error page
- `handlers/handlers.go` — added `statusTemplate()` function

---

## How the code works

### index.html layout

The page has two sections inside `#main-content` (Flexbox, column direction):

- `#controls` — the form: `<form method="POST" action="/ascii-art">` with the name attributes **exactly** as the handler expects them: `text`, `banner`, `color`, `letters`
- `#preview` — displays either `{{.Error}}` or `{{.Result}}` inside a `<pre>` tag

The `<pre>` uses `white-space: pre-wrap` so the ASCII art keeps its whitespace but doesn't break the layout on long text.

### Keeping values after submit

- The textarea keeps the text: `<textarea name="text">{{.Text}}</textarea>`
- The letters input keeps its value: `value="{{.Letters}}"`
- The color select keeps the selection via `{{if eq .Color "value"}}selected{{end}}` on every option

**Known TODO:** the banner select does **not** yet have `{{if eq .Banner "value"}}selected{{end}}` — after submit it resets to Standard. Same pattern as the color select if someone adds it.

### Allowed color values

`black` (default), `red`, `yellow`, `blue`, `green`, `purple`, `orange`, `gray`, `pink`, `lightblue`, `lightgreen` — exactly the values defined by Person 1.

### Error pages (400 / 404 / 405 / 500)

Each status code has its own template in `templates/`. They are simple static pages — a title, a short message, and a link back to the home page (`<a href="/">`). They take no data (executed with `nil`).

### The statusTemplate function (handlers.go)

```go
func statusTemplate(w http.ResponseWriter, status int)
```

- Switches on the status code and loads the matching template: `400.html`, `404.html`, `405.html`, `500.html`
- For any other status (default case), or if `ParseFiles` fails, it falls back to a plain `http.Error` with 500
- Writes `w.WriteHeader(status)` **before** executing the template, so the client receives the correct HTTP status code (not 200)

It is already used in two places:

- `Home` → unknown path → `statusTemplate(w, http.StatusNotFound)`, wrong method → `statusTemplate(w, http.StatusMethodNotAllowed)`
- `AsciiArt` → wrong method → `statusTemplate(w, http.StatusMethodNotAllowed)`

**Note:** form validation errors (empty text, invalid banner) do **not** go through `statusTemplate` — they are shown inside index.html via `{{.Error}}` with status 400 (through `renderTemplate`), so the user doesn't lose the form.

---

## How this connects to the other parts

### Connection to Person 1 (handlers.go)

The template receives the `PageData` struct as defined in person1.md:

| Field | Where it is used |
|-------|------------------|
| `{{.Text}}` | textarea content |
| `{{.Color}}` | selected option in the color select |
| `{{.Letters}}` | value of the letters input |
| `{{.Result}}` | inside the `<pre>` tag |
| `{{.Error}}` | error message in `#preview` |

### Connection to Person 2 (ascii.go)

`{{.Result}}` is rendered **without** escaping (`{{.Result}}`, not `{{html .Result}}`), because the field is `template.HTML` and contains the `<span style="color: X">` tags produced by `ascii.Generate`. This way the browser renders the colors correctly inside the `<pre>`.

---

## Git workflow

**1. Create your branch:**
```bash
git checkout -b feature/frontend
```

**2. Files you touch:**
- `templates/index.html`
- `templates/400.html`, `templates/404.html`, `templates/405.html`, `templates/500.html`
- `handlers/handlers.go` (only the `statusTemplate()` addition)

**3. Commit as you go:**
```bash
git add templates/index.html
git commit -m "feat: implement styled HTML template"

git add templates/400.html templates/404.html templates/405.html templates/500.html handlers/handlers.go
git commit -m "feat: add error pages and statusTemplate handler"
```

**4. Push and merge:**
```bash
git push origin feature/frontend
git checkout main
git merge feature/frontend
git push origin main
```

Let Person 1 know when you have merged so they can run `git pull` and do the final test (`go run .`, `go vet ./...`).

**⚠️ Do not touch** `main.go` or `ascii/`.
**Do not change** the form name attributes (`text`, `banner`, `color`, `letters`) or the Go template tags without notifying Person 1.
