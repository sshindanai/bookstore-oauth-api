package errors

import (
	"fmt"
	"strings"
)

const (
	ErrorNoData = "redis: nil"
)

func ParseError(err error, obj ...interface{}) *RestErr {
	if strings.Contains(err.Error(), ErrorNoData) {
		return NewUnauthorizedError(fmt.Sprintf("Access token doesn't existed for id %v", obj[0]))
	}

	return NewInternalServerError(fmt.Sprintf(err.Error()))
}
