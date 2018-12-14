package prompt

import "fmt"

func addFilter() {
	inputtedRegExpr := inputByPrompt("new reg. expr.:", func(input string) error {
		return nil
	})

	inputtedDesc := inputByPrompt("description of the new reg. expr.:", func(input string) error {
		return nil
	})

	fmt.Printf("::::::::::::: filter will add %s (%s) in the filter list\n", inputtedRegExpr, inputtedDesc)
}

func listFilter() {
	fmt.Printf("::::::::::::: print the filter list\n")
}

type filterListType struct {
	regExpr string
	desc    string
}

func removeFilter() {
	// list filter

	// select

	// remove the filter
	fmt.Printf("::::::::::::: remove the filter\n")
}
