package mongoid

import (
	"math/cmplx"
	"reflect"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("structValuesFromBsonM()", func() {

	Context("updating types of struct field values by name", func() {
		test := func(structPtr interface{}, fieldPtr interface{}, exBson bson.M, bsonFieldName string) {
			fieldValue := reflect.Indirect(reflect.ValueOf(fieldPtr))
			_, ok := exBson[bsonFieldName]
			Expect(ok).To(Equal(true), "given bsonFieldName should be a valid key to the target value, so the test can validate successful assignment")
			Expect(fieldValue.Interface()).ToNot(Equal(exBson[bsonFieldName]), "initial struct field value should not already equal the target value of the test")
			structValuesFromBsonM(structPtr, exBson)
			Expect(fieldValue.Interface()).To(Equal(exBson[bsonFieldName]), "struct field value should equal the target value after assignment")
		}

		It("bool field", func() {
			boolFieldEx := struct{ BoolField bool }{true}
			test(&boolFieldEx, &boolFieldEx.BoolField, bson.M{"bool_field": false}, "bool_field")
		})

		It("int field", func() {
			intFieldEx := struct{ IntField int }{7}
			test(&intFieldEx, &intFieldEx.IntField, bson.M{"int_field": 42}, "int_field")
		})

		It("int8 field", func() {
			tinyIntFieldEx := struct{ TinyIntField int8 }{7}
			test(&tinyIntFieldEx, &tinyIntFieldEx.TinyIntField, bson.M{"tiny_int_field": int8(127)}, "tiny_int_field")
		})

		It("int16 field", func() {
			smallIntFieldEx := struct{ SmallIntField int16 }{7}
			test(&smallIntFieldEx, &smallIntFieldEx.SmallIntField, bson.M{"small_int_field": int16(32767)}, "small_int_field")
		})

		It("int32 field", func() {
			explicitIntFieldEx := struct{ AnIntField int32 }{7}
			test(&explicitIntFieldEx, &explicitIntFieldEx.AnIntField, bson.M{"an_int_field": int32(2147483647)}, "an_int_field")
		})

		It("int64 field", func() {
			bigIntFieldEx := struct{ BigIntField int64 }{7}
			test(&bigIntFieldEx, &bigIntFieldEx.BigIntField, bson.M{"big_int_field": int64(4294967296)}, "big_int_field")
		})

		It("uint field", func() {
			uintFieldEx := struct{ UIntField uint }{7}
			test(&uintFieldEx, &uintFieldEx.UIntField, bson.M{"u_int_field": uint(42)}, "u_int_field")
		})

		It("uint8 field", func() {
			tinyUIntFieldEx := struct{ TinyUIntField uint8 }{7}
			test(&tinyUIntFieldEx, &tinyUIntFieldEx.TinyUIntField, bson.M{"tiny_u_int_field": uint8(255)}, "tiny_u_int_field")
		})

		It("uint16 field", func() {
			smallUIntFieldEx := struct{ SmallUIntField uint16 }{7}
			test(&smallUIntFieldEx, &smallUIntFieldEx.SmallUIntField, bson.M{"small_u_int_field": uint16(65535)}, "small_u_int_field")
		})

		It("uint32 field", func() {
			explicitUIntFieldEx := struct{ AnUIntField uint32 }{7}
			test(&explicitUIntFieldEx, &explicitUIntFieldEx.AnUIntField, bson.M{"an_u_int_field": uint32(4294967295)}, "an_u_int_field")
		})

		It("uint64 field", func() {
			bigUIntFieldEx := struct{ BigUIntField uint64 }{7}
			test(&bigUIntFieldEx, &bigUIntFieldEx.BigUIntField, bson.M{"big_u_int_field": uint64(18446744073709551615)}, "big_u_int_field")
		})

		It("float32 field", func() {
			float32FieldEx := struct{ FloatField float32 }{0.0}
			test(&float32FieldEx, &float32FieldEx.FloatField, bson.M{"float_field": float32(99.99)}, "float_field")
		})

		It("float64 field", func() {
			float64FieldEx := struct{ Float64Field float64 }{0.0}
			test(&float64FieldEx, &float64FieldEx.Float64Field, bson.M{"float_64_field": float64(-99.99)}, "float_64_field")
		})

		It("complex64 field", func() {
			complex64FieldEx := struct{ Complex64Field complex64 }{0.0}
			test(&complex64FieldEx, &complex64FieldEx.Complex64Field, bson.M{"complex_64_field": complex64(-99.99)}, "complex_64_field")
		})

		It("complex128 field", func() {
			complex128FieldEx := struct{ Complex128Field complex128 }{0.0}
			test(&complex128FieldEx, &complex128FieldEx.Complex128Field, bson.M{"complex_128_field": cmplx.Sqrt(-1.0)}, "complex_128_field")
		})

		It("string field", func() {
			strFieldEx := struct{ StrField string }{"example"}
			test(&strFieldEx, &strFieldEx.StrField, bson.M{"str_field": "forty two"}, "str_field")
		})

		It("array field", func() {
			By("single type bson.A")
			arrayFieldEx := struct{ ArrayField bson.A }{ArrayField: bson.A{"array", "example"}}
			test(&arrayFieldEx, &arrayFieldEx.ArrayField, bson.M{"array_field": bson.A{"bar", "world", "3.14159"}}, "array_field")

			By("mixed type bson.A")
			test(&arrayFieldEx, &arrayFieldEx.ArrayField, bson.M{"array_field": bson.A{"bar", "world", 3.14159}}, "array_field")

			By("native []string")
			stringArrayFieldEx := struct{ StringArrayField []string }{StringArrayField: []string{"array", "example"}}
			test(&stringArrayFieldEx, &stringArrayFieldEx.StringArrayField, bson.M{"string_array_field": []string{"play", "nice"}}, "string_array_field")

			By("native []int")
			intArrayFieldEx := struct{ IntArrayField []int }{IntArrayField: []int{1, 2, 3}}
			test(&intArrayFieldEx, &intArrayFieldEx.IntArrayField, bson.M{"int_array_field": []int{42}}, "int_array_field")
		})

		It("map field", func() {
			By("bson.M with string values")
			mapFieldEx := struct{ MapField bson.M }{MapField: bson.M{"bson.M(ap)": "example", "equivalentTo": "map[string]interface{}"}}
			test(&mapFieldEx, &mapFieldEx.MapField, bson.M{"map_field": bson.M{"play": "nice"}}, "map_field")

			By("bson.M with mixed values")
			mapMixedFieldEx := struct{ MixedMapField bson.M }{MixedMapField: bson.M{"bson.M(ap)": "example", "equivalentTo": "map[string]interface{}", "soThisIsFine": 7, "thisToo": 42.0}}
			test(&mapMixedFieldEx, &mapMixedFieldEx.MixedMapField, bson.M{"mixed_map_field": bson.M{"integers": 99, "floats": 99.99, "strings": "oh my!"}}, "mixed_map_field")

			By("native map[string]string")
			mapStringFieldEx := struct{ StringMapField map[string]string }{StringMapField: map[string]string{"native": "map[string]string"}}
			test(&mapStringFieldEx, &mapStringFieldEx.StringMapField, bson.M{"string_map_field": map[string]string{"shouldBe": "wellSupported"}}, "string_map_field")
		})

		It("struct with embedded int field", func() {
			structFieldEx := struct{ StructField struct{ IntField int } }{StructField: struct{ IntField int }{IntField: 7}}
			exBson := bson.M{"struct_field": bson.M{"int_field": 42}}
			Expect(structFieldEx.StructField.IntField).ToNot(Equal(exBson["struct_field"].(bson.M)["int_field"]), "initial struct field value should not already equal the target value of the test")
			structValuesFromBsonM(&structFieldEx, exBson)
			Expect(structFieldEx.StructField.IntField).To(Equal(exBson["struct_field"].(bson.M)["int_field"]), "initial struct field value should not already equal the target value of the test")
		})

		It("struct with embedded string field", func() {
			structFieldEx := struct{ StructField struct{ StringField string } }{StructField: struct{ StringField string }{StringField: "foo"}}
			exBson := bson.M{"struct_field": bson.M{"string_field": "bar"}}
			Expect(structFieldEx.StructField.StringField).ToNot(Equal(exBson["struct_field"].(bson.M)["string_field"]), "initial struct field value should not already equal the target value of the test")
			structValuesFromBsonM(&structFieldEx, exBson)
			Expect(structFieldEx.StructField.StringField).To(Equal(exBson["struct_field"].(bson.M)["string_field"]), "initial struct field value should not already equal the target value of the test")
		})

	}) // Context("updating a struct field value", func() {

})
