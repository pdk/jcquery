package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	decoder := json.NewDecoder(os.Stdin)

	// dumpTokens(decoder)

	keys, err := collectKeys(uniqueStrings{}, decoder, "")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%s\n", strings.Join(keys.values, "\n"))
}

func dumpTokens(d *json.Decoder) {

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

type uniqueStrings struct {
	present map[string]bool
	values  []string
}

func (u uniqueStrings) append(s string) uniqueStrings {
	if u.present == nil {
		u.present = map[string]bool{}
	}

	if !u.present[s] {
		u.present[s] = true
		u.values = append(u.values, s)
	}

	return u
}

func collectKeys(u uniqueStrings, d *json.Decoder, path string) (uniqueStrings, error) {

	t, err := d.Token()
	if err != nil {
		return u, err
	}

	switch t {
	case json.Delim('{'):
		return objectKeys(u, d, path)
	case json.Delim('['):
		return arrayKeys(u, d, path)
	default:
		// it's just a value, not a key
		return u, nil
	}
}

func objectKeys(u uniqueStrings, d *json.Decoder, path string) (uniqueStrings, error) {

	for {
		t, err := d.Token()
		if err != nil {
			return u, err
		}

		if t == json.Delim('}') {
			// end of the object
			return u, nil
		}

		// each item is a new key
		newPath := path + "/" + t.(string)
		u = u.append(newPath)

		// descend the tree looking for other objects
		u, err = collectKeys(u, d, newPath)
		if err != nil {
			return u, err
		}
	}
}

func arrayKeys(u uniqueStrings, d *json.Decoder, path string) (uniqueStrings, error) {

	path = path + "/"
	u = u.append(path)

	for {
		t, err := d.Token()
		if err != nil {
			return u, err
		}

		switch t {
		case json.Delim(']'):
			// end of the array
			return u, nil
		case json.Delim('{'):
			// start a new object
			u, err = objectKeys(u, d, path)
			if err != nil {
				return u, err
			}
		case json.Delim('['):
			// it's another array
			u, err = arrayKeys(u, d, path)
			if err != nil {
				return u, err
			}
		default:
			// it's just a scalar value
		}
	}
}
