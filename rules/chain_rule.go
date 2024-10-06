package rules

// ChainRule will check an entire list of rules, if none pass and the check value is not a float64, an InvalidCheckType error is returned
type ChainRule struct {
	Rules []Rule
}

func (c *ChainRule) GetScore(check interface{}) (bool, float64, error) {
	for _, r := range c.Rules {
		checkPass, result, err := r.GetScore(check)
		if err != nil {
			return false, 0, err
		}
		if checkPass {
			return checkPass, result, nil
		}
	}

	if checkVal, ok := check.(float64); ok {
		return false, checkVal, nil
	}
	return false, 0, InvalidCheckType
}

func NewChainRule(r ...Rule) Rule {
	return &ChainRule{
		Rules: r,
	}
}
