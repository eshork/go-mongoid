package mongoid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCriteria(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongoid")
}
