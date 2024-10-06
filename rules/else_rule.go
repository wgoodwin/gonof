package rules

// ElseRule will always return the result value when GetScore is called
type ElseRule struct {
	Result float64
}

func (e *ElseRule) GetScore(interface{}) (bool, float64, error) {
	return true, e.Result, nil
}

func NewElseRule(res float64) Rule {
	return &ElseRule{
		Result: res,
	}
}
