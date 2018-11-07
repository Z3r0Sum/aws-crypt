package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Z3r0Sum/aws-crypt/crypt"
	"github.com/Z3r0Sum/aws-crypt/utils"
	"github.com/Z3r0Sum/aws-crypt/version"
)

func main() {
	action := flag.String("action", "", "Either decrypt or encrypt")
	alias := flag.String("alias", "", "KMS Key Alias to use")
	file := flag.String("file", "", "File to decrypt/encrypt (e.g. /tmp/secret)")
	region := flag.String("region", "us-east-1", "AWS Region")
	ver := flag.Bool("version", false, "Displays Current Version of AWS Crypt")
	flag.Parse()

	if *ver {
		version.Print()
		os.Exit(0)
	}

	if *action == "" || *alias == "" || *file == "" || *region == "" {
		flag.Usage()
		os.Exit(128)
	}

	if err := utils.InputValidator(*action, *alias, *file, *region); err != nil {
		fmt.Println(err)
		os.Exit(128)
	}

	// Get Absolute file path
	fp, err := filepath.Abs(*file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	kmsFile := crypt.Properties{
		AwsRegion: *region,
		File:      fp,
		KmsAlias:  *alias,
	}

	switch {
	case *action == "decrypt":
		if err := kmsFile.Decrypter(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *action == "encrypt":
		if err := kmsFile.Encrypter(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unsupported action:", action)
		os.Exit(1)
	}

}
