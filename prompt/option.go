package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
)

func optionToggleDetails() {
	details := config.ToggleDetails()
	fmt.Printf("Option 'details' was set to be '%v'\n", details)
}
