package main

import (
	"fmt"
	"strings"
)

func inStringArray(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func appendSingle(arr []string, str string) []string {
	if !inStringArray(arr, str) {
		arr = append(arr, str)
	}
	return arr
}

func stringJoin(array interface{}, seq string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", seq, -1)
}
