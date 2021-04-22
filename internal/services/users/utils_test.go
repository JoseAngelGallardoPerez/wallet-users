package users

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Confialink/wallet-users/internal/db/models"
)

var _ = Describe("users package", func() {
	Context("utils", func() {
		Context("ToTypedValue", func() {
			Context("we expect `string` result", func() {
				Context("we have a string value", func() {
					val := "simple string"
					When("we use known type - `string`", func() {
						It("should return string", func() {
							res := ToTypedValue(val, models.AttributeTypeString)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.String))
							Expect(res).Should(Equal(val))
						})
					})
					When("we use unknown type", func() {
						It("should return string", func() {
							res := ToTypedValue(val, "unknown type")
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.String))
							Expect(res).Should(Equal(val))
						})
					})
				})
			})

			Context("we expect `boolean` result", func() {
				Context("we use known type - `bool`", func() {
					When("we have a boolean value", func() {
						It("should return boolean `true`", func() {
							res := ToTypedValue("true", models.AttributeTypeBool)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Bool))
							Expect(res).Should(BeTrue())
						})
						It("should return boolean `false`", func() {
							res := ToTypedValue("false", models.AttributeTypeBool)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Bool))
							Expect(res).Should(BeFalse())
						})
					})
					When("we have a not boolean value", func() {
						It("should return boolean `false`", func() {
							res := ToTypedValue("random-string", models.AttributeTypeBool)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Bool))
							Expect(res).Should(BeFalse())
						})
					})
				})
			})

			Context("we expect `integer` result", func() {
				Context("we use known type - `int`", func() {
					When("we have a valid integer value", func() {
						It("should return int64 value", func() {
							res := ToTypedValue("35", models.AttributeTypeInt)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Int64))
							Expect(res).Should(Equal(int64(35)))
						})
					})
					When("we have a not `int` value", func() {
						It("should return zero value", func() {
							res := ToTypedValue("random-string", models.AttributeTypeInt)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Int64))
							Expect(res).Should(Equal(int64(0)))
						})
					})
				})
			})

			Context("we expect `float` result", func() {
				Context("we use known type - `float`", func() {
					When("we have a valid `float` value", func() {
						It("should return `float64` value", func() {
							res := ToTypedValue("35.67", models.AttributeTypeFloat)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Float64))
							Expect(res).Should(Equal(float64(35.67)))
						})
					})
					When("we have a not `float` value", func() {
						It("should return zero value", func() {
							res := ToTypedValue("random-string", models.AttributeTypeFloat)
							Expect(reflect.TypeOf(res).Kind()).Should(Equal(reflect.Float64))
							Expect(res).Should(Equal(float64(0)))
						})
					})
				})
			})
		})
	})
})
