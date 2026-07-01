// Package logo renders a Glash wordmark in a stylized way.
package logo

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"glash/internal/ui/styles"
	"github.com/charmbracelet/x/ansi"
)

// blockLetters represents the block font letterforms for "GLASH".
var blockLetters = map[string]string{
	"G": `██████  
██      
██ ████ 
██   ██ 
 ██████ `,
	"L": `██      
██      
██      
██      
███████ `,
	"A": `██████  
██   ██ 
███████ 
██   ██ 
███████ `,
	"S": `██████  
      ██
███████ 
██      
███████ `,
	"H": `██   ██ 
██   ██ 
███████ 
██   ██ 
██   ██ `,
}

// renderBlockWord renders the word using block font letters.
func renderBlockWord(word string) string {
	var lines []string
	for lineIdx := 0; lineIdx < 5; lineIdx++ {
		var lineBuilder strings.Builder
		for i := 0; i < len(word); i++ {
			r := rune(word[i])
			if r == 'G' || r == 'g' {
				lineBuilder.WriteString(getBlockLine("G", lineIdx))
			} else if r == 'L' || r == 'l' {
				lineBuilder.WriteString(getBlockLine("L", lineIdx))
			} else if r == 'A' || r == 'a' {
				lineBuilder.WriteString(getBlockLine("A", lineIdx))
			} else if r == 'S' || r == 's' {
				lineBuilder.WriteString(getBlockLine("S", lineIdx))
			} else if r == 'H' || r == 'h' {
				lineBuilder.WriteString(getBlockLine("H", lineIdx))
			}
		}
		lines = append(lines, lineBuilder.String())
	}
	return strings.Join(lines, "\n")
}

func getBlockLine(letter string, lineIdx int) string {
	letterStr, ok := blockLetters[letter]
	if !ok {
		return strings.Repeat(" ", 8)
	}
	lines := strings.Split(letterStr, "\n")
	if lineIdx < 0 || lineIdx >= len(lines) {
		return strings.Repeat(" ", 8)
	}
	return lines[lineIdx]
}

const diag = `╱`

// Opts are the options for rendering the Glash title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	SuperColor   color.Color // Super™ text color
	VersionColor color.Color // version text color
	Width        int         // width of the rendered logo, used for truncation
	Hyper        bool        // whether it is Glash or Hyperglash

	// When true, stretch a random letterform on each render. Has no effect in
	// compact mode. Mainly for testing. In production you will want to cache
	// the stretched letterform to keep the logo from jittering on resize.
	Unstable bool
}

// Render renders the Glash logo. Set the argument to true to render the narrow
// version, intended for use in a sidebar.
//
// The compact argument determines whether it renders compact for the sidebar
// or wider for the main pane.
func Render(base lipgloss.Style, version string, compact bool, o Opts) string {
	super := "Super™"
	if !o.Hyper {
		super = " " + super
	}

	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title using block font.
	glashWord := "GLASH"
	if o.Hyper {
		glashWord = "HYPERGLASH"
	}
	
	glash := renderBlockWord(glashWord)
	
	glashWidth := lipgloss.Width(glash)
	b := new(strings.Builder)
	for r := range strings.SplitSeq(glash, "\n") {
		fmt.Fprintln(b, styles.ApplyForegroundGrad(base, r, o.TitleColorA, o.TitleColorB))
	}
	glash = b.String()

	// Super and version.
	metaRowGap := 1
	maxVersionWidth := glashWidth - lipgloss.Width(super) - metaRowGap
	version = ansi.Truncate(version, maxVersionWidth, "…") // truncate version if too long.
	if o.Hyper && compact {
		version += " "
	}
	gap := max(0, glashWidth-lipgloss.Width(super)-lipgloss.Width(version))
	metaRow := fg(o.SuperColor, super) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)

	// Join the meta row and big Glash title.
	glash = strings.TrimSpace(metaRow + "\n" + glash)

	// Narrow version. If this is Hyperglash, this is also a stacked version.
	if compact {
		field := fg(o.FieldColor, strings.Repeat(diag, glashWidth))
		return strings.Join([]string{field, field, glash, field, ""}, "\n")
	}

	fieldHeight := lipgloss.Height(glash)

	// Left field.
	const leftWidth = 6
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightWidth := max(15, o.Width-glashWidth-leftWidth-2) // 2 for the gap.
	const stepDownAt = 0
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		width := rightWidth
		if i >= stepDownAt {
			width = rightWidth - (i - stepDownAt)
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diag, width)), "\n")
	}

	// Return the wide version.
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, glash, hGap, rightField.String())
	if o.Width > 0 {
		// Truncate the logo to the specified width.
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

// SmallRender renders a smaller version of the Glash logo, suitable for
// smaller windows or sidebar usage.
func SmallRender(t *styles.Styles, width int, o Opts) string {
	name := "Glash"
	if o.Hyper {
		name = "HYPERGLASH"
	}
	super := "Super™"
	if !o.Hyper {
		super = " " + super
	}
	title := t.Logo.SmallSuper.Render(super)
	title = fmt.Sprintf("%s %s", title, styles.ApplyBoldForegroundGrad(t.Logo.GradCanvas, name, t.Logo.SmallGradFromColor, t.Logo.SmallGradToColor))
	remainingWidth := width - lipgloss.Width(title) - 1 // 1 for the space after the name
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, t.Logo.SmallDiagonals.Render(lines))
	}
	return title
}
