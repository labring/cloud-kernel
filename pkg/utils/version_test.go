package utils

import "testing"

func TestGetMajorMinorInt(t *testing.T) {
	a, b := GetMajorMinorInt("1.19.0")
	println(a)
	println(b)
}
