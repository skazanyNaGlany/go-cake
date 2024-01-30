package utils

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type WalkIterableFunc func(value reflect.Value) bool
type ReflectUtils struct{}

var addNumericsError string = "Cannot add two numeric values %v + %v (%v)"
var ReflectUtilsInstance ReflectUtils

func (ru ReflectUtils) CallReflectMethod(method reflect.Value, args []any) []any {
	inputs := make([]reflect.Value, len(args))

	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	valued_results := method.Call(inputs)

	results := make([]any, len(valued_results))

	for j := range valued_results {
		results[j] = valued_results[j].Interface()
	}

	return results
}

func (ru ReflectUtils) CallReflectMethodNoArgs(method reflect.Value) []any {
	args := make([]any, 0)

	return ru.CallReflectMethod(method, args)
}

func (ru ReflectUtils) HasEmbeddedType(val any, _type reflect.Type) bool {
	valType := reflect.TypeOf(val)

	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}

	if valType.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < valType.NumField(); i++ {
		if valType.Field(i).Type == _type {
			return true
		}
	}

	return false
}

func (ru ReflectUtils) WalkIterable(v reflect.Value, targetFunc WalkIterableFunc, recursive bool) {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		if !targetFunc(v) {
			return
		}

		for i := 0; i < v.Len(); i++ {
			if !targetFunc(v.Index(i)) {
				break
			}

			ru.WalkIterable(v.Index(i), targetFunc, recursive)
		}
	case reflect.Map:
		if !targetFunc(v) {
			return
		}

		for _, k := range v.MapKeys() {
			if !targetFunc(v.MapIndex(k)) {
				break
			}

			ru.WalkIterable(v.MapIndex(k), targetFunc, recursive)
		}
	default:
		// handle other types
	}
}

func (ru ReflectUtils) GetValueLenght(value any) (int64, reflect.Type) {
	_type := reflect.TypeOf(value)
	lenght := int64(0)

	if value == nil {
		return lenght, _type
	}

	switch _type.Kind() {
	case reflect.Int:
		lenght = int64(value.(int))
	case reflect.Int8:
		lenght = int64(value.(int8))
	case reflect.Int16:
		lenght = int64(value.(int16))
	case reflect.Int32:
		lenght = int64(value.(int32))
	case reflect.Int64:
		lenght = int64(value.(int64))
	case reflect.Uint:
		lenght = int64(value.(uint))
	case reflect.Uint8:
		lenght = int64(value.(uint8))
	case reflect.Uint16:
		lenght = int64(value.(uint16))
	case reflect.Uint32:
		lenght = int64(value.(uint32))
	case reflect.Uint64:
		lenght = int64(value.(uint64))
	case reflect.Float32:
		lenght = int64(value.(float32))
	case reflect.Float64:
		lenght = int64(value.(float64))
	case reflect.Complex64:
		lenght = int64(real(value.(complex64)))
	case reflect.Complex128:
		lenght = int64(real(value.(complex128)))
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		funk.ForEach(value, func(v any) {
			lenght++
		})
	case reflect.Map:
		funk.ForEach(value, func(k, v any) {
			lenght++
		})
	case reflect.String:
		lenght = int64(len(value.(string)))
	default:
		return 0, nil
	}

	return lenght, _type
}

func (ru ReflectUtils) IsNumericType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return true
	}

	return false
}

func (ru ReflectUtils) IsIntType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return true
	}

	return false
}

func (ru ReflectUtils) IsUIntType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return true
	}

	return false
}

func (ru ReflectUtils) IsFloatType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return true
	}

	return false
}

func (ru ReflectUtils) IsComplexType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return true
	}

	return false
}

func (ru ReflectUtils) AddNumerics(value1, value2 any, _type reflect.Type) (any, error) {
	switch _type.Kind() {
	case reflect.Int:
		cval1, ok1 := value1.(int)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(int)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Int8:
		cval1, ok1 := value1.(int8)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(int8)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Int16:
		cval1, ok1 := value1.(int16)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(int16)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Int32:
		cval1, ok1 := value1.(int32)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(int32)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Int64:
		cval1, ok1 := value1.(int64)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(int64)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Uint:
		cval1, ok1 := value1.(uint)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(uint)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Uint8:
		cval1, ok1 := value1.(uint8)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(uint8)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Uint16:
		cval1, ok1 := value1.(uint16)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(uint16)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Uint32:
		cval1, ok1 := value1.(uint32)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(uint32)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Uint64:
		cval1, ok1 := value1.(uint64)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(uint64)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Float32:
		cval1, ok1 := value1.(float32)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(float32)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Float64:
		cval1, ok1 := value1.(float64)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(float64)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Complex64:
		cval1, ok1 := value1.(complex64)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(complex64)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	case reflect.Complex128:
		cval1, ok1 := value1.(complex128)

		if !ok1 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		cval2, ok2 := value2.(complex128)

		if !ok2 {
			return nil, errors.Errorf(addNumericsError, value1, value2, _type)
		}

		return cval1 + cval2, nil
	}

	return nil, errors.Errorf("Cannot add two numeric types %v + %v (%v)", value1, value2, _type)
}

func (ru ReflectUtils) IsArrayType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return true
	}

	return false
}

func (ru ReflectUtils) IsMapType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.Map:
		return true
	}

	return false
}

func (ru ReflectUtils) IsStringType(_type reflect.Type) bool {
	switch _type.Kind() {
	case reflect.String:
		return true
	}

	return false
}

func (ru ReflectUtils) IsIterableType(_type reflect.Type) bool {
	return ru.IsArrayType(_type) || ru.IsMapType(_type)
}

func (ru ReflectUtils) PtrToValue(ptr any) (any, error) {
	_value := reflect.ValueOf(ptr)

	if _value.Kind() != reflect.Ptr {
		return nil, errors.Errorf("%v is not a pointer", ptr)
	}

	return _value.Elem().Interface(), nil
}

func (ru ReflectUtils) GetFunctionName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
