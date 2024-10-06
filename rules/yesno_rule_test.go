package rules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/wgoodwin/gonof/rules"
)

var _ = Describe("YesnoRule", func() {
	var rule rules.YesNoRule

	BeforeEach(func() {
		rule = rules.YesNoRule{
			Yes: 42,
			No:  10,
		}
	})

	Context("GetScore", func() {
		It("should return yes value and true", func() {
			check, res, err := rule.GetScore("yes")

			Expect(check).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(rule.Yes))
		})

		It("should return no value and false", func() {
			check, res, err := rule.GetScore("no")

			Expect(check).To(BeFalse())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(rule.No))
		})

		It("should return false and invalid type error", func() {
			check, res, err := rule.GetScore(10)

			Expect(check).To(BeFalse())
			Expect(err).To(Equal(rules.InvalidCheckType))
			Expect(res).To(Equal(float64(0)))
		})

		It("should return false and invalid input error", func() {
			check, res, err := rule.GetScore("Bad")

			Expect(check).To(BeFalse())
			Expect(err).To(Equal(rules.InvalidInput))
			Expect(res).To(Equal(float64(0)))
		})
	})

})
