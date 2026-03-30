package rendering

import (
	"github.com/gdamore/tcell/v3"
	"github.com/mattn/go-runewidth"
)

type TextInput struct {
	X, Y      int    // Screen position of input field
	CursorPos int    // Position in runes (0 = before first char)
	Buffer    []rune // The text content
}

func (t *TextInput) HandleKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		// Insert character at cursor position
		t.Buffer = append(t.Buffer[:t.CursorPos],
			append([]rune(ev.Str()), t.Buffer[t.CursorPos:]...)...)
		t.CursorPos++

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if t.CursorPos > 0 {
			t.Buffer = append(t.Buffer[:t.CursorPos-1], t.Buffer[t.CursorPos:]...)
			t.CursorPos--
		}

	case tcell.KeyDelete:
		if t.CursorPos < len(t.Buffer) {
			t.Buffer = append(t.Buffer[:t.CursorPos], t.Buffer[t.CursorPos+1:]...)
		}

	case tcell.KeyLeft:
		if t.CursorPos > 0 {
			t.CursorPos--
		}

	case tcell.KeyRight:
		if t.CursorPos < len(t.Buffer) {
			t.CursorPos++
		}

	case tcell.KeyHome:
		t.CursorPos = 0

	case tcell.KeyEnd:
		t.CursorPos = len(t.Buffer)
	}
}

func (t *TextInput) Draw(s tcell.Screen) {
	// Draw the text
	text := string(t.Buffer)
	s.PutStr(t.X, t.Y, text)

	// Calculate cursor screen position (accounting for Unicode width!)
	cursorScreenX := t.X + t.runeWidth(t.Buffer[:t.CursorPos])
	s.ShowCursor(cursorScreenX, t.Y)
}

// Calculate display width of runes (CJK/emoji = 2 cells)
func (t *TextInput) runeWidth(runes []rune) int {
	w := 0
	for _, r := range runes {
		w += runewidth.RuneWidth(r) // github.com/mattn/go-runewidth
	}
	return w
}
