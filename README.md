# go-prompt

[![Go Report Card](https://goreportcard.com/badge/github.com/elk-language/go-prompt)](https://goreportcard.com/report/github.com/elk-language/go-prompt)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/elk-language/go-prompt?status.svg)](https://godoc.org/github.com/elk-language/go-prompt)
![tests](https://github.com/elk-language/go-prompt/workflows/tests/badge.svg)

This is a fork of [c-bata/go-prompt](https://github.com/c-bata/go-prompt).
It's a great library but it's been abandoned
for quite a while.
This project aims to continue its development.

The library has been rewritten in many aspects, fixing existing bugs and adding new essential functionality.

Most notable changes include:
- Support for custom syntax highlighting with a lexer
- Multiline editing
- A scrolling buffer is used for displaying the current content which makes it possible to edit text of arbitrary length (only the visible part of the text is rendered)
- Support for automatic indentation when pressing <kbd>Enter</kbd> and the input is incomplete or for executing the input when it is complete. This is determined by a custom callback function.

I highly encourage you to see the [changelog](CHANGELOG.md) which fully documents the changes that have been made.

---

A library for building powerful interactive prompts inspired by [python-prompt-toolkit](https://github.com/jonathanslenders/python-prompt-toolkit),
making it easier to build cross-platform command line tools using Go.

```go
package main

import (
	"fmt"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

func completer(d prompt.Document) ([]prompt.Suggest, pstrings.RuneNumber, pstrings.RuneNumber) {
	endIndex := d.CurrentRuneIndex()
	w := d.GetWordBeforeCursor()
	startIndex := endIndex - pstrings.RuneCount(w)

	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, w, true), startIndex, endIndex
}

func main() {
	fmt.Println("Please select table.")
	t := prompt.Input(
		prompt.WithPrefix("> "),
		prompt.WithCompleter(completer),
	)
	fmt.Println("You selected " + t)
}
```

## Features

### Automatic indentation with a custom callback

![automatic indentation](readme/automatic-indentation.gif)

### Multiline editing with scrolling

![multiline editing](readme/multiline-editing.gif)

### Custom syntax highlighting

![syntax highlighting](readme/syntax-highlighting.gif)

### Powerful auto-completion

[![autocompletion](https://github.com/c-bata/assets/raw/master/go-prompt/kube-prompt.gif)](https://github.com/c-bata/kube-prompt)

(This is a GIF animation of kube-prompt.)

### Flexible options

go-prompt provides many options. Please check [option section of GoDoc](https://godoc.org/github.com/elk-language/go-prompt#Option) for more details.

[![options](https://github.com/c-bata/assets/raw/master/go-prompt/prompt-options.png)](#flexible-options)

### Keyboard Shortcuts

Emacs-like keyboard shortcuts are available by default (these also are the default shortcuts in Bash shell).
You can customize and expand these shortcuts.

[![keyboard shortcuts](https://github.com/c-bata/assets/raw/master/go-prompt/keyboard-shortcuts.gif)](#keyboard-shortcuts)

Key Binding                                        | Description
----------------------------------------------------|---------------------------------------------------------
<kbd>Ctrl + A</kbd>                                | Go to the beginning of the line (Home)
<kbd>Ctrl + E</kbd>                                | Go to the end of the line (End)
<kbd>Ctrl + F</kbd>                                | Forward one character
<kbd>Ctrl + B</kbd>                                | Backward one character
<kbd>Ctrl + Right</kbd> / <kbd>Alt + Right</kbd> / <kbd>Esc</kbd> then <kbd>Right</kbd> / <kbd>Esc + F</kbd> | Forward one word
<kbd>Ctrl + Left</kbd> / <kbd>Alt + Left</kbd> / <kbd>Esc</kbd> then <kbd>Left</kbd> / <kbd>Esc + B</kbd> | Backward one word
<kbd>Ctrl + P</kbd>                                | Previous command (Up arrow)
<kbd>Ctrl + N</kbd>                                | Next command (Down arrow)
<kbd>Esc + &lt;</kbd>                              | Oldest history entry
<kbd>Esc + &gt;</kbd>                              | Newest history entry
<kbd>Ctrl + D</kbd>                                | Delete the character under the cursor
<kbd>Ctrl + H</kbd>                                | Delete the character under the cursor
<kbd>Ctrl + W</kbd> / <kbd>Alt + Backspace</kbd>   | Delete the word before the cursor
<kbd>Esc + D</kbd>                                 | Delete the word after the cursor
<kbd>Ctrl + K</kbd> / <kbd>Alt + Delete</kbd>     | Delete the line after the cursor
<kbd>Ctrl + U</kbd>                                | Delete the line before the cursor
<kbd>Ctrl + T</kbd>                                | Swap the last two characters before the cursor
<kbd>Esc + T</kbd>                                 | Swap the last two words before the cursor
<kbd>Esc + U</kbd>                                 | Uppercase the word after the cursor
<kbd>Esc + L</kbd>                                 | Lowercase the word after the cursor
<kbd>Esc + C</kbd>                                 | Capitalize the word after the cursor
<kbd>Ctrl + O</kbd> / <kbd>Shift + Enter</kbd> / <kbd>Esc + Enter</kbd> | Insert a newline
<kbd>Ctrl + L</kbd>                                | Clear the screen

### History

You can use <kbd>Up arrow</kbd> and <kbd>Down arrow</kbd> to walk through the history of commands executed.

[![History](https://github.com/c-bata/assets/raw/master/go-prompt/history.gif)](#history)

### Multiple platform support

We have confirmed go-prompt works fine in the following terminals:

* iTerm2 (macOS)
* Terminal.app (macOS)
* Command Prompt (Windows)
* gnome-terminal (Ubuntu)

## Links

* [Change Log](./CHANGELOG.md)
* [GoDoc](http://godoc.org/github.com/elk-language/go-prompt)
* [gocover.io](https://gocover.io/github.com/elk-language/go-prompt)

## License

This software is licensed under the MIT license, see [LICENSE](./LICENSE) for more information.

## Original Author

Masashi Shibata

* Twitter: [@c\_bata\_](https://twitter.com/c_bata_/)
* Github: [@c-bata](https://github.com/c-bata/)
