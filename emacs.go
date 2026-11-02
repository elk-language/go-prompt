package prompt

/*

========
PROGRESS
========

Moving the cursor
-----------------

* [x] Ctrl + a   Go to the beginning of the line (Home)
* [x] Ctrl + e   Go to the End of the line (End)
* [x] Ctrl + p   Previous command (Up arrow)
* [x] Ctrl + n   Next command (Down arrow)
* [x] Ctrl + f   Forward one character
* [x] Ctrl + b   Backward one character
* [x] Ctrl + xx  Toggle between the start of line and current cursor position

Editing
-------

* [x] Ctrl + L   Clear the Screen, similar to the clear command
* [x] Ctrl + d   Delete character under the cursor
* [x] Ctrl + h   Delete character before the cursor (Backspace)

* [x] Ctrl + w   Cut the Word before the cursor to the clipboard.
* [x] Ctrl + k   Cut the Line after the cursor to the clipboard.
* [x] Ctrl + u   Cut/delete the Line before the cursor to the clipboard.

* [x] Ctrl + t   Swap the last two characters before the cursor (typo).
* [x] Esc  + t   Swap the last two words before the cursor.
* [x] Esc  + d   Delete the word after the cursor.
* [x] Esc  + u   Uppercase the word after the cursor.
* [x] Esc  + l   Lowercase the word after the cursor.
* [x] Esc  + c   Capitalize the word after the cursor.
* [x] Esc  + <   Oldest history entry.
* [x] Esc  + >   Newest history entry.

* [ ] ctrl + y   Paste the last thing to be cut (yank)
* [ ] ctrl + _   Undo

*/

var emacsKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: ControlE,
		Fn:  GoLineEnd,
	},
	// Go to the beginning of the line
	{
		Key: ControlA,
		Fn:  GoLineBeginning,
	},
	// Cut the Line after the cursor
	{
		Key: ControlK,
		Fn:  DeleteLineAfterCursor,
	},
	// Cut/delete the Line before the cursor
	{
		Key: ControlU,
		Fn:  DeleteLineBeforeCursor,
	},
	// Delete character under the cursor
	{
		Key: ControlD,
		Fn:  DeleteChar,
	},
	// Backspace
	{
		Key: ControlH,
		Fn:  DeleteCharacterBeforeCursor,
	},
	// Right allow: Forward one character
	{
		Key: ControlF,
		Fn:  GoRightChar,
	},
	// Alt Right allow: Forward one word
	{
		Key: AltRight,
		Fn:  GoRightWord,
	},
	// Escape, then the Right arrow: forward one word.
	{
		Key: EscapeRight,
		Fn:  GoRightWord,
	},
	// Escape+f: forward one word.
	{
		Key: EscapeF,
		Fn:  GoRightWord,
	},
	// Ctrl+Right: forward one word.
	{
		Key: ControlRight,
		Fn:  GoRightWord,
	},
	// Left allow: Backward one character
	{
		Key: ControlB,
		Fn:  GoLeftChar,
	},
	// Alt Left allow: Backward one word
	{
		Key: AltLeft,
		Fn:  GoLeftWord,
	},
	// Escape, then the Left arrow: backward one word.
	{
		Key: EscapeLeft,
		Fn:  GoLeftWord,
	},
	// Escape+b: backward one word.
	{
		Key: EscapeB,
		Fn:  GoLeftWord,
	},
	// Control Left allow: Backward one word
	{
		Key: ControlLeft,
		Fn:  GoLeftWord,
	},
	// Swap the last two characters before the cursor.
	{
		Key: ControlT,
		Fn:  SwapCharactersBeforeCursor,
	},
	// Swap the last two words before the cursor.
	{
		Key: EscapeT,
		Fn:  SwapWordsBeforeCursor,
	},
	// Delete the word after the cursor.
	{
		Key: EscapeD,
		Fn:  DeleteWordAfterCursor,
	},
	// Uppercase the word after the cursor.
	{
		Key: EscapeU,
		Fn:  UpcaseWordAfterCursor,
	},
	// Lowercase the word after the cursor.
	{
		Key: EscapeL,
		Fn:  DowncaseWordAfterCursor,
	},
	// Capitalize the word after the cursor.
	{
		Key: EscapeC,
		Fn:  CapitalizeWordAfterCursor,
	},
	// Oldest history entry.
	{
		Key: EscapeLess,
		Fn:  HistoryOldest,
	},
	// Newest history entry.
	{
		Key: EscapeGreater,
		Fn:  HistoryNewest,
	},
	// Insert a newline
	{
		Key: EscapeEnter,
		Fn:  InsertNewline,
	},
	// Insert a newline
	{
		Key: ShiftEnter,
		Fn:  InsertNewline,
	},
	// Insert a newline
	{
		Key: ControlO,
		Fn:  InsertNewline,
	},
	// Cut the Word before the cursor.
	{
		Key: ControlW,
		Fn:  DeleteWordBeforeCursor,
	},
	{
		Key: AltBackspace,
		Fn:  DeleteWordBeforeCursor,
	},
	// Clear the Screen, similar to the clear command
	{
		Key: ControlL,
		Fn:  ClearScreen,
	},
}
