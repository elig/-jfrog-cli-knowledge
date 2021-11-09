package utils

import (
	"bytes"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"io"
	"os"
)

const (
	promptItemTemplate = " {{ .Option | cyan }}{{if .TargetValue}}({{ .TargetValue }}){{end}}"
)

type PromptItem struct {
	Id           int64
	Option       string
	TargetValue  *string
	DefaultValue string
}

func PromptStringsNew(items []PromptItem, label string) (PromptItem, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🐸" + promptItemTemplate,
		Inactive: "  " + promptItemTemplate,
	}
	prompt := promptui.Select{
		Label:        label,
		Templates:    templates,
		Stdout:       &bellSkipper{},
		HideSelected: true,
		Size:         10,
		Items:        items,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return PromptItem{}, nil
	}

	return items[i], nil
}

// In MacOS, Terminal bell is ringing when trying to select items using up and down arrows.
// Using bellSkipper as Stdout is a workaround for this issue.
type bellSkipper struct{ io.WriteCloser }

var charBell = []byte{readline.CharBell}

func (bs *bellSkipper) Write(b []byte) (int, error) {
	if bytes.Equal(b, charBell) {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
