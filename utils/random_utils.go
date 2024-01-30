package utils

import "math/rand"

type RandomUtils struct{}

var RandomUtilsInstance RandomUtils

func (ru RandomUtils) RandomInt64(min int64, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func (ru RandomUtils) RandomInt32(min int32, max int32) int32 {
	return rand.Int31n(max-min) + min
}
