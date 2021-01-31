package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OperationTimedOut", func() {
	It("behaves", func() {
		Expect(IsMongoidError(OperationTimedOut{})).To(BeTrue())
		Expect(IsMongoidError(&OperationTimedOut{})).To(BeTrue())
		Expect(IsOperationTimedOut(OperationTimedOut{})).To(BeTrue())
		Expect(IsOperationTimedOut(&OperationTimedOut{})).To(BeTrue())
	})
})
