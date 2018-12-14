package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

var s3service = map[s3.BucketLocationConstraint]*s3.S3{}

func GetBucketRegion(name string) s3.BucketLocationConstraint {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("failed to load config, %v\n", err)
		os.Exit(1)
	}
	req := s3.New(cfg).GetBucketLocationRequest(&s3.GetBucketLocationInput{Bucket: aws.String(name)})
	resp, err := req.Send()
	if err != nil {
		fmt.Println("failed to get bucket location in GetBucketRegion(),", name)
	} else if len(resp.LocationConstraint) == 0 {
		fmt.Println("invalid s3 bucket location")
	} else {
		return resp.LocationConstraint
	}
	return ""
}

func ListObjects(bucket, region, prefix string) s3.ListObjectsV2Pager {
	svc := getService(region)

	req := svc.ListObjectsV2Request(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	pages := req.Paginate()
	if err := pages.Err(); err != nil {
		log.Fatalln("failed to list objects,", err)
	}
	return pages
}

func getService(region string) *s3.S3 {
	svc, ok := s3service[s3.BucketLocationConstraint(region)]
	if ok == false {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			log.Fatalln("failed to load config,", err)
			os.Exit(1)
		}
		cfg.Region = region
		svc = s3.New(cfg)
		s3service[s3.BucketLocationConstraint(region)] = svc
	}
	return svc
}
