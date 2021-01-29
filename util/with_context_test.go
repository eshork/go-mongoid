package util

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContextWithContext()", func() {
	var leftContext, rightContext, withContext context.Context

	unresolvedScenario := func() {
		Describe("Done()", func() {
			It("returns not-nil", func() {
				Expect(withContext.Done()).ToNot(BeNil())
			})
		})
		Describe("Err()", func() {
			It("returns nil", func() {
				Expect(withContext.Err()).ToNot(HaveOccurred())
			})
		})
	}

	leftWinScenario := func() {
		Describe("Done()", func() {
			It("returns left side channel (not nil)", func() {
				Expect(withContext.Done()).To(Equal(leftContext.Done()))
			})
		})
		Describe("Err()", func() {
			It("returns leftContext.Err()", func() {
				Expect(withContext.Err()).To(HaveOccurred())
				Expect(withContext.Err()).To(Equal(leftContext.Err()))
			})
		})
	}
	rightWinScenario := func() {
		Describe("Done()", func() {
			It("returns right side channel (not nil)", func() {
				Expect(withContext.Done()).To(Equal(rightContext.Done()))
			})
		})
		Describe("Err()", func() {
			It("returns rightContext.Err()", func() {
				Expect(withContext.Err()).To(HaveOccurred())
				Expect(withContext.Err()).To(Equal(rightContext.Err()))
			})
		})
	}

	BeforeEach(func() {
		leftContext = context.Background()
		rightContext = context.Background()
	})

	JustBeforeEach(func() {
		withContext = ContextWithContext(leftContext, rightContext)
	})

	Describe("neither side is cancellable", func() {
		Describe("Done()", func() {
			It("returns nil", func() {
				Expect(leftContext.Done()).To(BeNil())
				Expect(rightContext.Done()).To(BeNil())
				Expect(withContext.Done()).To(BeNil())
			})
		})
		Describe("Err()", func() {
			It("returns nil", func() {
				Expect(leftContext.Done()).To(BeNil())
				Expect(rightContext.Done()).To(BeNil())
				Expect(withContext.Err()).To(BeNil())
			})
		})
	})

	Describe("both sides cancellable", func() {
		var leftCancel, rightCancel context.CancelFunc
		BeforeEach(func() {
			leftContext, leftCancel = context.WithCancel(context.Background())
			rightContext, rightCancel = context.WithCancel(context.Background())
		})
		AfterEach(func() {
			leftCancel()
			rightCancel()
		})
		Describe("before any cancel", func() {
			unresolvedScenario()
		})
		Describe("after left cancel", func() {
			BeforeEach(func() {
				leftCancel()
			})
			leftWinScenario()
			Describe("followed by right cancel", func() {
				BeforeEach(func() {
					rightCancel()
				})
				leftWinScenario()
			})
		})
		Describe("after right cancel", func() {
			BeforeEach(func() {
				rightCancel()
			})
			rightWinScenario()
			Describe("followed by left cancel", func() {
				BeforeEach(func() {
					leftCancel()
				})
				rightWinScenario()
			})
		})
	})

	Describe("only left side cancellable", func() {
		var leftCancel context.CancelFunc
		BeforeEach(func() {
			leftContext, leftCancel = context.WithCancel(context.Background())
		})
		AfterEach(func() {
			leftCancel()
		})
		Describe("before cancel", func() {
			unresolvedScenario()
		})
		Describe("after left cancel", func() {
			BeforeEach(func() {
				leftCancel()
			})
			leftWinScenario()
		})
	})

	Describe("only right side cancellable", func() {
		var rightCancel context.CancelFunc
		BeforeEach(func() {
			rightContext, rightCancel = context.WithCancel(context.Background())
		})
		AfterEach(func() {
			rightCancel()
		})
		Describe("before cancel", func() {
			unresolvedScenario()
		})
		Describe("after right cancel", func() {
			BeforeEach(func() {
				rightCancel()
			})
			rightWinScenario()
		})
	})

	Describe("silly example", func() {
		Describe("left side cancellable by proxy", func() {
			var leftCancel, leftParentCancel context.CancelFunc
			BeforeEach(func() {
				var leftParentContext context.Context
				leftParentContext, leftParentCancel = context.WithCancel(context.Background())
				leftContext, leftCancel = context.WithCancel(leftParentContext)
			})
			AfterEach(func() {
				leftParentCancel()
				leftCancel()
			})
			Describe("before cancel", func() {
				Specify("Done() returns not-nil", func() {
					Expect(withContext.Done()).ToNot(BeNil())
				})
				Specify("Err() returns nil", func() {
					Expect(withContext.Err()).To(BeNil())
				})
			})
			Describe("after leftParentContext cancel", func() {
				BeforeEach(func() {
					leftParentCancel()
				})
				Specify("Done() returns left side channel (not nil)", func() {
					Expect(withContext.Done()).To(Equal(leftContext.Done()))
				})
				Specify("Err() returns leftContext.Err()", func() {
					Expect(withContext.Err()).To(HaveOccurred())
					Expect(withContext.Err()).To(Equal(leftContext.Err()))
				})
			})
		})
	})

})
