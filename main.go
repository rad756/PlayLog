package main

import (
	"strings"
)

func main() {
	ui()
}

// checks if string contains a comma
func noComma(s string) bool {
	res := strings.Contains(s, ",")

	if res {
		return false
	} else {
		return true
	}
}
