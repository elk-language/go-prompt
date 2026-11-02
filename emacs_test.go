package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
)

func TestEmacsKeyBindings(t *testing.T) {
	p := New(NoopExecutor)
	buf := p.buffer
	buf.InsertTextMoveCursor("abcde", DefColCount, DefRowCount, false)
	if buf.cursorPosition != istrings.RuneNumber(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}

	// Go to the beginning of the line
	applyEmacsKeyBind(p, ControlA)
	if buf.cursorPosition != 0 {
		t.Errorf("Want %d, but got %d", 0, buf.cursorPosition)
	}

	// Go to the end of the line
	applyEmacsKeyBind(p, ControlE)
	if buf.cursorPosition != istrings.RuneNumber(len("abcde")) {
		t.Errorf("Want %d, but got %d", len("abcde"), buf.cursorPosition)
	}
}

func TestEmacsSwapCharactersBeforeCursor(t *testing.T) {
	p := New(NoopExecutor)
	buf := p.buffer
	buf.InsertTextMoveCursor("あいう", DefColCount, DefRowCount, false)

	applyEmacsKeyBind(p, ControlT)
	if buf.Text() != "あうい" {
		t.Errorf("text %#v, want %#v", buf.Text(), "あうい")
	}
	if buf.cursorPosition != 3 {
		t.Errorf("cursor %#v, want 3", buf.cursorPosition)
	}

	buf.cursorPosition = 1
	applyEmacsKeyBind(p, ControlT)
	if buf.Text() != "あうい" {
		t.Errorf("text %#v, want %#v", buf.Text(), "あうい")
	}
}

func TestEmacsSwapWordsBeforeCursor(t *testing.T) {
	p := New(NoopExecutor)
	buf := p.buffer
	buf.InsertTextMoveCursor("foo, bar baz", DefColCount, DefRowCount, false)
	buf.cursorPosition = istrings.RuneCountInString("foo, bar")

	applyEmacsKeyBind(p, EscapeT)
	if buf.Text() != "bar, foo baz" {
		t.Errorf("text %#v, want %#v", buf.Text(), "bar, foo baz")
	}
	if buf.cursorPosition != istrings.RuneCountInString("bar, foo") {
		t.Errorf("cursor %#v, want %d", buf.cursorPosition, istrings.RuneCountInString("bar, foo"))
	}

	buf.cursorPosition = istrings.RuneCountInString("bar, foo baz")
	applyEmacsKeyBind(p, EscapeT)
	if buf.Text() != "bar, baz foo" {
		t.Errorf("text %#v, want %#v", buf.Text(), "bar, baz foo")
	}

	buf.cursorPosition = istrings.RuneCountInString("bar,")
	applyEmacsKeyBind(p, EscapeT)
	if buf.Text() != "bar, baz foo" {
		t.Errorf("single word should stay %#v, got %#v", "bar, baz foo", buf.Text())
	}
}

func TestEscapeThenTSwapsWords(t *testing.T) {
	p := New(NoopExecutor)
	p.renderer.col = DefColCount
	p.renderer.row = DefRowCount
	buf := p.buffer
	buf.InsertTextMoveCursor("foo, bar baz", DefColCount, DefRowCount, false)
	buf.cursorPosition = istrings.RuneCountInString("foo, bar")

	// Escape and t are delivered by separate reads.
	p.feed([]byte{0x1b})
	if buf.Text() != "foo, bar baz" {
		t.Fatalf("text after Escape: %#v", buf.Text())
	}

	p.feed([]byte{'T'})
	if buf.Text() != "bar, foo baz" {
		t.Fatalf("text after Escape then T: %#v", buf.Text())
	}

	// The same chord delivered in one read.
	buf.cursorPosition = istrings.RuneCountInString(buf.Text())
	p.feed([]byte{0x1b, 't'})
	if buf.Text() != "bar, baz foo" {
		t.Fatalf("text after combined Escape t: %#v", buf.Text())
	}
}

func TestEmacsWordCaseAndDelete(t *testing.T) {
	p := newEmacsPrompt()
	buf := p.buffer
	buf.InsertTextMoveCursor("say hello-WORLD", DefColCount, DefRowCount, false)
	buf.cursorPosition = istrings.RuneCountInString("say ")

	applyEmacsKeyBind(p, EscapeU)
	if buf.Text() != "say HELLO-WORLD" {
		t.Fatalf("upcase %#v", buf.Text())
	}
	if buf.cursorPosition != istrings.RuneCountInString("say HELLO") {
		t.Fatalf("upcase cursor %#v", buf.cursorPosition)
	}

	buf.cursorPosition = istrings.RuneCountInString("say HELLO-")
	applyEmacsKeyBind(p, EscapeL)
	if buf.Text() != "say HELLO-world" {
		t.Fatalf("downcase %#v", buf.Text())
	}

	buf.cursorPosition = istrings.RuneCountInString("say ")
	applyEmacsKeyBind(p, EscapeC)
	if buf.Text() != "say Hello-world" {
		t.Fatalf("capitalize %#v", buf.Text())
	}
	if buf.cursorPosition != istrings.RuneCountInString("say Hello") {
		t.Fatalf("capitalize cursor %#v", buf.cursorPosition)
	}
}

func TestEmacsDeleteWordAfterCursor(t *testing.T) {
	p := newEmacsPrompt()
	buf := p.buffer
	buf.InsertTextMoveCursor("foo bar", DefColCount, DefRowCount, false)
	buf.cursorPosition = 0

	applyEmacsKeyBind(p, EscapeD)
	if buf.Text() != " bar" {
		t.Fatalf("text %#v, want %#v", buf.Text(), " bar")
	}
	if buf.cursorPosition != 0 {
		t.Fatalf("cursor %#v, want 0", buf.cursorPosition)
	}
}

func TestEmacsHistoryEnds(t *testing.T) {
	p := newEmacsPrompt()
	p.history.Add("one")
	p.history.Add("two")
	p.buffer.InsertTextMoveCursor("draft", DefColCount, DefRowCount, false)

	applyEmacsKeyBind(p, EscapeLess)
	if p.buffer.Text() != "one" {
		t.Fatalf("oldest %#v", p.buffer.Text())
	}

	applyEmacsKeyBind(p, EscapeGreater)
	if p.buffer.Text() != "draft" {
		t.Fatalf("newest %#v", p.buffer.Text())
	}
}

func TestEscapeWordKeysFromInput(t *testing.T) {
	p := newEmacsPrompt()
	p.buffer.InsertTextMoveCursor("hello", DefColCount, DefRowCount, false)
	p.buffer.cursorPosition = 0

	p.feed([]byte{0x1b})
	p.feed([]byte{'u'})
	if p.buffer.Text() != "HELLO" {
		t.Fatalf("separate reads %#v", p.buffer.Text())
	}

	p.buffer.cursorPosition = 0
	p.feed([]byte{0x1b, 'l'})
	if p.buffer.Text() != "hello" {
		t.Fatalf("combined read %#v", p.buffer.Text())
	}
}

func newEmacsPrompt() *Prompt {
	p := New(NoopExecutor)
	p.renderer.col = DefColCount
	p.renderer.row = DefRowCount
	return p
}

func applyEmacsKeyBind(p *Prompt, key Key) {
	for i := range emacsKeyBindings {
		kb := emacsKeyBindings[i]
		if kb.Key == key {
			kb.Fn(p)
		}
	}
}
