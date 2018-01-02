package terraform_helpers_test

import (
	"github.com/hashicorp/terraform/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jmcgill/formation/terraform_helpers"
)

var _ = Describe("Terraform Helpers", func() {
	It("should sort numeric keys numerically", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"foo.#":  "value1",
				"foo.10": "value3",
				"foo.2":  "value2",
			},
		}

		expected := terraform_helpers.SortedInstanceState{
			Attributes: []terraform_helpers.KeyValue{
				{
					Key:   "foo.#",
					Value: "value1",
				},
				{
					Key:   "foo.2",
					Value: "value2",
				},
				{
					Key:   "foo.10",
					Value: "value3",
				},
			},
		}

		var wrappedState terraform_helpers.InstanceState
		wrappedState = terraform_helpers.InstanceState(state)
		sortedState := wrappedState.ToSorted()

		Expect(*sortedState).To(Equal(expected))
	})
})