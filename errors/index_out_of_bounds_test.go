package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IndexOutOfBounds", func() {
	It("behaves", func() {
		Expect(IsMongoidError(IndexOutOfBounds{})).To(BeTrue())
		Expect(IsMongoidError(&IndexOutOfBounds{})).To(BeTrue())
		Expect(IsIndexOutOfBounds(IndexOutOfBounds{})).To(BeTrue())
		Expect(IsIndexOutOfBounds(&IndexOutOfBounds{})).To(BeTrue())
		Expect(IsIndexOutOfBounds(ErrIndexOutOfBounds)).To(BeTrue())
	})
})
