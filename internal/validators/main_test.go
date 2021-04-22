package validators

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mocks "github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/go-playground/validator/v10"
)

var _ = Describe("validators package", func() {
	Context("specialCharacterRequired", func() {
		var (
			fl *mocks.FieldLevel
		)
		BeforeEach(func() {
			fl = &mocks.FieldLevel{}
		})
		Context("the value contains a special character", func() {
			specialCharacters := `!"#$%&'\()*+,-./:;<=>?@[\]^_{|}~`

			for _, character := range specialCharacters {
				stringChar := string(character)
				When(fmt.Sprintf("the value contains only `%s` character", stringChar), func() {
					It("should return true", func() {
						fl.On("Field").Return(reflect.ValueOf(stringChar))
						Expect(specialCharacterRequired(fl)).Should(BeTrue())
					})
				})

				When(fmt.Sprintf("the value contains the character `%s` in the start of the string", stringChar), func() {
					It("should return true", func() {
						fl.On("Field").Return(reflect.ValueOf(stringChar + "random"))
						Expect(specialCharacterRequired(fl)).Should(BeTrue())
					})
				})

				When(fmt.Sprintf("the value contains the character `%s` in the end of the string", stringChar), func() {
					It("should return true", func() {
						fl.On("Field").Return(reflect.ValueOf("random" + stringChar))
						Expect(specialCharacterRequired(fl)).Should(BeTrue())
					})
				})

				When(fmt.Sprintf("the value contains the character `%s` in the middle of the string", stringChar), func() {
					It("should return true", func() {
						fl.On("Field").Return(reflect.ValueOf("random" + stringChar + "random"))
						Expect(specialCharacterRequired(fl)).Should(BeTrue())
					})
				})
			}
		})

		When("the value does not contain a special character", func() {
			It("should return false", func() {
				fl.On("Field").Return(reflect.ValueOf("simpleString"))
				Expect(specialCharacterRequired(fl)).Should(BeFalse())
			})
		})

		When("the value is empty", func() {
			It("should return true", func() {
				fl.On("Field").Return(reflect.ValueOf(""))
				Expect(specialCharacterRequired(fl)).Should(BeTrue())
			})
		})
	})
})
