package rules_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wgoodwin/gonof/rules"
)

var _ = Describe("ChainRule", func() {
	var rule rules.ChainRule
	var eqRule rules.EqualRule
	var ltRule rules.LTRule
	var gtRule rules.GTRule

	BeforeEach(func() {
		eqRule = rules.EqualRule{
			Value:  0,
			Result: -10,
		}
		ltRule = rules.LTRule{
			Value:  15,
			Result: 3,
		}
		gtRule = rules.GTRule{
			Value:  42,
			Result: 10,
		}
		rule = rules.ChainRule{
			Rules: []rules.Rule{
				&eqRule,
				&ltRule,
				&gtRule,
			},
		}
	})

	Context("GetScore", func() {
		It("should return result for equal", func() {
			check, res, err := rule.GetScore(float64(0))

			Expect(check).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(eqRule.Result))
		})

		It("should return result for lt", func() {
			check, res, err := rule.GetScore(float64(8))

			Expect(check).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(ltRule.Result))
		})

		It("should return result for gt", func() {
			check, res, err := rule.GetScore(float64(43))

			Expect(check).To(BeTrue())
			Expect(err).To(BeNil())
			Expect(res).To(Equal(gtRule.Result))
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
