package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// InputValidator Validates CLI input
func InputValidator(action string, alias string, file string, region string) error {
	if err := actionValidator(action); err != nil {
		return err
	}

	if err := aliasValidator(alias, region); err != nil {
		return err
	}

	if err := fileValidatior(file); err != nil {
		return err
	}

	return nil
}

func actionValidator(action string) error {
	validValues := []string{"decrypt", "encrypt"}

	for _, v := range validValues {
		if v == action {
			return nil
		}
	}
	return errors.New("Invalid action specified - accepted actions: decrypt or encrypt")
}

func aliasValidator(alias string, region string) error {
	if matched, _ := regexp.MatchString("alias/.+", alias); !matched {
		return errors.New("KMS Key Alias must be formatted 'alias/<alias name>'")
	}

	input := &kms.DescribeKeyInput{KeyId: aws.String(alias)}
	kmsClient := kms.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(region))

	_, err := kmsClient.DescribeKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if matched, _ := regexp.MatchString("NoCredentialProviders: no valid providers in chain. Deprecated.*", aerr.Error()); matched {
				fmt.Println("Could not locate valid AWS Credentials")
				os.Exit(126)
			}
			return aerr
		}
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func fileValidatior(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("Issues locating the file: %s, does it exist?\n", file)
		return err
	}
	return nil
}
