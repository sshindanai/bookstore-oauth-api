package errors

import (
	"fmt"
	"strings"

	"github.com/sshindanai/bookstore-utils-go/resterrors"
)

const (
	ErrorNoData = "redis: nil"
)

func ParseError(err error, obj ...interface{}) *resterrors.RestErr {
	if strings.Contains(err.Error(), ErrorNoData) {
		return resterrors.NewUnauthorizedError(fmt.Sprintf("Access token doesn't existed for id %v", obj[0]))
	}

	return resterrors.NewInternalServerError("database error", err)
}
