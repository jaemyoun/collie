package config

import (
	"fmt"
	"strings"
	"time"
)

func GetCheckDate() [2]CheckDateInfo {
	return setting.checkDate
}

func ParseDate(input string) (t time.Time, err error) {
	layout := []string{"2006-01-02T15:04:05Z", "2006-01-02 15:04:05 Z", "2006/01/02/15"}
	for _, e := range layout {
		t, err = time.Parse(e, input)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Now(), fmt.Errorf("parsing time '%s' as '%s': cannot parse it",
			input, strings.Join(layout, "', or '"))
	}
	return t, err
}

func AddCheckDate(operationPosition, operationInt int, timeString string) (time.Time, error) {
	t, err := ParseDate(timeString)
	if err != nil {
		fmt.Println("failed to parse inputted time string,", timeString, err)
		return time.Now(), err
	}

	setting.checkDate[operationPosition] = CheckDateInfo{Operation: OperationType(operationInt), Time: t}
	return t, nil
}

func DeleteCheckDate(operationPosition int) (string, string) {
	switch operationPosition {
	case LeftOperation, RightOperation:
		op := OperationString[setting.checkDate[operationPosition].Operation]
		t := setting.checkDate[operationPosition].Time.String()
		setting.checkDate[operationPosition] = CheckDateInfo{}
		return op, t
	}
	return "", ""
}
