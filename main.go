package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

var creds = &aws.Config{}
var region = &aws.Config{Region: aws.String("us-east-1")}
var sess = session.Must(session.NewSession())

type Token struct {
	AwsAccessKeyId     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	AwsSecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
}

func main() {
	os.Exit(process(os.Args[1:]))
}

func process(args []string) int {
	var token Token
	if err := envconfig.Process("", &token); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	creds = &aws.Config{Credentials: credentials.NewStaticCredentials(token.AwsAccessKeyId, token.AwsSecretAccessKey, "")}

	iamChecker()
	s3Checker()

	return 0
}

func iamChecker() {
	iamsvc := iam.New(sess, creds, region)

	result, err := iamsvc.GetUser(&iam.GetUserInput{})

	if err != nil {
		fmt.Println("Failed to get user", err)
		return
	}

	fmt.Printf("%s\n", result.GoString())
}

func s3Checker() {
	s3svc := s3.New(sess, creds, region)

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
