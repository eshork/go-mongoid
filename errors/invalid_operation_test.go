package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InvalidOperation", func() {
	It("behaves", func() {
		Expect(IsMongoidError(InvalidOperation{})).To(BeTrue())
		Expect(IsMongoidError(&InvalidOperation{})).To(BeTrue())
		Expect(IsInvalidOperation(InvalidOperation{})).To(BeTrue())
		Expect(IsInvalidOperation(&InvalidOperation{})).To(BeTrue())
	})
})
