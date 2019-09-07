package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/url"
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
	if len(args) == 2 {
		token.AwsAccessKeyId = args[0]
		token.AwsSecretAccessKey = args[1]
	} else {
		if err := envconfig.Process("", &token); err != nil {
			log.Printf("[ERROR] Failed to process env var: %s", err)
			return 1
		}
	}

	creds = &aws.Config{Credentials: credentials.NewStaticCredentials(token.AwsAccessKeyId, token.AwsSecretAccessKey, "")}

	iamChecker()

	return 0
}

func iamChecker() {
	iamsvc := iam.New(sess, creds, region)

	user, err := iamsvc.GetUser(&iam.GetUserInput{})
	if err != nil {
		fmt.Println("Failed to get user", err)
		return
	}
	fmt.Println("User:")
	fmt.Printf("%s\n", user.GoString())

	userGroups, err := iamsvc.ListGroupsForUser(&iam.ListGroupsForUserInput{UserName: user.User.UserName})
	if err != nil {
		fmt.Println("Failed to get user groups", err)
		return
	}
	fmt.Println("User groups:")
	fmt.Printf("%s\n", userGroups.Groups)

	for _, group := range userGroups.Groups {
		groupPolicies, err := iamsvc.ListGroupPolicies(&iam.ListGroupPoliciesInput{GroupName: group.GroupName})
		if err != nil {
			fmt.Println("Failed to get group policies", err)
			return
		}
		fmt.Printf("Group %s policies:\n", *group.GroupName)
		fmt.Printf("%s\n", groupPolicies.String())

		groupAttachedPolicies, err := iamsvc.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{GroupName: group.GroupName})
		if err != nil {
			fmt.Println("Failed to get group attached policies", err)
			return
		}
		fmt.Printf("Group %s attached policies:\n", *group.GroupName)
		fmt.Printf("%s\n", groupAttachedPolicies.String())

		for _, policy := range groupAttachedPolicies.AttachedPolicies {
			fetchedPolicy, err := iamsvc.GetPolicy(&iam.GetPolicyInput{PolicyArn: policy.PolicyArn})
			if err != nil {
				fmt.Println("Failed to get policies", err)
				return
			}
			fmt.Printf("Policy for %s permissions:\n", *policy.PolicyName)
			version, _ := iamsvc.GetPolicyVersion(&iam.GetPolicyVersionInput{PolicyArn: policy.PolicyArn, VersionId: fetchedPolicy.Policy.DefaultVersionId})
			document, _ := parsePolicyDocument(version.PolicyVersion.Document)
			fmt.Printf("version %s\n", document)
		}

	}
}

func parsePolicyDocument(p *string) (string, error) {
	parsedDoc, _ := url.QueryUnescape(*p)
	return parsedDoc, nil
}
