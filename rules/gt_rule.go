package rules

type GTRule struct {
	Value  float64
	Result float64
}

func (gt *GTRule) GetScore(check interface{}) (bool, float64, error) {
	if checkVal, ok := check.(float64); ok {
		if checkVal > gt.Value {
			return true, gt.Result, nil
		}
		return false, checkVal, nil
	}
	return false, 0, InvalidCheckType
}

func NewGTRule(val, res float64) Rule {
	return &GTRule{
		Value:  val,
		Result: res,
	}
}
