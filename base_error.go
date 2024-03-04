package go_cake

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/skazanyNaGlany/go-cake/utils"
)

type BaseError struct {
	Message string
}

func (be BaseError) Error() string {
	return be.Message
}

func (be *BaseError) FormatStatusMessage(message string, baseError, internalError error) string {
	baseErrorType := utils.StructUtilsInstance.GetCleanType(baseError)
	baseErrorType = strings.Replace(baseErrorType, "Error", "", 1)

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

func (be *BaseError) logError(childErr error, object any) {
	message := fmt.Sprintf("HTTPError: %T\n", childErr)
	message += fmt.Sprintf("Stacktrace: %v\n", strings.TrimSpace(string(debug.Stack())))

	if object != nil {
		message += fmt.Sprintf("Object: %T (%v)\n", object, object)
	}

	log.Print(message)
}
