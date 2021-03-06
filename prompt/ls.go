package prompt

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/jaemyoun/collie/aws"
	"github.com/jaemyoun/collie/config"
	"github.com/jaemyoun/collie/gofunc"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func ls() {
	listObjectsInGoRoutine(false)
}

func lsRecursively() {
	listObjectsInGoRoutine(true)
}

func listObjectsInGoRoutine(rec bool) {
	duplicatedManager = duplicatedMangerType{scannedPath: make(map[string]struct{})}

	for bucket, info := range config.GetSelectedBucket() {
		color.Cyan("%s (%s):\n", bucket, info.Region)

		ch := make(chan interface{})
		routine := gofunc.New()
		routine.ChOutput = ch
		routine.Do = handleListObjects
		routine.AddInput(config.GetCurrentLocation())

		start := time.Now()
		routine.Run(string(bucket), string(info.Region), rec) //todo: get from config

		for v := range ch {
			output := v.(string)
			fmt.Print(output)
		}
		fmt.Printf("Request Count: %d, Elapse time: %v\n", routine.GetRunningCount(), time.Since(start))
		fmt.Printf("Total size: %v\n", totalSize)
		fmt.Println("")
	}
}

func handleListObjects(input interface{}, output chan<- interface{}, recursiveFunc func(v interface{}), v ...interface{}) {
	prefix := input.(string)
	vars := []interface{}(v)
	bucket := vars[0].(string)
	region := vars[1].(string)
	recursive := vars[2].(bool)

	details := config.GetDetails()
	duplication := config.GetDuplication()
	depth := config.GetDepth()

	pages := aws.ListObjects(bucket, region, prefix)
	for pages.Next() {
		page := pages.CurrentPage()
		for _, e := range page.CommonPrefixes {
			if recursive && depth >= 0 && depth > strings.Count(*e.Prefix, "/") {
				recursiveFunc(*e.Prefix)
			} else if !validateFilter(*e.Prefix) {
			} else if duplication && !duplicatedManager.isDuplicated(*e.Prefix) {
			} else {
				output <- sPrintCommonPrefixes(e, details)
			}
		}

		for _, e := range page.Contents {
			if *e.Key == *page.Prefix {
				continue
			}
			if !validateFilter(*e.Key) {
			} else if !validateDate(e.LastModified) {
			} else if duplication && !duplicatedManager.isDuplicated(*e.Key) {
			} else {
				output <- sPrintContents(details, e)
			}
		}
	}
}

func sPrintCommonPrefixes(e s3.CommonPrefix, details bool) string {
	if details {
		return fmt.Sprintln("d", *e.Prefix)
	} else {
		return fmt.Sprintln(*e.Prefix)
	}
}

var totalSize int64

func sPrintContents(details bool, e s3.Object) string {
	atomic.AddInt64(&totalSize, *e.Size)
	if details {
		size := fmt.Sprintf("%d", *e.Size)
		if len(size) < 10 {
			size = strings.Repeat(" ", 10-len(size)) + size
		}
		return fmt.Sprintf("- %s %s %s\n", (*e.LastModified).Format(time.RFC3339), size, *e.Key)
	} else {
		return fmt.Sprintln(*e.Key)
	}
}

type duplicatedMangerType struct {
	lock        sync.Mutex
	scannedPath map[string]struct{}
}

var duplicatedManager duplicatedMangerType

func (d *duplicatedMangerType) isDuplicated(path string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, ok := d.scannedPath[path]
	d.scannedPath[path] = struct{}{}
	return ok
}
