package terraform_helpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terraform Helpers Suite")
}
