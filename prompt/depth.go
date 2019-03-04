package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
	"strconv"
)

func setDepth() {
	depth := inputByPrompt("set limit of depths:", func(input string) error {
		if len(input) == 0 {
			return nil
		}
		_, err := strconv.Atoi(input)
		return err
	})
	if len(depth) == 0 {
		config.RemoveDepth()
	} else {
		n, _ := strconv.Atoi(depth)
		config.SetDepth(n)
	}
}

func getStatusForDepth() (ret string) {
	depth := config.GetDepth()
	if depth == -1 {
		return ""
	}
	ret += fmt.Sprintf("\U00002023 Depth: %d\n", depth)
	return ret
}
