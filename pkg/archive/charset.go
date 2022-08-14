package archive

import "golang.org/x/text/encoding/ianaindex"

// String converts the given encoded string to UTF-8. It returns the converted
// string or "", err if any error occurred.
func DecodeString(src string, name string) (string, error) {
	encoding, err := ianaindex.IANA.Encoding(name)
	if err != nil {
		return "", err
	}
	target, err := encoding.NewDecoder().String(src)
	if err != nil {
		return "", err
	}
	return target, nil
}
