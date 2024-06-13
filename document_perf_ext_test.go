package mongoid_test

import (
	"mongoid"

	// gofakeit "github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	"fmt"
)

type PerfSimpleExampleDocument struct {
	mongoid.Document
	IntField    int
	StringField string
}

var PerfSimpleExampleDocuments = mongoid.Collection(&PerfSimpleExampleDocument{})

// helper to provide better reporting for sub-second bench times
func BenchLoop(b Benchmarker, interations int, fn func()) {
	const constSecToMilliSecond = 1000
	const constSecToNanoSecond = 1000000000
	s := fmt.Sprintf("walltime per sample (%d interations/sample)", interations)
	runtime := b.Time(s, func() {
		for i := 0; i < interations; i++ {
			fn()
		}
	})
	b.RecordValue("millisec/iteration (sample avg)", runtime.Seconds()*constSecToMilliSecond/float64(interations))
	b.RecordValue("nanosec/iteration (sample avg)", runtime.Seconds()*constSecToNanoSecond/float64(interations))
}

var _ = Describe("DocumentBase", func() {
	Context("performance", func() {
		// `type ExampleDocument struct` is from document_test.go:~16
		Measure("ExampleDocuments.New() creation", func(b Benchmarker) {
			BenchLoop(b, 200, func() {
				ExampleDocuments.New()
			})
		}, 500)

		Measure("PerfSimpleExampleDocuments.New() creation", func(b Benchmarker) {
			BenchLoop(b, 200, func() {
				PerfSimpleExampleDocuments.New()
			})
		}, 500)

		Measure("ExampleDocuments.ToBson()", func(b Benchmarker) {
			testDoc := ExampleDocuments.New()
			BenchLoop(b, 200, func() {
				testDoc.ToBson()
			})
		}, 500)

		Measure("PerfSimpleExampleDocuments.ToBson()", func(b Benchmarker) {
			testDoc := PerfSimpleExampleDocuments.New()
			BenchLoop(b, 200, func() {
				testDoc.ToBson()
			})
		}, 500)
	})
})
