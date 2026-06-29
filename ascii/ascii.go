package ascii

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generate is the main entry point for this package.
// It reads the requested banner file, splits the input text into lines,
// and converts each line into 8 rows of ASCII art glyphs.
// The result is returned as a single string ready to be placed inside an
// HTML <pre> tag. When color is requested, the relevant glyphs are wrapped
// in <span style="color: ..."> tags so the browser renders the color.
//
// Parameters:
//   - text:    the user's input; may contain real newlines (\n) from a textarea
//   - banner:  which banner file to use — "standard", "shadow", or "thinkertoy"
//   - color:   a CSS color value (e.g. "red", "#ff0000"); empty means no coloring
//   - letters: a substring of text whose ASCII glyphs should be colored;
//     empty means color the entire output (when color is also set)
//
// Returns an error if the banner file cannot be read, the file looks corrupt,
// or the text contains characters outside the printable ASCII range (32–126).
func Generate(text, banner, color, letters string) (string, error) {
	// Build the path to the banner file relative to the working directory.
	// The server is always started from the project root, so "banners/name.txt"
	// resolves correctly regardless of the OS path separator.
	filename := filepath.Join("banners", banner+".txt")

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not open banner file: %w", err)
	}

	// A valid banner file always has exactly 855 newline characters
	// (95 characters × 9 lines each, where line 0 is blank = 855 \n total).
	// This catches truncated or wrong files before we do any indexing.
	if strings.Count(string(data), "\n") != 855 {
		return "", fmt.Errorf("banner file is corrupt or invalid")
	}

	// Split the file into individual lines. We normalise \r\n first so that
	// banner files saved on Windows don't produce lines ending in \r, which
	// would appear as garbage characters in the ASCII art output.
	bannerLines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	// Some browsers send \r\n for textarea newlines (notably on Windows).
	// Normalise to \n before splitting so we always get clean line strings.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var result strings.Builder

	for i, line := range lines {
		// An empty string here means the user pressed Enter at this position.
		// We output one blank line as a separator, but skip the very first
		// element to avoid a leading blank line when the input starts with \n.
		if line == "" {
			if i > 0 {
				result.WriteByte('\n')
			}
			continue
		}

		// Reject any character outside the printable ASCII range.
		// The banner files only contain glyphs for codes 32–126, so anything
		// outside that range would cause an out-of-bounds index.
		for _, r := range line {
			if r < 32 || r > 126 {
				return "", fmt.Errorf("character %q is out of supported range", r)
			}
		}

		// Each line of text becomes 8 rows of ASCII art.
		// We build one row at a time (row 1 through 8) by picking the
		// corresponding glyph line for every character in the input.
		for row := 1; row <= 8; row++ {
			result.WriteString(buildRow(line, bannerLines, color, letters, row))
			result.WriteByte('\n')
		}
	}

	return result.String(), nil
}

// buildRow assembles a single horizontal row of ASCII art for one line of text.
//
// How it works: for each character in the line it looks up the correct glyph
// line inside the banner using the formula:
//
//	index = (charCode - 32) * 9 + row
//
// where charCode is the ASCII value of the character (space = 32 is index 0),
// the *9 skips over the 8 glyph lines + 1 blank separator per character,
// and +row selects which of the 8 glyph lines we want (1 = top, 8 = bottom).
//
// When coloring is requested, the function wraps the relevant glyph segments
// in <span> tags. Opening and closing tags are inserted only when the "colored"
// state changes between consecutive characters, so adjacent colored characters
// share a single <span> rather than each getting their own.
func buildRow(line string, bannerLines []string, color, letters string, row int) string {
	// Fast path: no color requested — just concatenate glyph lines directly.
	if color == "" {
		var sb strings.Builder
		for _, r := range line {
			sb.WriteString(bannerLines[(int(r)-32)*9+row])
		}
		return sb.String()
	}

	// Build the opening tag once; it is reused for every span in this row.
	colorTag := `<span style="color: ` + color + `">`

	// Color everything: wrap the entire row in a single span.
	if letters == "" {
		var sb strings.Builder
		sb.WriteString(colorTag)
		for _, r := range line {
			sb.WriteString(bannerLines[(int(r)-32)*9+row])
		}
		sb.WriteString("</span>")
		return sb.String()
	}

	// Partial coloring: find which character positions belong to a match of
	// `letters` in the line, then wrap only those glyph segments in spans.
	runes := []rune(line)
	colored := coloredPositions(runes, []rune(letters))

	var sb strings.Builder
	inSpan := false // tracks whether a <span> is currently open

	for i, r := range runes {
		// Toggle the span open/closed whenever the colored state changes.
		if colored[i] != inSpan {
			if colored[i] {
				sb.WriteString(colorTag)
			} else {
				sb.WriteString("</span>")
			}
			inSpan = colored[i]
		}
		sb.WriteString(bannerLines[(int(r)-32)*9+row])
	}

	// Close any span that is still open after the last character.
	if inSpan {
		sb.WriteString("</span>")
	}
	return sb.String()
}

// coloredPositions returns a boolean slice the same length as chars.
// A true value at index i means character i is part of at least one occurrence
// of letters inside chars. Overlapping matches are both marked true.
//
// The search is a straightforward O(n*m) scan: for each starting position i,
// check whether chars[i:i+n] matches letters character by character.
// This is fine for the short strings expected from a web form.
func coloredPositions(chars, letters []rune) []bool {
	colored := make([]bool, len(chars))
	n := len(letters)

	// The outer loop only reaches positions where a full match can still fit.
	// If len(chars) < n, len(chars)-n is negative and the loop never runs.
	for i := 0; i <= len(chars)-n; i++ {
		match := true
		for j := range n {
			if chars[i+j] != letters[j] {
				match = false
				break
			}
		}
		if match {
			// Mark every character position covered by this match.
			for j := range n {
				colored[i+j] = true
			}
		}
	}
	return colored
}
