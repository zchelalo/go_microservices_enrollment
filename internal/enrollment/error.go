package enrollment

import (
	"errors"
	"fmt"
)

var ErrUserIdRequired = errors.New("user id is required")
var ErrCouseIdRequired = errors.New("course id is required")
var ErrStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	EnrollmentId string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("enrollment '%s' doesn't exist", e.EnrollmentId)
}
