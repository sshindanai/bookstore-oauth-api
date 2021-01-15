package errors

import (
	"fmt"
	"strings"
)

const (
	ErrorNoData = "redis: nil"
)

func ParseError(err error) *RestErr {
	if strings.Contains(err.Error(), ErrorNoData) {
		return NewUnauthorizedError("Access token doesn't existed")
	}

	return NewInternalServerError(fmt.Sprintf(err.Error()))
}
