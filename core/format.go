package core

import "strings"

func Format(name string) string {
	return strings.Replace(strings.Replace(name, "-", "_", -1), ".", "_", -1)
}