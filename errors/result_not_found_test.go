package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResultNotFound", func() {
	It("behaves", func() {
		Expect(IsMongoidError(ResultNotFound{})).To(BeTrue())
		Expect(IsMongoidError(&ResultNotFound{})).To(BeTrue())
		Expect(IsResultNotFound(ResultNotFound{})).To(BeTrue())
		Expect(IsResultNotFound(&ResultNotFound{})).To(BeTrue())
		Expect(IsResultNotFound(ErrResultNotFound)).To(BeTrue())
	})
})
