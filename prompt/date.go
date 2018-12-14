package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
	"time"
)

type checkDateItemType struct {
	Command   string
	Operation config.OperationType
}

func checkDate() {
	configCheckDate := config.GetCheckDate()
	left := configCheckDate[config.LeftOperation]
	right := configCheckDate[config.RightOperation]

	cd := make([]checkDateItemType, 3)
	cd[0] = checkDateItemType{Command: "(back)"}
	if left.Operation == config.OperationNotExist {
		cd[1] = checkDateItemType{Command: "(set left date. e.g. `left` <= x)", Operation: config.OperationNotExist}
	} else {
		cd[1] = checkDateItemType{
			Command:   fmt.Sprintf("%s %s", left.Time.String(), config.OperationString[left.Operation]),
			Operation: left.Operation}
	}
	if right.Operation == config.OperationNotExist {
		cd[2] = checkDateItemType{Command: "(set right date. e.g. x < `right`)", Operation: config.OperationNotExist}
	} else {
		cd[2] = checkDateItemType{
			Command:   fmt.Sprintf("%s %s", config.OperationString[right.Operation], right.Time.String()),
			Operation: right.Operation}
	}

	s := newSelectDraft("Command")
	s.Label = "Which filter do you want to set"
	s.Items = cd
	s.Searcher = search(func(index int) string {
		return cd[index].Command
	})
	idx, _, err := s.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch idx {
	case 0: // (back)
	case 1: // (left)
		if cd[idx].Operation == config.OperationNotExist {
			inputCheckDate(config.LeftOperation)
		} else {
			if op, t := config.DeleteCheckDate(config.LeftOperation); len(op) != 0 {
				fmt.Printf("New '%s %s' is removed\n", op, t)
			}
		}
	case 2: // (right)
		if cd[idx].Operation == config.OperationNotExist {
			inputCheckDate(config.RightOperation)
		} else {
			if op, t := config.DeleteCheckDate(config.RightOperation); len(op) != 0 {
				fmt.Printf("New '%s %s' is removed\n", op, t)
			}
		}
	}
}

func inputCheckDate(operationPosition int) {

	idx := 0
	var err error
	op := make([]checkDateItemType, 0)

	funcInputOperation := func() (int, error) {
		for idx, e := range config.OperationString {
			switch config.OperationType(idx) {
			case config.OperationEqual, config.OperationNotEqual, config.OperationLess, config.OperationEqualLess:
				op = append(op, checkDateItemType{Command: e})
			}
		}
		s := newSelectDraft("Command")
		s.Label = "Which filter do you want to set"
		s.Items = op
		s.Searcher = search(func(index int) string {
			return op[index].Command
		})
		idx, _, err = s.Run()
		if err != nil {
			fmt.Printf("failed to input operation: %v", err)
		}
		return idx + int(config.OperationEqual), err
	}

	if operationPosition == config.RightOperation {
		if idx, err = funcInputOperation(); err != nil {
			return
		}
	}

	input := inputByPrompt("new date:", func(input string) error {
		_, err := config.ParseDate(input)
		return err
	})

	if operationPosition == config.LeftOperation {
		if idx, err = funcInputOperation(); err != nil {
			return
		}
	}

	if t, err := config.AddCheckDate(operationPosition, idx, input); err == nil {
		fmt.Printf("New '%s %s' is added\n", config.OperationString[idx], t)
	}
}

func getStatusForCheckDate() (ret string) {
	checkDate := config.GetCheckDate()
	checkDateLeft := checkDate[config.LeftOperation]
	checkDateRight := checkDate[config.RightOperation]
	if checkDateLeft.Operation != config.OperationNotExist || checkDateRight.Operation != config.OperationNotExist {
		ret += "\U00002023 Checking Date List: "
		if checkDateLeft.Operation != config.OperationNotExist {
			ret += fmt.Sprintf("%s %s ",
				checkDateLeft.Time.String(), config.OperationString[checkDateLeft.Operation])
		}
		ret += "x"
		if checkDateRight.Operation != config.OperationNotExist {
			ret += fmt.Sprintf(" %s %s",
				config.OperationString[checkDateRight.Operation], checkDateRight.Time.String())
		}
		ret += "\n"
	}
	return ret
}

func validateDate(lastModified *time.Time) bool {
	ret := true
	checkDate := config.GetCheckDate()
	left := checkDate[config.LeftOperation]
	right := checkDate[config.RightOperation]

	T := left.Time
	switch left.Operation {
	case config.OperationEqual:
		ret = ret && T.Equal(*lastModified)
	case config.OperationNotEqual:
		ret = ret && !T.Equal(*lastModified)
	case config.OperationLess:
		ret = ret && T.Before(*lastModified)
	case config.OperationEqualLess:
		ret = ret && (T.Equal(*lastModified) || T.Before(*lastModified))
	}

	T = right.Time
	switch right.Operation {
	case config.OperationEqual:
		ret = ret && T.Equal(*lastModified)
	case config.OperationNotEqual:
		ret = ret && !T.Equal(*lastModified)
	case config.OperationLess:
		ret = ret && T.After(*lastModified)
	case config.OperationEqualLess:
		ret = ret && (T.Equal(*lastModified) || T.After(*lastModified))
	}

	return ret
}
