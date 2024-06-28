package validation

import (
	"fmt"
	"regexp"
)

type IValidation interface {
	ValidateUserName(string) ValidateResult
}

type Validation struct {
}

type ValidateResult struct {
	IsSuccess bool
	Message   string
}

func NewValidation() IValidation {
	return &Validation{}
}

func (v Validation) ValidateUserName(userName string) (result ValidateResult) {
	l := len(userName)
	if l < 3 || l > 20 {
		result.Message = fmt.Sprintf("The [%s] invalid length", userName)
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(userName) {
		result.Message = fmt.Sprintf("The [%s] contain invalid chars", userName)
		return
	}

	result.IsSuccess = true
	return
}
