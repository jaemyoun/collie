package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
)

func cd() {
	location := inputByPrompt("new location (prefix):", func(input string) error {
		return nil // todo: ??
	})

	config.SetCurrentLocation(location)
}

func getStatusForLocation() (ret string) {
	location := config.GetCurrentLocation()
	if len(location) == 0 {
		location = "(root)"
	}
	ret += fmt.Sprintf("\U00002023 Current location: %s\n", location)
	return ret
}
