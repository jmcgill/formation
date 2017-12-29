package core

import (
	"regexp"
	"strings"
)

func Format(name string) string {
	// Terraform identifiers can only contain letters, numbers, dashes and underscores

	// Make @ symboles human readable
	r := strings.Replace(name, "@", "_at_", -1)

	// Transform all other invalid characters into underscores
	re := regexp.MustCompile("[^A-Za-z0-9-_]")
	return string(re.ReplaceAll([]byte(r), []byte("_")))
}