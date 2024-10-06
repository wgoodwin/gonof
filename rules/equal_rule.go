package rules

type EqualRule struct {
	Value  float64
	Result float64
}

func (e *EqualRule) GetScore(check interface{}) (bool, float64, error) {
	if checkVal, ok := check.(float64); ok {
		if checkVal == e.Value {
			return true, e.Result, nil
		}
		return false, checkVal, nil
	}
	return false, 0, InvalidCheckType
}

func NewEqualRule(val, res float64) Rule {
	return &EqualRule{
		Value:  val,
		Result: res,
	}
}
