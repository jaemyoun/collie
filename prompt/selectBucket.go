package prompt

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jaemyoun/collie/aws"
	"github.com/jaemyoun/collie/config"
	"strings"
)

type selectBucketItemType struct {
	BucketName string
}

func selectBucket() {
	items := getSelectBucketItems()
	p := newSelectDraft("BucketName")
	p.Items = items
	p.Searcher = search(func(index int) string {
		return items[index].BucketName
	})
	idx, _, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	result := items[idx].BucketName
	if idx == 0 {
		result = inputByPrompt("input new bucket name:", func(input string) (err error) {
			if len(input) == 0 {
				err = fmt.Errorf("no length")
			}
			return
		})
	}
	toggleBucket(result)
}

func getSelectBucketItems() []selectBucketItemType {
	items := make([]selectBucketItemType, 0)
	items = append(items, selectBucketItemType{BucketName: "(new bucket)"})
	buckets := config.Predefined.GetBuckets()
	for _, e := range buckets {
		items = append(items, selectBucketItemType{BucketName: e})
	}
	return items
}

func toggleBucket(newBucketName string) {
	if config.ContainBucket(config.BucketNameType(newBucketName)) {
		config.UnselectBucket(config.BucketNameType(newBucketName))
		fmt.Printf("The bucket (%s) was removed in the bucket list\n\n", newBucketName)
	} else {
		region := aws.GetBucketRegion(newBucketName)
		if len(region) == 0 {
			region = s3.BucketLocationConstraint(strings.TrimSpace(strings.ToLower(inputByPrompt(
				fmt.Sprintf("input the region of %s:", newBucketName), func(input string) error {
					input = strings.TrimSpace(strings.ToLower(input))
					switch s3.BucketLocationConstraint(input) {
					case s3.BucketLocationConstraintEu, s3.BucketLocationConstraintEuWest1,
						s3.BucketLocationConstraintUsWest1, s3.BucketLocationConstraintUsWest2,
						s3.BucketLocationConstraintApSouth1, s3.BucketLocationConstraintApSoutheast1,
						s3.BucketLocationConstraintApSoutheast2, s3.BucketLocationConstraintApNortheast1,
						s3.BucketLocationConstraintSaEast1, s3.BucketLocationConstraintCnNorth1,
						s3.BucketLocationConstraintEuCentral1, "ap-northeast-2":
						return nil
					default:
						return fmt.Errorf("invalid region")
					}
				}))))
		}
		config.SelectBucket(config.BucketNameType(newBucketName), region)
		fmt.Printf("The bucket (%s) was added in the bucket list\n", newBucketName)
	}
}

func getStatusForSelectBucket() (ret string) {
	title := "\U00002023 Selected S3 Buckets: "
	ret += fmt.Sprintf(title)
	bucketList := config.GetSelectedBucket()
	if len(bucketList) == 0 {
		ret += fmt.Sprintf("(nothing)\n")
	} else {
		first := true
		for bucket, info := range bucketList {
			if first == false {
				ret += fmt.Sprintf(",\n%s", strings.Repeat(" ", len(title)-2))
			}
			ret += fmt.Sprintf("%s (%s)", bucket, info.Region)
			first = false
		}
		ret += "\n"
	}
	return ret
}
