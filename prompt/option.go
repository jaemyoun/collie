package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
)

func optionDetails() {
	details := config.ToggleDetails()
	fmt.Printf("Option 'details' was set to be '%v'\n", details)
}

func getStatusOptionDetails() (ret string) {
	if config.GetDetails() {
		ret = "\U00002023 print in details\n"
	}
	return ret
}

func optionDuplication() {
	c := config.ToggleDuplication()
	fmt.Printf("Option 'Printing duplicated files' was set to be '%v'\n", c)
}

func getStatusOptionCheckDuplicated() (ret string) {
	if config.GetDuplication() {
		ret = "\U00002023 print only duplicated files\n"
	}
	return ret
}
