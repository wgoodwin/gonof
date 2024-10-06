package rules

import "strings"

type YesNoRule struct {
	Yes float64
	No  float64
}

func (yn *YesNoRule) GetScore(check interface{}) (bool, float64, error) {
	if checkVal, ok := check.(string); ok {
		switch strings.ToLower(checkVal) {
		case "yes":
			return true, yn.Yes, nil
		case "no":
			return false, yn.No, nil
		default:
			return false, 0, InvalidInput
		}
	}
	return false, 0, InvalidCheckType
}

func NewYesNoRule(yes, no float64) Rule {
	return &YesNoRule{
		Yes: yes,
		No:  no,
	}
}
