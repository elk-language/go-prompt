package prompt

import (
	"testing"

	istrings "github.com/elk-language/go-prompt/strings"
)

// quotedSuggestion is a completion candidate that contains the completion word
// separator (a space). Committing such a suggestion is exactly the case that
// previously got duplicated on the command line.
const quotedSuggestion = `"file with space.txt"`

// newCompletionTestPrompt builds a Prompt wired with a completer that always
// offers a single suggestion replacing the range [start, cursor], and with a
// renderer sized like a normal terminal so buffer edits behave normally.
func newCompletionTestPrompt(suggestion string, start istrings.RuneNumber) *Prompt {
	p := New(
		NoopExecutor,
		WithPrefix(""),
		WithCompleter(func(d Document) ([]Suggest, istrings.RuneNumber, istrings.RuneNumber) {
			return []Suggest{{Text: suggestion}}, start, d.CurrentRuneIndex()
		}),
	)
	p.renderer.col = DefColCount
	p.renderer.row = DefRowCount

	return p
}

// TestCompletionCommitDoesNotDuplicateQuotedSuggestion reproduces the reported
// bug: selecting a quoted path with spaces via Tab (which writes the inline
// preview into the buffer) and then committing it must leave the line intact
// instead of duplicating the suggestion.
func TestCompletionCommitDoesNotDuplicateQuotedSuggestion(t *testing.T) {
	p := newCompletionTestPrompt(quotedSuggestion, 3)

	p.buffer.InsertTextMoveCursor("rm ", DefColCount, DefRowCount, false)
	p.completion.Update(*p.buffer.Document())

	// Tab selects the suggestion; the inline preview inserts it into the buffer.
	p.handleCompletionKeyBinding(nil, Tab, p.completion.Completing())

	want := "rm " + quotedSuggestion
	if got := p.buffer.Text(); got != want {
		t.Fatalf("after Tab preview: buffer = %q, want %q", got, want)
	}

	// Enter commits the selection and must not re-insert (duplicate) the text.
	p.handleCompletionKeyBinding(nil, Enter, p.completion.Completing())

	if got := p.buffer.Text(); got != want {
		t.Fatalf("after commit: buffer = %q, want %q", got, want)
	}
}

// TestCompletionCommitWithoutPreviewInsertsSuggestion covers the defensive
// fallback: if a suggestion is selected without the inline preview having
// written it into the buffer, committing must still insert it (replacing the
// partially typed token) rather than silently do nothing.
func TestCompletionCommitWithoutPreviewInsertsSuggestion(t *testing.T) {
	p := newCompletionTestPrompt(quotedSuggestion, 3)

	// A partially typed token, cursor at the end.
	p.buffer.InsertTextMoveCursor(`rm "fi`, DefColCount, DefRowCount, false)
	p.completion.Update(*p.buffer.Document())

	// Select directly, bypassing updateSuggestions, so no inline preview is
	// written to the buffer.
	p.completion.Next()

	if !p.completion.Completing() {
		t.Fatal("expected a suggestion to be selected")
	}

	// The line before the cursor does not yet end with the suggestion, so the
	// fallback must replace the partial token with the full suggestion.
	p.handleCompletionKeyBinding(nil, Enter, p.completion.Completing())

	want := "rm " + quotedSuggestion
	if got := p.buffer.Text(); got != want {
		t.Fatalf("after commit without preview: buffer = %q, want %q", got, want)
	}
}
