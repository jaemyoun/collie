package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"strings"
)

func search(elem func(index int) string) list.Searcher {
	return func(input string, index int) bool {
		name := strings.Replace(strings.ToLower(elem(index)), " ", "", -1)
		return strings.Contains(name, strings.Replace(strings.ToLower(input), " ", "", -1))
	}
}

func newSelectDraft(mainVar string) promptui.Select {
	return promptui.Select{
		Label: "> ",
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   fmt.Sprintf("\U00002714 {{ .%s | cyan }}", mainVar),
			Inactive: fmt.Sprintf("  {{ .%s | cyan }}", mainVar),
			Selected: fmt.Sprintf("\U00002714 {{ .%s | red }}", mainVar),
		},
		Size: 10,
	}
}

func inputByPrompt(label string, validate func(input string) error) string {
	prompt := promptui.Prompt{
		Label: label,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | yellow }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return result
}
