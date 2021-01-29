package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResultNotExpected", func() {
	It("behaves", func() {
		Expect(IsMongoidError(ResultNotExpected{})).To(BeTrue())
		Expect(IsMongoidError(&ResultNotExpected{})).To(BeTrue())
		Expect(IsResultNotExpected(ResultNotExpected{})).To(BeTrue())
		Expect(IsResultNotExpected(&ResultNotExpected{})).To(BeTrue())
		Expect(IsResultNotExpected(ErrResultNotExpected)).To(BeTrue())
	})
})
