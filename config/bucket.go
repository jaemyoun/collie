package config

import "github.com/aws/aws-sdk-go-v2/service/s3"

func SelectBucket(bucketName BucketNameType, region s3.BucketLocationConstraint) {
	setting.selectedBucket[bucketName] = BucketInfo{Region: region}
}

func UnselectBucket(bucketName BucketNameType) {
	delete(setting.selectedBucket, bucketName)
}

func ContainBucket(bucketName BucketNameType) bool {
	_, ok := setting.selectedBucket[bucketName]
	return ok
}

func GetSelectedBucket() map[BucketNameType]BucketInfo {
	return setting.selectedBucket
}
