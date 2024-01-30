package utils

import (
	"net/mail"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type StringUtils struct{}

var StringUtilsInstance StringUtils

func (su StringUtils) ConvertToType(value string, toType reflect.Type) (any, error) {
	return su.ConvertToKind(value, toType.Kind())
}

func (su StringUtils) NewUUID() string {
	uuidStr := uuid.New().String()

	return strings.ReplaceAll(uuidStr, "-", "")
}

func (su StringUtils) StringFirstToken(s string, separator string) string {
	tokens := strings.Split(s, separator)
	len_tokens := len(tokens)

	if len_tokens <= 1 {
		return s
	}

	return tokens[0]
}

func (su StringUtils) OptimizeString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func (su StringUtils) ToggleCase(s string) string {
	result := make([]byte, 0)

	for _, v := range s {
		result = append(result, byte(v)^32)
	}

	return string(result)
}

func (su StringUtils) IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)

	return err == nil
}

func (su StringUtils) IsLogin(s string) bool {
	return su.IsEmail(s + "@com")
}

func (su StringUtils) ConvertToKind(value string, toKind reflect.Kind) (any, error) {
	switch toKind {
	case reflect.Uint:
		{
			v, err := strconv.ParseUint(value, 10, 64)

			return uint(v), err
		}
	case reflect.Uint8:
		{
			v, err := strconv.ParseUint(value, 10, 8)

			return uint8(v), err
		}
	case reflect.Uint16:
		{
			v, err := strconv.ParseUint(value, 10, 16)

			return uint16(v), err
		}
	case reflect.Uint32:
		{
			v, err := strconv.ParseUint(value, 10, 32)

			return uint32(v), err
		}
	case reflect.Uint64:
		{
			v, err := strconv.ParseUint(value, 10, 64)

			return uint64(v), err
		}
	case reflect.Int:
		{
			v, err := strconv.ParseInt(value, 10, 64)

			return int(v), err
		}
	case reflect.Int8:
		{
			v, err := strconv.ParseInt(value, 10, 8)

			return int8(v), err
		}
	case reflect.Int16:
		{
			v, err := strconv.ParseInt(value, 10, 16)

			return int16(v), err
		}
	case reflect.Int32:
		{
			v, err := strconv.ParseInt(value, 10, 32)

			return int32(v), err
		}
	case reflect.Int64:
		{
			v, err := strconv.ParseInt(value, 10, 64)

			return int64(v), err
		}
	case reflect.Float32:
		{
			v, err := strconv.ParseFloat(value, 32)

			return float32(v), err
		}
	case reflect.Float64:
		{
			v, err := strconv.ParseFloat(value, 64)

			return float64(v), err
		}
	case reflect.Complex64:
		{
			v, err := strconv.ParseComplex(value, 64)

			return complex64(v), err
		}
	case reflect.Complex128:
		{
			v, err := strconv.ParseComplex(value, 128)

			return complex128(v), err
		}
	case reflect.Bool:
		{
			return strconv.ParseBool(value)
		}
	}

	return value, nil
}
