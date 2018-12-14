package config

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"regexp"
	"time"
)

var setting = struct {
	selectedBucket  map[BucketNameType]BucketInfo
	currentLocation string
	addedFilters    []FilterInfo
	details         bool
	checkDate       [2]CheckDateInfo
}{}

type BucketNameType string

type BucketInfo struct {
	Region s3.BucketLocationConstraint
}

type FilterInfo struct {
	Filter *regexp.Regexp
	Desc   string
}

type CheckDateInfo struct {
	Operation OperationType
	Time      time.Time
}

type OperationType int

const (
	LeftOperation = iota
	RightOperation
)
const (
	OperationNotExist OperationType = iota
	OperationEqual
	OperationNotEqual
	OperationLess
	OperationEqualLess
)

var OperationString = []string{"", "==", "!=", "<", "<="}

func init() {
	setting.selectedBucket = make(map[BucketNameType]BucketInfo)
	setting.addedFilters = make([]FilterInfo, 0)
	readConfigFile()
}
