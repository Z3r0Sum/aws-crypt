package version

import (
	"fmt"
)

// AwsCryptVersion is the current version of the CLI
var AwsCryptVersion = "v0.1.0"

// Print print the version to STDOUT
func Print() {
	fmt.Println("AWS Crypt:", AwsCryptVersion)
}
