# Update from Person 2 — ASCII Art Logic

## What was done

The ASCII art generation logic is complete and lives in `ascii/ascii.go`.

### Files created / modified

- `ascii/ascii.go` — replaced the Person 1 stub with the real implementation
- `banners/standard.txt` — banner file (95 characters × 9 lines = 856 lines total)
- `banners/shadow.txt` — same structure
- `banners/thinkertoy.txt` — same structure

---

## How the code works

### Banner file structure

Each banner file contains the printable ASCII characters from space (32) to tilde (126).
Every character occupies exactly **8 lines** and is separated from the next character by **1 blank line**.
Total: 95 characters × 9 lines = 855 newlines (856 lines), starting with a blank line.

```
[blank line]              ← line 0
[8 lines for ' ']         ← lines 1–8   (ASCII 32)
[blank line]              ← line 9
[8 lines for '!']         ← lines 10–17 (ASCII 33)
...
```

### Index formula

To find the glyph data for character `r` at row `row` (1–8):

```
bannerLines[(int(r) - 32) * 9 + row]
```

This single formula is the entire lookup mechanism — no map or struct needed.

### Coloring

The output is returned as a plain string containing `<span style="color: X">` HTML tags (not ANSI escape codes), so that colors render correctly in the browser inside a `<pre>` tag.

| `color` | `letters` | Behaviour |
|---------|-----------|-----------|
| empty | — | plain ASCII art, no tags |
| e.g. `red` | empty | the entire output is wrapped in a `<span>` |
| e.g. `red` | e.g. `ell` | only occurrences of `"ell"` are wrapped in `<span>` |

---

## How this connects to the other parts

### Connection to Person 1 (handlers.go)

`handlers/handlers.go` calls:

```go
result, err := ascii.Generate(text, banner, color, letters)
```

- `text` — the text from the textarea; may contain `\n` (Enter key)
- `banner` — `"standard"`, `"shadow"`, or `"thinkertoy"`; already validated by the handler before the call
- `color` — a CSS color value (e.g. `"red"`, `"blue"`); empty string if the user did not pick one
- `letters` — the substring to color; empty string if the user left the field blank

On success: `Generate` returns `(string, nil)` — the handler stores the string in `data.Result` as `template.HTML`.  
On error: `Generate` returns `("", error)` — the handler displays the error message and responds with HTTP 400.

**Note:** `Generate` does not re-validate the banner name (the handler already does that). If a banner file is missing from disk, `Generate` returns an error.

### Connection to Person 3 (index.html)

The result of `Generate` is displayed in the template as:

```html
<pre>{{.Result}}</pre>
```

`.Result` is of type `template.HTML` (not a plain `string`), which tells Go's template engine **not** to escape the content — so the `<span>` tags are rendered as real HTML and the browser applies the colors.

If Person 3 writes `{{html .Result}}` or `{{.Result | html}}` instead, the `<span>` tags will appear as literal text in the page and no colors will show.

---

## Banner file location

The `banners/*.txt` files must live in the **project root** (next to `main.go`), not inside `ascii/`. The server is started with `go run .` from the root, so the relative path `"banners/standard.txt"` resolves there.

```
ascii-art-web/          ← go run . is run from here
├── banners/            ← HERE (not inside ascii/)
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
├── ascii/
│   └── ascii.go
...
```

---

## Git workflow

**1. Create your branch:**
```bash
git checkout -b feature/ascii-logic
```

**2. Files you touch:**
- `ascii/ascii.go`
- `banners/standard.txt`, `banners/shadow.txt`, `banners/thinkertoy.txt`

**3. Commit as you go:**
```bash
git add ascii/ascii.go
git commit -m "feat: implement ascii Generate function"

git add banners/
git commit -m "feat: add banner files (standard, shadow, thinkertoy)"
```

**4. Push and merge:**
```bash
git push origin feature/ascii-logic
git checkout main
git merge feature/ascii-logic
git push origin main
```

Let Person 1 know when you have merged so they can run `git pull` and replace their stub.

**⚠️ Do not touch** `main.go`, `handlers/handlers.go`, or `templates/`.  
**Do not change** the signature of `Generate()` without notifying Person 1.
