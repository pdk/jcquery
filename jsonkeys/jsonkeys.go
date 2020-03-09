package jsonkeys

import (
	"encoding/json"
	"io"
	"log"
	"strings"
)

// GetKeys returns all the key-paths in the JSON input stream. "/" is used as
// path delimiter, ala unix paths. If keys in the input have slashes, this won't
// work right.
func GetKeys(input io.Reader) ([]string, error) {

	decoder := json.NewDecoder(input)

	keys := uniqueStrings{}
	var err error

	for {
		keys, err = collectKeys(keys, decoder, "")

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		log.Printf("got non EOF: %#v", err)
		return nil, err
	}

	// remove non-leafs
	result := []string{}
	for _, p := range keys.values {
		if isParentPath(p, keys.values) {
			continue
		}
		result = append(result, p)
	}

	return result, nil
}

func isParentPath(path string, allPaths []string) bool {

	for _, p := range allPaths {
		if path != p && strings.HasPrefix(p, path) {
			return true
		}
	}

	return false
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
