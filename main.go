package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sess = session.Must(session.NewSession())
var region = &aws.Config{Region: aws.String("us-east-1")}

func main() {
	iamChecker()
	s3Checker()
}

func iamChecker() {
	iamsvc := iam.New(sess, region)

	result, err := iamsvc.GetUser(&iam.GetUserInput{})

	if err != nil {
		fmt.Println("Failed to get user", err)
		return
	}

	fmt.Printf("%s\n", result.GoString())
}

func s3Checker() {
	s3svc := s3.New(sess, region)

	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf("%s\n", aws.StringValue(bucket.Name))
	}
}
