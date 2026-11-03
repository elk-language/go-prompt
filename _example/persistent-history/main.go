package main

import (
	"fmt"
	"os"
	"path/filepath"

	prompt "github.com/elk-language/go-prompt"
)

var LivePrefix *prompt.Prefix = &prompt.Prefix{
	Content:   ">>> ",
	TextColor: prompt.Blue,
}

func executor(in string) {
	fmt.Println("Your input: " + in)

	if len(in)%2 == 0 {
		LivePrefix = &prompt.Prefix{
			Content:   `\\\ `,
			TextColor: prompt.Yellow,
		}
	} else {
		LivePrefix = &prompt.Prefix{
			Content:   ">>> ",
			TextColor: prompt.Blue,
		}
	}
}

func changeLivePrefix() *prompt.Prefix {
	return LivePrefix
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	historyPath := filepath.Join(cwd, "history")

	p := prompt.New(
		executor,
		prompt.WithExtendedPrefixCallback(changeLivePrefix),
		prompt.WithTitle("persistent-history"),
		prompt.WithFSHistory(historyPath),
	)
	p.Run()
}
