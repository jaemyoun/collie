package prompt

import (
	"fmt"
	"os"
)

type promptCommandType struct {
	CommandName string
	Desc        string
	Do          func()
}

var promptCommands = []promptCommandType{
	{CommandName: "select-bucket", Desc: "Select/unselect S3 bucket to explore", Do: DoSelectBucket},
	{CommandName: "ls", Desc: "List objects", Do: DoLs},
	{CommandName: "ls-rec", Desc: "List objects recursively", Do: DoLsRecursively},
	{CommandName: "add-filter", Desc: "Add filter as regular expression", Do: addFilter},
	{CommandName: "list-filter", Desc: "List all filters", Do: listFilter},
	{CommandName: "rm-filter", Desc: "Remove filter in the list"},
	{CommandName: "check-date", Desc: "Check recent modified date of all objects with filtering"},
}

func Run() bool {
	p := newSelectDraft("CommandName")
	p.Items = promptCommands
	p.Searcher = search(func(index int) string {
		return promptCommands[index].CommandName
	})
	p.Templates.Details = fmt.Sprintf("%s\nDescription: {{ .Desc }}", barString)
	idx, _, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(-1)
	}
	promptCommands[idx].Do()
	return true
}
