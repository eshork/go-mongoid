package util

import (
	"math/cmplx"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetValueByInterfacePtr()", func() {
	// this test() is a bit complex due to the various scenarios it supports
	// - testStructPtr may be a pointer to a value or a pointer to pointer to a value, when a pointer to a pointer, the pointer it points to may be nil
	// - newVal may be a value or a pointer to a value, when a pointer, it may be nil
	// in all of these cases, it needs to compare values before and after calling SetValueByInterfacePtr() to verify the expected behavior
	test := func(testStructPtr interface{}, newVal interface{}) {
		handleStructValue := reflect.Indirect(reflect.ValueOf(testStructPtr))
		handleStructField := handleStructValue.Field(0)
		if handleStructField.Kind() == reflect.Ptr {
			if handleStructField.IsNil() {
				ExpectWithOffset(1, handleStructField.Interface()).ToNot(Equal(newVal), "testStructPtr.Field(nil) and newVal should not be equal before assignment")
			} else {
				ExpectWithOffset(1, reflect.Indirect(handleStructField).Interface()).ToNot(Equal(newVal), "testStructPtr.Field and newVal should not be equal before assignment")
			}
			SetValueByInterfacePtr(handleStructField.Addr().Interface(), newVal)
			if handleStructField.IsNil() {
				ExpectWithOffset(1, handleStructField.Interface()).To(BeNil(), "testStructPtr.Field(nil) and newVal should both be nil after nil assignment")
				ExpectWithOffset(1, newVal).To(BeNil(), "testStructPtr.Field(nil) and newVal should both be nil after nil assignment")
			} else {
				if reflect.ValueOf(newVal).Kind() == reflect.Ptr {
					ExpectWithOffset(1, reflect.Indirect(handleStructField).Interface()).To(Equal(reflect.Indirect(reflect.ValueOf(newVal)).Interface()), "*testStructPtr.Field and *newVal should be equal after assignment")
				} else {
					ExpectWithOffset(1, reflect.Indirect(handleStructField).Interface()).To(Equal(newVal), "testStructPtr.Field and newVal should be equal after assignment")
				}
			}
		} else {
			ExpectWithOffset(1, handleStructField.Interface()).ToNot(Equal(newVal), "src and dst should not be equal before assignment")
			SetValueByInterfacePtr(handleStructField.Addr().Interface(), newVal)
			ExpectWithOffset(1, handleStructField.Interface()).To(Equal(newVal), "src and dst should be equal after assignment")
		}
	}

	Context("given concrete dst", func() {
		It("bool", func() {
			tStruct := struct {
				Field bool
			}{}
			newVal := true
			test(&tStruct, newVal)
		})
		It("int", func() {
			tStruct := struct {
				Field int
			}{}
			newVal := 42
			test(&tStruct, newVal)
		})
		It("int8", func() {
			tStruct := struct {
				Field int8
			}{}
			newVal := int8(127)
			test(&tStruct, newVal)
		})
		It("int16", func() {
			tStruct := struct {
				Field int16
			}{}
			newVal := int16(32767)
			test(&tStruct, newVal)
		})
		It("int32", func() {
			tStruct := struct {
				Field int32
			}{}
			newVal := int32(2147483647)
			test(&tStruct, newVal)
		})
		It("int64", func() {
			tStruct := struct {
				Field int64
			}{}
			newVal := int64(4294967296)
			test(&tStruct, newVal)
		})
		It("uint", func() {
			tStruct := struct {
				Field uint
			}{}
			newVal := uint(42)
			test(&tStruct, newVal)
		})
		It("uint8", func() {
			tStruct := struct {
				Field uint8
			}{}
			newVal := uint8(255)
			test(&tStruct, newVal)
		})
		It("uint16", func() {
			tStruct := struct {
				Field uint16
			}{}
			newVal := uint16(65535)
			test(&tStruct, newVal)
		})
		It("uint32", func() {
			tStruct := struct {
				Field uint32
			}{}
			newVal := uint32(4294967295)
			test(&tStruct, newVal)
		})
		It("uint64", func() {
			tStruct := struct {
				Field uint64
			}{}
			newVal := uint64(18446744073709551615)
			test(&tStruct, newVal)
		})
		It("float32", func() {
			tStruct := struct {
				Field float32
			}{}
			newVal := float32(99.99)
			test(&tStruct, newVal)
		})
		It("float64", func() {
			tStruct := struct {
				Field float64
			}{}
			newVal := float64(-99.99)
			test(&tStruct, newVal)
		})
		It("complex64", func() {
			tStruct := struct {
				Field complex64
			}{}
			newVal := complex64(-99.99)
			test(&tStruct, newVal)
		})
		It("complex128", func() {
			tStruct := struct {
				Field complex128
			}{}
			newVal := cmplx.Sqrt(-1.0)
			test(&tStruct, newVal)
		})
		It("string", func() {
			tStruct := struct {
				Field string
			}{}
			newVal := "forty two"
			test(&tStruct, newVal)
		})

	})

	Context("given ptr to concrete dst", func() {
		It("*bool", func() {
			tStruct := struct {
				Field *bool
			}{}
			newVal := true
			test(&tStruct, newVal)
			newVal = false
			test(&tStruct, &newVal)
			test(&tStruct, true)
			test(&tStruct, nil)
		})
		It("*int", func() {
			tStruct := struct {
				Field *int
			}{}
			newVal := 42
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, 42)
			test(&tStruct, nil)
		})
		It("*int8", func() {
			tStruct := struct {
				Field *int8
			}{}
			newVal := int8(127)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, int8(127))
			test(&tStruct, nil)
		})
		It("*int16", func() {
			tStruct := struct {
				Field *int16
			}{}
			newVal := int16(32767)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, int16(32767))
			test(&tStruct, nil)
		})
		It("*int32", func() {
			tStruct := struct {
				Field *int32
			}{}
			newVal := int32(2147483647)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, int32(2147483647))
			test(&tStruct, nil)
		})
		It("*int64", func() {
			tStruct := struct {
				Field *int64
			}{}
			newVal := int64(4294967296)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, int64(4294967296))
			test(&tStruct, nil)
		})
		It("*uint", func() {
			tStruct := struct {
				Field *uint
			}{}
			newVal := uint(42)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, uint(42))
			test(&tStruct, nil)
		})
		It("*uint8", func() {
			tStruct := struct {
				Field *uint8
			}{}
			newVal := uint8(255)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, uint8(255))
			test(&tStruct, nil)
		})
		It("*uint16", func() {
			tStruct := struct {
				Field *uint16
			}{}
			newVal := uint16(65535)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, uint16(65535))
			test(&tStruct, nil)
		})
		It("*uint32", func() {
			tStruct := struct {
				Field *uint32
			}{}
			newVal := uint32(4294967295)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, uint32(4294967295))
			test(&tStruct, nil)
		})
		It("*uint64", func() {
			tStruct := struct {
				Field *uint64
			}{}
			newVal := uint64(18446744073709551615)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, uint64(18446744073709551615))
			test(&tStruct, nil)
		})
		It("*float32", func() {
			tStruct := struct {
				Field *float32
			}{}
			newVal := float32(99.99)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, float32(99.99))
			test(&tStruct, nil)
		})
		It("*float64", func() {
			tStruct := struct {
				Field *float64
			}{}
			newVal := float64(-99.99)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, float64(-99.99))
			test(&tStruct, nil)
		})
		It("*complex64", func() {
			tStruct := struct {
				Field *complex64
			}{}
			newVal := complex64(-99.99)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, complex64(-99.99))
			test(&tStruct, nil)
		})
		It("*complex128", func() {
			tStruct := struct {
				Field *complex128
			}{}
			newVal := cmplx.Sqrt(-1.0)
			test(&tStruct, newVal)
			newVal = 0
			test(&tStruct, &newVal)
			test(&tStruct, cmplx.Sqrt(-1.0))
			test(&tStruct, nil)
		})
		It("*string", func() {
			tStruct := struct {
				Field *string
			}{}
			newVal := "forty two"
			test(&tStruct, newVal)
			newVal = ""
			test(&tStruct, &newVal)
			test(&tStruct, "something else")
			test(&tStruct, nil)
		})
	})

	Context("structs", func() {
		It("simple assignment", func() {
			type innerStruct struct{ AnotherField int }
			tStruct := struct {
				Field innerStruct
			}{}
			newVal := innerStruct{AnotherField: 7}
			test(&tStruct, newVal)
		})
		It("pointer assignment", func() {
			type innerStruct struct{ AnotherField int }
			tStruct := struct {
				Field *innerStruct
			}{}
			newVal := innerStruct{AnotherField: 7}
			test(&tStruct, newVal)
			newVal = innerStruct{AnotherField: 42}
			test(&tStruct, &newVal)
			test(&tStruct, innerStruct{AnotherField: 0})
			test(&tStruct, nil)
		})
	})

	Context("map", func() {
		Context("map[string]int", func() {
			type innerMap map[string]int
			It("simple assignment", func() {
				tStruct := struct {
					Field innerMap
				}{}
				newVal := innerMap{"Example": 7}
				test(&tStruct, newVal)
			})
			It("pointer assignment", func() {
				tStruct := struct {
					Field *innerMap
				}{}
				newVal := innerMap{"Example": 7}
				test(&tStruct, newVal)
				newVal = innerMap{"DifferentExample": 42}
				test(&tStruct, &newVal)
				test(&tStruct, innerMap{"ExamplesAbound": 0})
				test(&tStruct, nil)
			})
		})
		Context("map[string]interface{} (bson.M)", func() {
			type innerMap map[string]interface{}
			It("simple assignment", func() {
				tStruct := struct {
					Field innerMap
				}{}
				newVal := innerMap{"Example": 7}
				test(&tStruct, newVal)
				newVal = innerMap{"ExampleTwo": "SomethingDifferent"}
				test(&tStruct, newVal)
			})
			It("pointer assignment", func() {
				tStruct := struct {
					Field *innerMap
				}{}
				newVal := innerMap{"Example": 7}
				test(&tStruct, newVal)
				newVal = innerMap{"DifferentExample": "42"}
				test(&tStruct, &newVal)
				test(&tStruct, innerMap{"ExamplesAbound": nil})
				test(&tStruct, nil)
			})
		})
	})

	Context("array", func() {
		Context("[]int", func() {
			type innerArray []int
			It("simple assignment", func() {
				tStruct := struct {
					Field innerArray
				}{}
				newVal := innerArray{7, 42}
				test(&tStruct, newVal)
			})
			It("pointer assignment", func() {
				tStruct := struct {
					Field *innerArray
				}{}
				newVal := innerArray{7, 42}
				test(&tStruct, newVal)
				newVal = innerArray{7, 42, 99}
				test(&tStruct, &newVal)
				test(&tStruct, innerArray{})
				test(&tStruct, nil)
			})
		})
		Context("[]interface{} (bson.A)", func() {
			type innerArray []interface{}
			It("simple assignment", func() {
				tStruct := struct {
					Field innerArray
				}{}
				newVal := innerArray{7, "42"}
				test(&tStruct, newVal)
			})
			It("pointer assignment", func() {
				tStruct := struct {
					Field *innerArray
				}{}
				newVal := innerArray{7, 42}
				test(&tStruct, newVal)
				newVal = innerArray{7, nil, 99}
				test(&tStruct, &newVal)
				test(&tStruct, innerArray{})
				test(&tStruct, nil)
			})
		})
	})

})
