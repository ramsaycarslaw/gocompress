// Copyright (C) 2020 Ramsay Carslaw

package dict

import (
	"strconv"
	"strings"
)

// DictionaryCompression implements a dictionary compression in golang
func DictionaryCompression(src string) (string, map[string]int) {
	s := strings.Split(src, " ")
	var m = make(map[string]int)
	var result string

	// create dictionary
	for i, v := range s {
		if m[v] == 0 {
			m[v] = i
		}
	}

	for _, v := range s {
		result += strconv.Itoa(m[v]) + " "
	}
	return result, m
}

// DictionaryDecompression the reverse of compression
func DictionaryDecompression(src string, m map[string]int) (string, error) {
	s := strings.Split(src, " ")
	var out string

	for _, v := range s {
		if v == "" {
			break
		}
		actual, err := strconv.Atoi(string(v))
		if err != nil {
			return "", err
		}
		for key, value := range m {
			if actual == value {
				out += key + " "
			}
		}
	}
	return out, nil
}
