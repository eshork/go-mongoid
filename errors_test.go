package mongoid_test

import (
	"mongoid"
	"mongoid/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Result", func() {
	type ErrorTestModel struct {
		mongoid.Document
		ID mongoid.ObjectID `bson:"_id"`
	}
	var ErrorTestModels = mongoid.Collection(&ErrorTestModel{})
	Context(".Streaming()", func() {
		Context(".At()", func() {
			It("should panic InvalidOperation", func() {
				OnlineDatabaseOnly(func() {
					res := ErrorTestModels.Find()
					res.Streaming()
					Expect(func() {
						// normally a panic event within mongoid would write to the log output,
						// but for this test a panic is expected (required) in order to pass,
						// so we temporarily mute logging here to keep test suite output clean(er)
						log.WithMute(func() {
							res.At(0)
						})
					}).To(Panic())
				})
			})
		})
	})
})
