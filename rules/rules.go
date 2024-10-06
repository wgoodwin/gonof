package rules

import "errors"

var (
	InvalidCheckType = errors.New("rules: invalid check type")
	InvalidInput     = errors.New("rules: invalid input")
)

type Rule interface {
	// GetScore takes in an interface and returns an evaluated score or an error if it couldn't be derived
	GetScore(interface{}) (bool, float64, error)
}
