package utils

import "fmt"

type ErrorUtils struct{}

var ErrorUtilsInstance ErrorUtils

func (eu ErrorUtils) ErrorDetailedMessage(err error) string {
	errMsg := err.Error()

	if errMsg != "" {
		return fmt.Sprintf("%T: %s", err, err.Error())
	} else {
		return fmt.Sprintf("%T", err)
	}
}
