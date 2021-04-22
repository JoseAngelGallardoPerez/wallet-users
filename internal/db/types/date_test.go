package types

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("date db type", func() {
	Context("Scan()", func() {
		When("value is nil", func() {
			It("returns nil", func() {
				dateType := &Date{}
				Expect(dateType.Scan(nil)).Should(BeNil())
			})
		})

		When("value is invalid", func() {
			It("returns an error", func() {
				dateType := &Date{}
				Expect(dateType.Scan("random string")).Should(MatchError("cannot Scan Date"))
			})
		})

		When("value is valid", func() {
			It("returns nil", func() {
				currentTime := time.Now()
				dateType := &Date{}
				Expect(dateType.Scan(currentTime)).Should(BeNil())
				Expect(dateType.Time).Should(Equal(&currentTime))
			})
		})
	})

	Context("String()", func() {
		When("struct is nil", func() {
			It("returns empty string", func() {
				var dateType *Date
				Expect(dateType.String()).Should(Equal(""))
			})
		})

		When("struct is not nil", func() {
			It("returns a formatted string", func() {
				currentTime := time.Now()
				dateType := &Date{Time: &currentTime}
				Expect(dateType.String()).Should(Equal(currentTime.Format("2006-01-02")))
			})
		})
	})

	Context("Value()", func() {
		When("struct is nil", func() {
			It("returns nil value and nil error", func() {
				var dateType *Date
				val, err := dateType.Value()
				Expect(val).Should(BeNil())
				Expect(err).Should(BeNil())
			})
		})

		When("struct is not nil", func() {
			It("returns a formatted string and nil error", func() {
				currentTime := time.Now()
				dateType := &Date{Time: &currentTime}
				val, err := dateType.Value()
				Expect(err).Should(BeNil())
				Expect(val).Should(Equal(currentTime.Format("2006-01-02")))
			})
		})
	})

	Context("UnmarshalJSON()", func() {
		When("data is not valid", func() {
			It("returns an error", func() {
				dateType := &Date{}
				data := []byte(`"11-11-11"`) // invalid date format
				Expect(dateType.UnmarshalJSON(data)).Should(HaveOccurred())
			})
		})

		When("data is valid", func() {
			It("returns a nil error", func() {
				currentTime := time.Now()
				dateType := &Date{}
				data := []byte(`"` + currentTime.Format("2006-01-02") + `"`)
				Expect(dateType.UnmarshalJSON(data)).ShouldNot(HaveOccurred())
				Expect(dateType.String()).Should(Equal(currentTime.Format("2006-01-02")))
			})
		})
	})

	Context("MarshalJSON()", func() {
		When("struct is nil", func() {
			It("returns nil data and nil error", func() {
				var dateType *Date
				val, err := dateType.MarshalJSON()
				Expect(val).Should(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("struct is not nil", func() {
			It("returns a bite slice and a nil error", func() {
				currentTime := time.Now()
				dateType := &Date{Time: &currentTime}
				expectedData := []byte(`"` + currentTime.Format("2006-01-02") + `"`)
				val, err := dateType.MarshalJSON()
				Expect(val).Should(Equal(expectedData))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
