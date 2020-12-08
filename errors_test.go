package mongoid_test

import (
	"mongoid"
	// "github.com/brianvoe/gofakeit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ErrorTestModel1 struct {
	mongoid.Base
	ID mongoid.ObjectID `bson:"_id"`
}

var ErrorTestModel1s = mongoid.Register(&ErrorTestModel1{})
var _ = Describe("Result", func() {
	Context(".Streaming()", func() {
		Context(".At()", func() {
			It("should panic InvalidOperation", func() {
				res, err := ErrorTestModel1s.Find()
				if err != nil {
					panic(err)
				}
				defer res.Close()
				res.Streaming()
				Expect(func() {
					res.At(0)
				}).To(Panic())
			})
		})
	})
})
