package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
	"regexp"
	"strings"
)

type filterItemType struct {
	Expr           string
	Desc           string
	PredefinedMark string
}

func setFilter() {
	added := make(map[string]bool)
	filters := append([]filterItemType{}, filterItemType{Expr: "(new)", Desc: "Define a new filter"})
	for _, e := range config.GetFilters() {
		filters = append(filters, filterItemType{Expr: e.Filter.String(), Desc: e.Desc, PredefinedMark: " "})
		added[e.Filter.String()] = true
	}
	predefinedOffset := len(filters)
	for _, e := range config.Predefined.GetFilters() {
		if _, ok := added[e.Expr]; ok == false {
			filters = append(filters, filterItemType{Expr: e.Expr, Desc: e.Desc, PredefinedMark: "\U00002605"})
		}
	}
	f := newSelectDraft("Expr")
	f.Label = "Which filter do you want to set"
	f.Items = filters
	f.Searcher = search(func(index int) string {
		return filters[index].Expr
	})
	f.Templates.Active = "\U00002714{{ .PredefinedMark | cyan }} {{ .Expr | cyan }}"
	f.Templates.Inactive = " {{ .PredefinedMark | cyan }} {{ .Expr | cyan }}"
	f.Templates.Details = fmt.Sprintf("%s\nDescription: {{ .Desc }}", barString)

	idx, _, err := f.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if idx == 0 {
		appendFilter("", "")
	} else if idx >= predefinedOffset {
		appendFilter(filters[idx].Expr, filters[idx].Desc)
	} else {
		subSetFilter(filters[idx].Expr)
	}
}

type filterSubMenuItemType struct {
	OptionName string
	Desc       string
}

var filterItems = []filterSubMenuItemType{
	{OptionName: "(back)", Desc: "Doesn't do anything"},
	{OptionName: "modify reg. expr.", Desc: "Modify the regular expresion of the filter"},
	{OptionName: "modify description", Desc: "Modify the description of the filter"},
	{OptionName: "delete", Desc: "Remove filter from the filter list"},
}

func subSetFilter(regExpr string) {
	p := newSelectDraft("OptionName")
	p.Items = filterItems
	p.Searcher = search(func(index int) string {
		return filterItems[index].OptionName
	})
	p.Templates.Details = fmt.Sprintf("%s\nDescription: {{ .Desc }}", barString)
	idx, _, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var ok bool
	switch idx {
	case 0:
		ok = true
	case 1:
		newRegExpr := inputRegExpr()
		if ok = config.SetRegExprOfFilter(regExpr, newRegExpr); ok {
			fmt.Printf("The filter '%s' was modified to '%s\n", regExpr, newRegExpr)
		}
	case 2:
		newDesc := inputDesc()
		if ok = config.SetDescOfFilter(regExpr, newDesc); ok {
			fmt.Printf("The description of the filter '%s' was modified to '%s'\n", regExpr, newDesc)
		}
	case 3:
		if ok = config.DeleteFilter(regExpr); ok {
			fmt.Printf("The filter '%s' was deleted\n", regExpr)
		}
	}
	if ok == false {
		fmt.Printf("The filter '%s' cannot found in the filter list", regExpr)
	}
}

func appendFilter(expr, desc string) {
	if len(expr) == 0 {
		expr = inputRegExpr()
		desc = inputDesc()
	}

	config.AddFilter(expr, desc)
	fmt.Printf("The filter '%s' was added\n", expr)
}

func inputDesc() string {
	input := inputByPrompt("description of the new reg. expr.:", func(input string) error {
		return nil
	})
	return input
}

func inputRegExpr() string {
	input := inputByPrompt("new reg. expr.:", func(input string) error {
		_, err := regexp.CompilePOSIX(input)
		return err
	})
	return input
}

func getStatusForFilter() (ret string) {
	filters := config.GetFilters()
	if len(filters) != 0 {
		title := "\U00002023 Filter List: "
		ret += fmt.Sprintf(title)
		first := true
		for idx, e := range filters {
			padding := strings.Repeat(" ", len(title)-2)
			if first == false {
				ret += fmt.Sprintf("\n%s", padding)
			}
			ret += fmt.Sprintf("%d) %s\n", idx+1, e.Filter.String())
			ret += fmt.Sprintf("%s  : %s", padding, e.Desc)
			first = false
		}
		ret += "\n"
	}
	return ret
}

func validateFilter(target string) bool {
	filters := config.GetFilters()
	if len(filters) == 0 {
		return true
	}
	for _, e := range filters {
		if e.Filter.Match([]byte(target)) {
			return true
		}
	}
	return false
}
