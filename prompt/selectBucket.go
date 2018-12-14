package prompt

import (
	"fmt"
	"github.com/jaemyoun/collie/config"
)

type selectBucketType struct {
	BucketName string
}

var selectBuckets = []selectBucketType{
	{BucketName: "(new bucket)"},
	{BucketName: "nxlog-streaming-ap-northeast-1"},
	{BucketName: "nxlog-streaming-ap-northeast-2"},
	{BucketName: "nxlog-standard-ap-northeast-1-json"},
	{BucketName: "nxlog-standard-ap-northeast-1-parquet"},
	{BucketName: "nxlog-standard-ap-northeast-2-json"},
	{BucketName: "nxlog-standard-ap-northeast-2-parquet"},
	{BucketName: "nxlog-v1-standard-ap-northeast-1-parquet"},
	{BucketName: "nxlog-v1-standard-ap-northeast-2-parquet"},
}

func DoSelectBucket() {
	p := newSelectDraft("BucketName")
	p.Items = selectBuckets
	p.Searcher = search(func(index int) string {
		return selectBuckets[index].BucketName
	})
	idx, _, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	result := selectBuckets[idx].BucketName
	if idx == 0 {
		result = inputByPrompt("input new bucket name:", func(input string) (err error) {
			if len(input) == 0 {
				err = fmt.Errorf("no length")
			}
			return
		})
	}
	config.ToogleBucket(result)
}
