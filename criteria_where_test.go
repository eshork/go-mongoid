package mongoid

import (
	// "fmt"
	"go.mongodb.org/mongo-driver/bson"
	// "strconv"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Criteria Building", func() {
	gofakeit.Seed(GinkgoRandomSeed())

	Describe("criteriaWhere", func() {
		It("returns a new Criteria pointer", func() {
			Expect(criteriaWhere(nil, nil)).To(BeAssignableToTypeOf(&criteriaStruct{}))
		})
		It("accepts a missing criteria input", func() {
			criteriaWhere(nil, nil)
		})
		It("accepts an empty Query{} input", func() {
			criteriaWhere(nil, nil, Query{})
		})
		It("accepts multiple blank criteria input", func() {
			criteriaWhere(nil, nil, Query{}, Query{})
		})
		It("chains to given Criteria", func() {
			existingCriteria := criteriaStruct{}
			newCriteriaPtr := criteriaWhere(nil, &existingCriteria, Query{})
			Expect(newCriteriaPtr.getPrevCriteria()).To(BeIdenticalTo(&existingCriteria))
		})
		It("builds driver-ready query BSON", func() {
			{
				By("No Query")
				criteriaPtr := criteriaWhere(nil, nil)
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaWhereToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult)).To(Equal(0))
			}
			{
				By("Empty Query")
				criteriaPtr := criteriaWhere(nil, nil, Q{})
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaWhereToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult)).To(Equal(0))
			}
			{
				By("With Query")
				criteriaPtr := criteriaWhere(nil, nil, Q{
					"Title.$eq":     "something",
					"Severity.$gte": 7,
					"Private.$ne":   true,
				})
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaWhereToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult) > 0).To(BeTrue(), "len(bsonResult) > 0")
			}
		})
	})

	Describe("normalizeQueryElement", func() {
		fieldKeyTypes := map[string]string{
			"top level field name (field)":                             "field",
			"nested field name (field.nestedField)":                    "field.nestedField",
			"double nested field name (field.nestedField.nestedField)": "field.nestedField.nestedField",
		}
		fieldValueTypes := map[string]interface{}{
			"string field value":  gofakeit.HipsterWord(),
			"int field value":     int(gofakeit.Int32()),
			"float64 field value": gofakeit.Float64(),
			"date field value":    gofakeit.Date(),
			"bool field value":    gofakeit.Bool(),
			"map field value":     map[string]string{gofakeit.Word(): gofakeit.HipsterWord()},
			"array field value":   []string{gofakeit.Word(), gofakeit.HipsterWord()},
		}

		standardOperatorTests := func(operator, fieldNameDesc, fieldName, fieldValueDesc string, fieldValue interface{}) {
			Context("."+operator, func() {
				fieldNameOperator := fieldName + "." + operator
				It("returns valid bson.D", func() {
					ret := normalizeQueryElement(fieldNameOperator, fieldValue)
					Expect(ret).To(BeAssignableToTypeOf(bson.D{}))
				})
				It("extracts trailing ."+operator+" from field name \""+fieldName+"."+operator+"\"", func() {
					ret := normalizeQueryElement(fieldNameOperator, fieldValue)
					Expect(ret[0].Key).To(Equal(fieldName))
				})
				It("embeds the original value into a map[operator]value", func() {
					ret := normalizeQueryElement(fieldNameOperator, fieldValue)
					By("Replacing the field value with a map")
					Expect(ret[0].Value).To(BeAssignableToTypeOf(bson.M{}))
					By("Using the operator as the map key")
					newFieldValueMap := ret[0].Value.(bson.M)
					val, ok := newFieldValueMap[operator]
					Expect(len(newFieldValueMap)).To(Equal(1), "newFieldValueMap only has 1 key")
					Expect(ok).To(BeTrue(), "newFieldValueMap key is the operator "+operator)
					By("Using the original field value as the map value")
					Expect(val).To(Equal(fieldValue))
				})
			})
		}

		fieldNameValueTableTests := func(fieldNameDesc, fieldName, fieldValueDesc string, fieldValue interface{}) {
			Context(fieldNameDesc, func() {
				Context(fieldNameDesc, func() {
					Context("with "+fieldValueDesc, func() {
						standardOperatorTests("$eq", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$ne", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$gt", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$gte", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$lt", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$lte", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$in", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$nin", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)
						standardOperatorTests("$exists", fieldNameDesc, fieldName, fieldValueDesc, fieldValue)

						Context("No operator", func() {
							It("returns the unaltered field name \""+fieldName+"\"", func() {
								ret := normalizeQueryElement(fieldName, fieldValue)
								Expect(ret[0].Key).To(Equal(fieldName))
							})
						})

					})
				})
			})
		}

		for fieldKeyContext, fieldKeyExampleValue := range fieldKeyTypes {
			for fieldValueContext, fieldValueExampleValue := range fieldValueTypes {
				fieldNameValueTableTests(fieldKeyContext, fieldKeyExampleValue, fieldValueContext, fieldValueExampleValue)
			}
		}

	})

})
