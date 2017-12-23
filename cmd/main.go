package main

import (
	"formation/core"
	"formation/terraform_helpers"
)

func main() {
	var g terraform_helpers.InstanceState
	g.Attributes = make(map[string]string)
	g.Attributes["simple_field"] = "simple_value"
	g.Attributes["another_field"] = "another_value"

	var i core.InstanceStateParser
	i.Parse(&g)
}
