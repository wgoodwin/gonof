package rules

type LTRule struct {
	Value  float64
	Result float64
}

func (lt *LTRule) GetScore(check interface{}) (bool, float64, error) {
	if checkVal, ok := check.(float64); ok {
		if checkVal < lt.Value {
			return true, lt.Result, nil
		}
		return false, checkVal, nil
	}
	return false, 0, InvalidCheckType
}

func NewLTRule(val, res float64) Rule {
	return &LTRule{
		Value:  val,
		Result: res,
	}
}
