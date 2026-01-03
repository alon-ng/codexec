package db

import (
	"fmt"
	"strings"
)

func IsDuplicateKeyError(err error) bool {
	return strings.HasSuffix(err.Error(), "(SQLSTATE 23505)")
}

func IsDuplicateKeyErrorWithConstraint(err error, constraint string) bool {
	return strings.HasSuffix(err.Error(), fmt.Sprintf("\"%s\" (SQLSTATE 23505)", constraint))
}
