package prompt

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

type promptCommandType struct {
	CommandName string
	Desc        string
	Do          func()
}

var promptCommands = []promptCommandType{
	{CommandName: "select bucket", Desc: "Select/unselect S3 bucket to explore", Do: selectBucket},
	{CommandName: "ls", Desc: "List objects", Do: ls},
	{CommandName: "ls recursively", Desc: "List objects recursively", Do: lsRecursively},
	{CommandName: "cd", Desc: "Change location (prefix) to list objects", Do: cd},
	{CommandName: "set filters", Desc: "Add/remove filter", Do: setFilter},
	{CommandName: "set date range", Desc: "Check recent modified date of all objects with filtering", Do: checkDate},
	{CommandName: "toggle details option", Desc: "Turn on/off printing list objects in details", Do: optionDetails},
	{CommandName: "toggle duplicated objects", Desc: "Toggle printing only duplicated files in multiple S3 Buckets", Do: optionDuplication},
}

func Run() bool {
	color.HiBlack(getStatus())

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

func getStatus() string {
	ret := barStringHead + "\n"
	ret += getStatusForSelectBucket()
	ret += getStatusForLocation()
	ret += getStatusForFilter()
	ret += getStatusForCheckDate()
	ret += getStatusOptionDetails()
	ret += getStatusOptionCheckDuplicated()
	return ret + barString
}
