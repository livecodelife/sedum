// Command sedum generates boilerplate code from provenance records using
// generator packages that teams author themselves.
package main

import (
	"os"

	"github.com/calebcowen/sedum/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
