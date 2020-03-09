package jsonkeys

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// DumpTokens is used to dump all the tokens for debugging.
func DumpTokens(d *json.Decoder) {

	for {
		t, err := d.Token()

		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case t == json.Delim('{'):
			fmt.Printf("{\n")
			continue
		case t == json.Delim('}'):
			fmt.Printf("}\n")
			continue
		case t == json.Delim('['):
			fmt.Printf("[\n")
			continue
		case t == json.Delim(']'):
			fmt.Printf("]\n")
			continue
		}

		fmt.Printf("%#v\n", t)
	}
}
