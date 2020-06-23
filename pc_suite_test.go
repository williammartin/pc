package pc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pc Suite")
}
