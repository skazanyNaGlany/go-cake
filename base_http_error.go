package go_cake

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/skazanyNaGlany/go-cake/utils"
)

type BaseHTTPError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

func (e BaseHTTPError) Error() string {
	return e.StatusMessage
}

func (e BaseHTTPError) GetStatusCode() int {
	return e.StatusCode
}

func (e BaseHTTPError) GetStatusMessage() string {
	return e.StatusMessage
}

func (e BaseHTTPError) logError(childErr HTTPError, object any) {
	message := fmt.Sprintf("HTTPError: %T\n", childErr)
	message += fmt.Sprintf("Stacktrace: %v\n", strings.TrimSpace(string(debug.Stack())))
	message += fmt.Sprintf("StatusCode: %v\n", childErr.GetStatusCode())
	message += fmt.Sprintf("StatusMessage: %v\n", childErr.GetStatusMessage())

	if object != nil {
		message += fmt.Sprintf("Object: %T (%v)\n", object, object)
	}

	log.Print(message)
}

func (e *BaseHTTPError) FormatStatusMessage(message string, baseError, internalError error) string {
	baseErrorType := utils.StructUtilsInstance.GetCleanType(baseError)
	baseErrorType = strings.Replace(baseErrorType, "HTTPError", "", 1)

	formatted := fmt.Sprintf("%s: %s", baseErrorType, message)
	formatted = strings.TrimSpace(formatted)

	if internalError != nil {
		internalErrorMessage := strings.TrimSpace(internalError.Error())

		if internalErrorMessage != "" {
			if strings.LastIndex(formatted, ":") != len(formatted)-1 {
				formatted += ":"
			}

			formatted += fmt.Sprintf(" %s", internalErrorMessage)
		}
	}

	formatted = strings.TrimSpace(formatted)

	return strings.TrimSuffix(formatted, ":")
}
