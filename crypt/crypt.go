package crypt

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// Properties associated with a file to decrypt and encrypt
type Properties struct {
	AwsRegion string
	File      string
	KmsAlias  string
}

// Encrypter - takes a KMS alias, file, and region to perform encryption of that
// file into the same directory with a .enc extension
func (p Properties) Encrypter() error {
	kmsClient := newClient(p.AwsRegion)

	content, err := ioutil.ReadFile(p.File)
	if err != nil {
		return err
	}

	input := &kms.EncryptInput{
		KeyId:     aws.String(p.KmsAlias),
		Plaintext: content,
	}

	result, err := kmsClient.Encrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr

		}
		return err
	}

	encryptedFilename := p.File + ".enc"
	if err := ioutil.WriteFile(encryptedFilename, result.CiphertextBlob, os.FileMode(0644)); err != nil {
		return err
	}

	fmt.Println("Wrote out encrypted file to:", encryptedFilename)
	return nil
}

// Decrypter decrypt the supplied file and print its contents to STDOUT
func (p Properties) Decrypter() error {
	kmsClient := newClient(p.AwsRegion)

	content, err := ioutil.ReadFile(p.File)
	if err != nil {
		return err
	}

	secret, err := decrypt(content, kmsClient)
	if err != nil {
		return err
	}

	fmt.Println(string(secret))
	return nil
}

func decrypt(f []byte, c *kms.KMS) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: f,
	}

	result, err := c.Decrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr
		}

		return nil, err
	}

	return result.Plaintext, nil
}

func newClient(region string) *kms.KMS {
	sess := session.Must(session.NewSession())
	return kms.New(sess, aws.NewConfig().WithRegion(region))
}
