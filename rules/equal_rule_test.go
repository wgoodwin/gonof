package rules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/wgoodwin/gonof/rules"
)

var _ = Describe("EqualRule", func() {
	var rule rules.EqualRule

	BeforeEach(func() {
		rule = rules.EqualRule{
			Value:  42,
			Result: 10,
		}
	})

	Context("GetScore", func() {
		It("should return result", func() {
			check, res, err := rule.GetScore(float64(42))

			Expect(check).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(rule.Result))
		})

		It("should return check value with false", func() {
			check, res, err := rule.GetScore(float64(40))

			Expect(check).To(BeFalse())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(float64(40)))
		})

		It("should return err and 0", func() {
			check, res, err := rule.GetScore("NaN")

			Expect(check).To(BeFalse())
			Expect(err).To(Equal(rules.InvalidCheckType))
			Expect(res).To(Equal(float64(0)))
		})
	})

})
