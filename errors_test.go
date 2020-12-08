package mongoid_test

import (
	"mongoid"
	"mongoid/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ErrorTestModel struct {
	mongoid.Base
	ID mongoid.ObjectID `bson:"_id"`
}

var ErrorTestModels = mongoid.Register(&ErrorTestModel{})
var _ = Describe("Result", func() {
	Context(".Streaming()", func() {
		Context(".At()", func() {
			It("should panic InvalidOperation", func() {
				res, err := ErrorTestModels.Find()
				if err != nil {
					panic(err)
				}
				defer res.Close()
				res.Streaming()
				Expect(func() {
					// normally a panic event within mongoid would write to the log output,
					// but for this test a panic is required in order to pass,
					// so we temporarily mute logging here to keep test suite output clean
					log.WithMute(func() {
						res.At(0)
					})
				}).To(Panic())
			})
		})
	})
})
