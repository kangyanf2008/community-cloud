package utils

import "errors"

func Zipmap(a, b []string) (map[string]string, error) {

	if len(a) != len(b) {
		return nil, errors.New("zip: arguments must be of same length")
	}

	r := make(map[string]string, len(a))

	for i, e := range a {
		r[e] = b[i]
	}

	return r, nil
}
