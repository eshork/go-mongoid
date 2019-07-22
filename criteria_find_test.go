package mongoid

import (
	// "fmt"
	"go.mongodb.org/mongo-driver/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Criteria Building", func() {
	Describe("criteriaFind", func() {
		It("returns a new Criteria pointer", func() {
			Expect(criteriaFind(nil)).To(BeAssignableToTypeOf(&criteriaStruct{}))
		})
		It("accepts an empty criteria input", func() {
			criteriaFind(nil)
		})
		It("accepts an ObjectID criteria input", func() {
			criteriaFind(nil, NewObjectID())
		})
		It("accepts multiple ObjectID criteria input", func() {
			criteriaFind(nil, NewObjectID(), NewObjectID(), NewObjectID(), NewObjectID())
		})
		It("chains to given Criteria", func() {
			existingCriteria := criteriaStruct{}
			newCriteriaPtr := criteriaFind(&existingCriteria)
			Expect(newCriteriaPtr.getPrevCriteria()).To(BeIdenticalTo(&existingCriteria))
		})
		It("builds driver-ready query BSON", func() {
			{
				By("No ObjectIDs")
				criteriaPtr := criteriaFind(nil)
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaFindToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult)).To(Equal(0))
			}
			{
				By("One ObjectIDs")
				newObjID := NewObjectID()
				criteriaPtr := criteriaFind(nil, newObjID)
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaFindToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult)).To(Equal(1))
				Expect(bsonResult[0].Key).To(Equal("_id"))
				Expect(bsonResult[0].Value).To(Equal(newObjID))
			}
			{
				By("Many ObjectIDs")
				criteriaPtr := criteriaFind(nil, NewObjectID(), NewObjectID(), NewObjectID(), NewObjectID())
				findPtr := criteriaPtr.(*criteriaStruct)
				bsonResult := criteriaFindToBsonD(findPtr)
				Expect(bsonResult).To(BeAssignableToTypeOf(bson.D{}))
				Expect(len(bsonResult)).To(Equal(1))
				Expect(bsonResult[0].Key).To(Equal("$or"))
				Expect(bsonResult[0].Value).To(BeAssignableToTypeOf(bson.A{}))
				Expect(len(bsonResult[0].Value.(bson.A))).To(Equal(4))
			}
		})
	})
})

// func criteriaFind(prevCriteria *Criteria, ids ...ObjectID) *Criteria {
// 	thisQuery := Query{}
// 	newCriteria := Criteria{
// 		prevCriteria: prevCriteria,
// 		thisVerb:     queryVerbFind,
// 		thisQuery:    thisQuery,
// 	}
// 	return &newCriteria
// }
