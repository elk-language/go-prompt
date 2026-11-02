package prompt

import (
	"bytes"
	"io"
)

// WinSize represents the width and height of terminal.
type WinSize struct {
	Row uint16
	Col uint16
}

const (
	DefColCount = 80 // Default column count of the terminal
	DefRowCount = 25 // Default row count of the terminal
)

// Reader is an interface to abstract input layer.
type Reader interface {
	// Open should be called before starting reading
	Open() error
	// GetWinSize returns WinSize object to represent width and height of terminal.
	GetWinSize() *WinSize
	io.ReadCloser
}

// GetKey returns Key correspond to input byte codes.
func GetKey(b []byte) Key {
	for _, k := range ASCIISequences {
		if bytes.Equal(k.ASCIICode, b) {
			return k.Key
		}
	}
	return NotDefined
}

// ASCIISequences holds mappings of the key and byte array.
var ASCIISequences = []*ASCIICode{
	{Key: Escape, ASCIICode: []byte{0x1b}},

	{Key: ControlSpace, ASCIICode: []byte{0x00}},
	{Key: ControlA, ASCIICode: []byte{0x1}},
	{Key: ControlB, ASCIICode: []byte{0x2}},
	{Key: ControlC, ASCIICode: []byte{0x3}},
	{Key: ControlD, ASCIICode: []byte{0x4}},
	{Key: ControlE, ASCIICode: []byte{0x5}},
	{Key: ControlF, ASCIICode: []byte{0x6}},
	{Key: ControlG, ASCIICode: []byte{0x7}},
	{Key: ControlH, ASCIICode: []byte{0x8}},
	//{Key: ControlI, ASCIICode: []byte{0x9}},
	//{Key: ControlJ, ASCIICode: []byte{0xa}},
	{Key: ControlK, ASCIICode: []byte{0xb}},
	{Key: ControlL, ASCIICode: []byte{0xc}},
	{Key: ControlM, ASCIICode: []byte{0xd}},
	{Key: ControlN, ASCIICode: []byte{0xe}},
	{Key: ControlO, ASCIICode: []byte{0xf}},
	{Key: ControlP, ASCIICode: []byte{0x10}},
	{Key: ControlQ, ASCIICode: []byte{0x11}},
	{Key: ControlR, ASCIICode: []byte{0x12}},
	{Key: ControlS, ASCIICode: []byte{0x13}},
	{Key: ControlT, ASCIICode: []byte{0x14}},
	{Key: ControlU, ASCIICode: []byte{0x15}},
	{Key: ControlV, ASCIICode: []byte{0x16}},
	{Key: ControlW, ASCIICode: []byte{0x17}},
	{Key: ControlX, ASCIICode: []byte{0x18}},
	{Key: ControlY, ASCIICode: []byte{0x19}},
	{Key: ControlZ, ASCIICode: []byte{0x1a}},

	{Key: ControlBackslash, ASCIICode: []byte{0x1c}},
	{Key: ControlSquareClose, ASCIICode: []byte{0x1d}},
	{Key: ControlCircumflex, ASCIICode: []byte{0x1e}},
	{Key: ControlUnderscore, ASCIICode: []byte{0x1f}},
	{Key: Backspace, ASCIICode: []byte{0x7f}},
	{Key: AltBackspace, ASCIICode: []byte{0x1b, 0x7f}},

	{Key: Up, ASCIICode: []byte{0x1b, '[', 'A'}},
	{Key: Down, ASCIICode: []byte{0x1b, '[', 'B'}},
	{Key: Right, ASCIICode: []byte{0x1b, '[', 'C'}},
	{Key: EscapeRight, ASCIICode: []byte{0x1b, 0x1b, '[', 'C'}},
	{Key: AltRight, ASCIICode: []byte{0x1b, '[', '1', ';', '3', 'C'}},
	{Key: Left, ASCIICode: []byte{0x1b, '[', 'D'}},
	{Key: AltLeft, ASCIICode: []byte{0x1b, '[', '1', ';', '3', 'D'}},
	{Key: EscapeLeft, ASCIICode: []byte{0x1b, 0x1b, '[', 'D'}},
	{Key: Home, ASCIICode: []byte{0x1b, '[', 'H'}},
	{Key: Home, ASCIICode: []byte{0x1b, '0', 'H'}},
	{Key: End, ASCIICode: []byte{0x1b, '[', 'F'}},
	{Key: End, ASCIICode: []byte{0x1b, '0', 'F'}},

	{Key: Enter, ASCIICode: []byte{'\n'}},
	{Key: ShiftEnter, ASCIICode: []byte{0x1b, '[', '1', '3', ';', '2', 'u'}},
	{Key: Delete, ASCIICode: []byte{0x1b, '[', '3', '~'}},
	{Key: AltDelete, ASCIICode: []byte{0x1b, '[', '3', ';', '3', '~'}},
	{Key: MetaDelete, ASCIICode: []byte{0x1b, '[', '3', ';', '9', '~'}},
	{Key: ShiftDelete, ASCIICode: []byte{0x1b, '[', '3', ';', '2', '~'}},
	{Key: ControlDelete, ASCIICode: []byte{0x1b, '[', '3', ';', '5', '~'}},
	{Key: Home, ASCIICode: []byte{0x1b, '[', '1', '~'}},
	{Key: End, ASCIICode: []byte{0x1b, '[', '4', '~'}},
	{Key: PageUp, ASCIICode: []byte{0x1b, '[', '5', '~'}},
	{Key: PageDown, ASCIICode: []byte{0x1b, '[', '6', '~'}},
	{Key: Home, ASCIICode: []byte{0x1b, '[', '7', '~'}},
	{Key: End, ASCIICode: []byte{0x1b, '[', '8', '~'}},
	{Key: Tab, ASCIICode: []byte{'\t'}},
	{Key: BackTab, ASCIICode: []byte{0x1b, '[', 'Z'}},
	{Key: Insert, ASCIICode: []byte{0x1b, '[', '2', '~'}},

	{Key: F1, ASCIICode: []byte{0x1b, 'O', 'P'}},
	{Key: F2, ASCIICode: []byte{0x1b, 'O', 'Q'}},
	{Key: F3, ASCIICode: []byte{0x1b, 'O', 'R'}},
	{Key: F4, ASCIICode: []byte{0x1b, 'O', 'S'}},

	{Key: F1, ASCIICode: []byte{0x1b, 'O', 'P', 'A'}}, // Linux console
	{Key: F2, ASCIICode: []byte{0x1b, '[', '[', 'B'}}, // Linux console
	{Key: F3, ASCIICode: []byte{0x1b, '[', '[', 'C'}}, // Linux console
	{Key: F4, ASCIICode: []byte{0x1b, '[', '[', 'D'}}, // Linux console
	{Key: F5, ASCIICode: []byte{0x1b, '[', '[', 'E'}}, // Linux console

	{Key: F1, ASCIICode: []byte{0x1b, '[', 0x11, '~'}}, // rxvt-unicode
	{Key: F2, ASCIICode: []byte{0x1b, '[', 0x12, '~'}}, // rxvt-unicode
	{Key: F3, ASCIICode: []byte{0x1b, '[', 0x13, '~'}}, // rxvt-unicode
	{Key: F4, ASCIICode: []byte{0x1b, '[', 0x14, '~'}}, // rxvt-unicode

	{Key: F5, ASCIICode: []byte{0x1b, '[', '1', '5', '~'}},
	{Key: F6, ASCIICode: []byte{0x1b, '[', '1', '7', '~'}},
	{Key: F7, ASCIICode: []byte{0x1b, '[', '1', '8', '~'}},
	{Key: F8, ASCIICode: []byte{0x1b, '[', '1', '9', '~'}},
	{Key: F9, ASCIICode: []byte{0x1b, '[', '2', '0', '~'}},
	{Key: F10, ASCIICode: []byte{0x1b, '[', '2', '1', '~'}},
	{Key: F11, ASCIICode: []byte{0x1b, '[', '2', '2', '~'}},
	{Key: F12, ASCIICode: []byte{0x1b, '[', '2', '4', '~', '\b'}},
	{Key: F13, ASCIICode: []byte{0x1b, '[', '%', '~'}},
	{Key: F14, ASCIICode: []byte{0x1b, '[', '&', '~'}},
	{Key: F15, ASCIICode: []byte{0x1b, '[', '(', '~'}},
	{Key: F16, ASCIICode: []byte{0x1b, '[', ')', '~'}},
	{Key: F17, ASCIICode: []byte{0x1b, '[', '1', '~'}},
	{Key: F18, ASCIICode: []byte{0x1b, '[', '2', '~'}},
	{Key: F19, ASCIICode: []byte{0x1b, '[', '3', '~'}},
	{Key: F20, ASCIICode: []byte{0x1b, '[', '4', '~'}},

	// Xterm
	{Key: F13, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'P'}},
	{Key: F14, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'Q'}},
	// &ASCIICode{Key: F15, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'R'}},  // Conflicts with CPR response
	{Key: F16, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'R'}},
	{Key: F17, ASCIICode: []byte{0x1b, '[', 0x15, ';', '2', '~'}},
	{Key: F18, ASCIICode: []byte{0x1b, '[', 0x17, ';', '2', '~'}},
	{Key: F19, ASCIICode: []byte{0x1b, '[', 0x18, ';', '2', '~'}},
	{Key: F20, ASCIICode: []byte{0x1b, '[', 0x19, ';', '2', '~'}},
	{Key: F21, ASCIICode: []byte{0x1b, '[', ' ', ';', '2', '~'}},
	{Key: F22, ASCIICode: []byte{0x1b, '[', '!', ';', '2', '~'}},
	{Key: F23, ASCIICode: []byte{0x1b, '[', '#', ';', '2', '~'}},
	{Key: F24, ASCIICode: []byte{0x1b, '[', '$', ';', '2', '~'}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, '[', '1', ';', '5', 'A'}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, '[', '1', ';', '5', 'B'}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, '[', '1', ';', '5', 'C'}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, '[', '1', ';', '5', 'D'}},

	{Key: ShiftUp, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'A'}},
	{Key: ShiftDown, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'B'}},
	{Key: ShiftRight, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'C'}},
	{Key: ShiftLeft, ASCIICode: []byte{0x1b, '[', '1', ';', '2', 'D'}},

	// Tmux sends following keystrokes when control+arrow is pressed, but for
	// Emacs ansi-term sends the same sequences for normal arrow keys. Consider
	// it a normal arrow press, because that's more important.
	{Key: Up, ASCIICode: []byte{0x1b, 'O', 'A'}},
	{Key: Down, ASCIICode: []byte{0x1b, 'O', 'B'}},
	{Key: Right, ASCIICode: []byte{0x1b, 'O', 'C'}},
	{Key: Left, ASCIICode: []byte{0x1b, 'O', 'D'}},

	{Key: ControlUp, ASCIICode: []byte{0x1b, '[', '5', 'A'}},
	{Key: ControlDown, ASCIICode: []byte{0x1b, '[', '5', 'B'}},
	{Key: ControlRight, ASCIICode: []byte{0x1b, '[', '5', 'C'}},
	{Key: ControlLeft, ASCIICode: []byte{0x1b, '[', '5', 'D'}},

	{Key: ControlRight, ASCIICode: []byte{0x1b, '[', 'O', 'c'}}, // rxvt
	{Key: ControlLeft, ASCIICode: []byte{0x1b, '[', 'O', 'd'}},  // rxvt

	{Key: Ignore, ASCIICode: []byte{0x1b, '[', 'E'}}, // Xterm
	{Key: Ignore, ASCIICode: []byte{0x1b, '[', 'F'}}, // Linux console
	{Key: EscapeT, ASCIICode: []byte{0x1b, 't'}},
	{Key: EscapeD, ASCIICode: []byte{0x1b, 'd'}},
	{Key: EscapeU, ASCIICode: []byte{0x1b, 'u'}},
	{Key: EscapeL, ASCIICode: []byte{0x1b, 'l'}},
	{Key: EscapeC, ASCIICode: []byte{0x1b, 'c'}},
	{Key: EscapeB, ASCIICode: []byte{0x1b, 'b'}},
	{Key: EscapeF, ASCIICode: []byte{0x1b, 'f'}},
	{Key: EscapeLess, ASCIICode: []byte{0x1b, '<'}},
	{Key: EscapeGreater, ASCIICode: []byte{0x1b, '>'}},
}
