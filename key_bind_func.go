package prompt

// GoLineEnd Go to the End of the line
func GoLineEnd(p *Prompt) bool {
	return p.CursorEndOfLine()
}

// GoLineBeginning Go to the beginning of the line
func GoLineBeginning(p *Prompt) bool {
	return p.CursorStartOfLine()
}

// DeleteChar Delete character under the cursor
func DeleteChar(p *Prompt) bool {
	return p.DeleteCharacterUnderCursor()
}

// DeleteCharacterBeforeCursor Delete character before the cursor
func DeleteCharacterBeforeCursor(p *Prompt) bool {
	return p.DeleteCharacterUnderCursor()
}

// DeleteBeforeChar Go to Backspace
func DeleteBeforeChar(p *Prompt) bool {
	return p.DeleteCharacterBeforeCursor()
}

func DeleteLineAfterCursor(p *Prompt) bool {
	return p.DeleteLineAfterCursor()
}

func DeleteLineBeforeCursor(p *Prompt) bool {
	return p.DeleteLineBeforeCursor()
}

// GoRightChar Forward one character
func GoRightChar(p *Prompt) bool {
	return p.CursorRight(1)
}

// GoRightWord Forward one word
func GoRightWord(p *Prompt) bool {
	return p.CursorRightWord()
}

// GoLeftChar Backward one character
func GoLeftChar(p *Prompt) bool {
	return p.CursorLeft(1)
}

// GoLeftWord Backward one word
func GoLeftWord(p *Prompt) bool {
	return p.CursorLeftWord()
}

func DeleteWordBeforeCursor(p *Prompt) bool {
	return p.DeleteWordBeforeCursor()
}

// Delete the word after the cursor.
func DeleteWordAfterCursor(p *Prompt) bool {
	return p.DeleteWordAfterCursor()
}

// Uppercase the word after the cursor.
func UpcaseWordAfterCursor(p *Prompt) bool {
	return p.UpcaseWordAfterCursor()
}

// Lowercase the word after the cursor.
func DowncaseWordAfterCursor(p *Prompt) bool {
	return p.DowncaseWordAfterCursor()
}

// Capitalize the word after the cursor.
func CapitalizeWordAfterCursor(p *Prompt) bool {
	return p.CapitalizeWordAfterCursor()
}

// Move to the oldest history entry.
func HistoryOldest(p *Prompt) bool {
	return p.HistoryOldest()
}

// Move to the newest history entry.
func HistoryNewest(p *Prompt) bool {
	return p.HistoryNewest()
}

// Swap the last two characters before the cursor.
func SwapCharactersBeforeCursor(p *Prompt) bool {
	return p.SwapCharactersBeforeCursor()
}

// Swap the last two words before the cursor.
func SwapWordsBeforeCursor(p *Prompt) bool {
	return p.SwapWordsBeforeCursor()
}

func ClearScreen(p *Prompt) bool {
	p.ClearScreen()
	return true
}

func InsertNewline(p *Prompt) bool {
	p.InsertNewline()
	return true
}
