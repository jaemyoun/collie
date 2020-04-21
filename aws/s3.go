package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"log"
	"os"
)

var s3service = map[s3.BucketLocationConstraint]*s3.S3{}

func GetBucketRegion(name string) s3.BucketLocationConstraint {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}
	ec2Region, err := ec2metadata.New(cfg).Region()
	if err != nil {
		ec2Region = "ap-northeast-1"
	}
	cfg.Region = ec2Region
	region, err := s3manager.GetBucketRegion(context.Background(), cfg, name, ec2Region)
	if err != nil {
		log.Fatalln("failed to get bucket location in getBucketRegion(),", name, err)
	} else if len(region) == 0 {
		log.Fatalln("invalid s3 bucket location")
	} else {
		return s3.BucketLocationConstraint(region)
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
