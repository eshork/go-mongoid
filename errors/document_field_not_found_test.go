package errors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DocumentFieldNotFound", func() {
	It("behaves", func() {
		Expect(IsMongoidError(DocumentFieldNotFound{})).To(BeTrue())
		Expect(IsMongoidError(&DocumentFieldNotFound{})).To(BeTrue())
		Expect(IsDocumentFieldNotFound(DocumentFieldNotFound{})).To(BeTrue())
		Expect(IsDocumentFieldNotFound(&DocumentFieldNotFound{})).To(BeTrue())
	})
})
