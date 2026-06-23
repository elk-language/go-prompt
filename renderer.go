package prompt

import (
	"strings"
	"unicode/utf8"

	"github.com/elk-language/go-prompt/debug"
	istrings "github.com/elk-language/go-prompt/strings"
)

const multilinePrefixCharacter = '.'

// Takes care of the rendering process
type Renderer struct {
	out                    Writer
	prefixCallback         PrefixCallback
	extendedPrefixCallback ExtendedPrefixCallback
	breakLineCallback      func(*Document)
	title                  string
	row                    int
	col                    istrings.Width
	indentSize             int // How many spaces constitute a single indentation level

	previousCursor Position

	// colors,
	prefixTextColor              Color
	prefixBGColor                Color
	inputTextColor               Color
	inputBGColor                 Color
	suggestionTextColor          Color
	suggestionBGColor            Color
	selectedSuggestionTextColor  Color
	selectedSuggestionBGColor    Color
	descriptionTextColor         Color
	descriptionBGColor           Color
	selectedDescriptionTextColor Color
	selectedDescriptionBGColor   Color
	scrollbarThumbColor          Color
	scrollbarBGColor             Color
}

// Build a new Renderer.
func NewRenderer() *Renderer {
	defaultWriter := NewStdoutWriter()
	registerWriter(defaultWriter)

	return &Renderer{
		out:                          defaultWriter,
		indentSize:                   DefaultIndentSize,
		prefixCallback:               DefaultPrefixCallback,
		prefixTextColor:              Blue,
		prefixBGColor:                DefaultColor,
		inputTextColor:               DefaultColor,
		inputBGColor:                 DefaultColor,
		suggestionTextColor:          White,
		suggestionBGColor:            Cyan,
		selectedSuggestionTextColor:  Black,
		selectedSuggestionBGColor:    Turquoise,
		descriptionTextColor:         Black,
		descriptionBGColor:           Turquoise,
		selectedDescriptionTextColor: White,
		selectedDescriptionBGColor:   Cyan,
		scrollbarThumbColor:          DarkGray,
		scrollbarBGColor:             Cyan,
	}
}

func (r *Renderer) getPrefix() *Prefix {
	if r.extendedPrefixCallback != nil {
		return r.extendedPrefixCallback()
	}

	return &Prefix{
		Content:         r.prefixCallback(),
		BackgroundColor: r.prefixBGColor,
		TextColor:       r.prefixTextColor,
	}
}

// Setup to initialize console output.
func (r *Renderer) Setup() {
	if r.title != "" {
		r.out.SetTitle(r.title)
		r.flush()
	}
}

func (r *Renderer) renderPrefix(prefix *Prefix) {
	r.out.SetColor(prefix.TextColor, prefix.BackgroundColor, false)
	if _, err := r.out.WriteString("\r"); err != nil {
		panic(err)
	}
	if _, err := r.out.WriteString(prefix.Content); err != nil {
		panic(err)
	}
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Close to clear title and erase.
func (r *Renderer) Close() {
	r.out.ClearTitle()
	r.out.EraseDown()
	r.flush()
}

func (r *Renderer) prepareArea(lines int) {
	for i := 0; i < lines; i++ {
		r.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		r.out.ScrollUp()
	}
}

// UpdateWinSize called when window size is changed.
func (r *Renderer) UpdateWinSize(ws *WinSize) {
	r.row = int(ws.Row)
	r.col = istrings.Width(ws.Col)
}

func (r *Renderer) renderCompletion(buf *Buffer, completions *CompletionManager) {
	suggestions := completions.GetSuggestions()
	if len(suggestions) == 0 {
		return
	}
	prefix := r.getPrefix()
	prefixWidth := istrings.GetWidth(prefix.Content)
	formatted, width := formatSuggestions(
		suggestions,
		r.col-istrings.GetWidth(prefix.Content)-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	windowHeight := len(formatted)
	if windowHeight > int(completions.max) {
		windowHeight = int(completions.max)
	}
	formatted = formatted[completions.verticalScroll : completions.verticalScroll+windowHeight]
	r.prepareArea(windowHeight)

	cursor := positionAtEndOfString(buf.Document().TextBeforeCursor(), r.col-prefixWidth)
	cursor.X += prefixWidth
	x := cursor.X
	if x+width >= r.col {
		cursor = r.backward(cursor, x+width-r.col)
	}

	contentHeight := len(completions.tmp)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.selected - completions.verticalScroll
	cursorColumnSpacing := cursor

	r.out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		alignNextLine(r, cursorColumnSpacing.X)

		if i == selected {
			r.out.SetColor(r.selectedSuggestionTextColor, r.selectedSuggestionBGColor, true)
		} else {
			r.out.SetColor(r.suggestionTextColor, r.suggestionBGColor, false)
		}
		if _, err := r.out.WriteString(formatted[i].Text); err != nil {
			panic(err)
		}

		if i == selected {
			r.out.SetColor(r.selectedDescriptionTextColor, r.selectedDescriptionBGColor, false)
		} else {
			r.out.SetColor(r.descriptionTextColor, r.descriptionBGColor, false)
		}
		if _, err := r.out.WriteString(formatted[i].Description); err != nil {
			panic(err)
		}

		if isScrollThumb(i) {
			r.out.SetColor(DefaultColor, r.scrollbarThumbColor, false)
		} else {
			r.out.SetColor(DefaultColor, r.scrollbarBGColor, false)
		}
		if _, err := r.out.WriteString(" "); err != nil {
			panic(err)
		}
		r.out.SetColor(DefaultColor, DefaultColor, false)

		c := cursor.Add(Position{X: width})
		r.backward(c, width)
	}

	if x+width >= r.col {
		r.out.CursorForward(int(x + width - r.col + 1))
	}

	r.out.CursorUp(windowHeight)
	r.out.SetColor(DefaultColor, DefaultColor, false)
}

// Render renders to the console.
func (r *Renderer) Render(buffer *Buffer, completion *CompletionManager, lexer Lexer) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if r.col == 0 {
		return
	}
	defer func() { r.flush() }()
	r.clear(r.previousCursor)

	text := buffer.Text()
	prefix := r.getPrefix()
	prefixWidth := istrings.GetWidth(prefix.Content)
	col := r.col - prefixWidth
	endLine := buffer.startLine + int(r.row) - 1
	cursor := positionAtEndOfStringLine(text, col, endLine)
	cursor.X += prefixWidth

	// Rendering
	r.out.HideCursor()
	defer r.out.ShowCursor()

	r.renderText(lexer, buffer.Text(), buffer.startLine)

	r.out.SetColor(DefaultColor, DefaultColor, false)

	targetCursor := buffer.DisplayCursorPosition(col)
	targetCursor.X += prefixWidth
	// Log("col: %#v, targetCursor: %#v, cursor: %#v\n", col, targetCursor, cursor)
	cursor = r.move(cursor, targetCursor)

	r.renderCompletion(buffer, completion)
	r.previousCursor = cursor
}

func (r *Renderer) renderText(lexer Lexer, input string, startLine int) {
	if lexer != nil {
		r.lex(lexer, input, startLine)
		return
	}

	prefix := r.getPrefix()
	prefixWidth := istrings.GetWidth(prefix.Content)
	col := r.col - prefixWidth
	multilinePrefix := r.getMultilinePrefix(prefix)
	if startLine != 0 {
		prefix = multilinePrefix
	}
	firstIteration := true
	endLine := startLine + int(r.row)
	var lineBuffer strings.Builder
	var lineCharIndex istrings.Width
	var lineNumber int

	for _, char := range input {
		if lineCharIndex >= col || char == '\n' {
			lineNumber++
			lineCharIndex = 0
			if lineNumber-1 < startLine {
				continue
			}
			if lineNumber >= endLine {
				break
			}
			lineBuffer.WriteRune('\n')
			r.renderLine(prefix, lineBuffer.String(), r.inputTextColor)
			lineBuffer.Reset()
			if char != '\n' {
				lineBuffer.WriteRune(char)
				lineCharIndex += istrings.GetRuneWidth(char)
			}
			if firstIteration {
				prefix = multilinePrefix
				firstIteration = false
			}
			continue
		}

		lineCharIndex += istrings.GetRuneWidth(char)
		if lineNumber < startLine {
			continue
		}
		lineBuffer.WriteRune(char)
	}

	r.renderLine(prefix, lineBuffer.String(), r.inputTextColor)
}

func (r *Renderer) flush() {
	debug.AssertNoError(r.out.Flush())
}

func (r *Renderer) renderLine(prefix *Prefix, line string, color Color) {
	r.renderPrefix(prefix)
	r.writeStringColor(line, color)
}

func (r *Renderer) writeStringColor(text string, color Color) {
	r.out.SetColor(color, r.inputBGColor, false)
	if _, err := r.out.WriteString(text); err != nil {
		panic(err)
	}
}

func (r *Renderer) write(b []byte) {
	if _, err := r.out.Write(b); err != nil {
		panic(err)
	}
}

func (r *Renderer) getMultilinePrefix(prefix *Prefix) *Prefix {
	content := prefix.Content

	var spaceCount int
	var dotCount int
	var nonSpaceCharSeen bool
	for len(content) != 0 {
		char, size := utf8.DecodeLastRuneInString(content)
		content = content[:len(content)-size]
		charWidth := istrings.GetRuneWidth(char)
		if nonSpaceCharSeen {
			dotCount += int(charWidth)
			continue
		}
		if char != ' ' {
			nonSpaceCharSeen = true
			dotCount += int(charWidth)
			continue
		}
		spaceCount += int(charWidth)
	}

	var multilinePrefixBuilder strings.Builder

	for i := 0; i < dotCount; i++ {
		multilinePrefixBuilder.WriteByte(multilinePrefixCharacter)
	}
	for i := 0; i < spaceCount; i++ {
		multilinePrefixBuilder.WriteByte(IndentUnit)
	}

	return &Prefix{
		Content:         multilinePrefixBuilder.String(),
		TextColor:       prefix.TextColor,
		BackgroundColor: prefix.BackgroundColor,
	}
}

// lex processes the given input with the given lexer
// and writes the result
func (r *Renderer) lex(lexer Lexer, input string, startLine int) {
	prefix := r.getPrefix()
	prefixWidth := istrings.GetWidth(prefix.Content)
	col := r.col - prefixWidth
	multilinePrefix := r.getMultilinePrefix(prefix)
	var lineCharIndex istrings.Width
	var lineNumber int
	endLine := startLine + int(r.row)
	previousByteIndex := istrings.ByteNumber(-1)
	lineBuffer := make([]byte, 8)
	runeBuffer := make([]byte, utf8.UTFMax)

	lexer.Init(input)

	if startLine != 0 {
		prefix = multilinePrefix
	}
	r.renderPrefix(prefix)

tokenLoop:
	for {
		token, ok := lexer.Next()
		var currentFirstByteIndex istrings.ByteNumber
		var currentLastByteIndex istrings.ByteNumber
		var tokenColor Color
		var tokenBackgroundColor Color
		var tokenDisplayAttributes []DisplayAttribute
		var noToken bool
		if ok {
			currentFirstByteIndex = token.FirstByteIndex()
			currentLastByteIndex = token.LastByteIndex()
			tokenColor = token.Color()
			tokenBackgroundColor = token.BackgroundColor()
			tokenDisplayAttributes = token.DisplayAttributes()
		} else if previousByteIndex == istrings.Len(input)-1 {
			break tokenLoop
		} else {
			currentFirstByteIndex = istrings.Len(input)
			tokenColor = r.inputTextColor
			tokenBackgroundColor = r.inputBGColor
			tokenDisplayAttributes = nil
			noToken = true
		}

		color := r.inputTextColor
		backgroundColor := r.inputBGColor
		displayAttributes := tokenDisplayAttributes
		text := input[previousByteIndex+1 : currentFirstByteIndex]
		previousByteIndex = currentLastByteIndex
		lineBuffer = lineBuffer[:0]
		interToken := true

	repeatLoop:
		for {

		charLoop:
			for _, char := range text {
				if lineCharIndex >= col || char == '\n' {
					lineNumber++
					lineCharIndex = 0
					if lineNumber-1 < startLine {
						continue charLoop
					}
					if lineNumber >= endLine {
						break tokenLoop
					}
					lineBuffer = append(lineBuffer, '\n')
					r.out.SetDisplayAttributes(color, backgroundColor, displayAttributes...)
					r.write(lineBuffer)
					r.resetFormatting()
					r.renderPrefix(multilinePrefix)
					lineBuffer = lineBuffer[:0]
					if char != '\n' {
						size := utf8.EncodeRune(runeBuffer, char)
						lineBuffer = append(lineBuffer, runeBuffer[:size]...)
						lineCharIndex += istrings.GetRuneWidth(char)
					}
					continue charLoop
				}

				lineCharIndex += istrings.GetRuneWidth(char)
				if lineNumber < startLine {
					continue charLoop
				}
				size := utf8.EncodeRune(runeBuffer, char)
				lineBuffer = append(lineBuffer, runeBuffer[:size]...)
			}
			if len(lineBuffer) > 0 {
				r.out.SetDisplayAttributes(color, backgroundColor, displayAttributes...)
				r.write(lineBuffer)
				r.resetFormatting()
			}

			if !interToken {
				break repeatLoop
			}

			if noToken {
				break tokenLoop
			}
			color = tokenColor
			backgroundColor = tokenBackgroundColor
			displayAttributes = tokenDisplayAttributes
			text = input[currentFirstByteIndex : currentLastByteIndex+1]
			lineBuffer = lineBuffer[:0]
			interToken = false
		}
	}

	r.resetFormatting()
}

func (r *Renderer) resetFormatting() {
	r.out.SetDisplayAttributes(r.inputTextColor, r.inputBGColor, DisplayReset)
}

// BreakLine to break line.
func (r *Renderer) BreakLine(buffer *Buffer, lexer Lexer) {
	// Erasing and Renderer
	prefix := r.getPrefix()
	prefixWidth := istrings.GetWidth(prefix.Content)
	cursor := positionAtEndOfString(buffer.Document().TextBeforeCursor(), r.col-prefixWidth)
	cursor.X += prefixWidth
	r.clear(cursor)

	r.renderText(lexer, buffer.Text(), buffer.startLine)
	if _, err := r.out.WriteString("\n"); err != nil {
		panic(err)
	}

	r.out.SetColor(DefaultColor, DefaultColor, false)

	r.flush()
	if r.breakLineCallback != nil {
		r.breakLineCallback(buffer.Document())
	}

	r.previousCursor = Position{}
}

// Get the number of columns that are available
// for user input.
func (r *Renderer) UserInputColumns() istrings.Width {
	prefix := r.getPrefix()
	return r.col - istrings.GetWidth(prefix.Content)
}

// clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (r *Renderer) clear(cursor Position) {
	r.move(cursor, Position{})
	r.out.EraseDown()
}

// backward moves cursor to backward from a current cursor position
// regardless there is a line break.
func (r *Renderer) backward(from Position, n istrings.Width) Position {
	return r.move(from, Position{X: from.X - n, Y: from.Y})
}

// move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (r *Renderer) move(from, to Position) Position {
	newPosition := from.Subtract(to)
	r.out.CursorUp(newPosition.Y)
	r.out.CursorBackward(int(newPosition.X))
	return to
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}

func alignNextLine(r *Renderer, col istrings.Width) {
	r.out.CursorDown(1)
	if _, err := r.out.WriteString("\r"); err != nil {
		panic(err)
	}
	r.out.CursorForward(int(col))
}
