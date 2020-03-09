package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pdk/jcquery/jsonkeys"
)

func main() {

	keys, err := jsonkeys.GetKeys(os.Stdin)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%s\n", strings.Join(keys, "\n"))
}
