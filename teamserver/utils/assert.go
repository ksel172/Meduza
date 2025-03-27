package utils

import "fmt"

func AssertNotNil(value any) {
	if value == nil {
		panic(fmt.Sprintf("assertion error: %v is nil", value))
	}
}

func AssertEquals(firstValue, secondValue any) {
	if firstValue != secondValue {
		panic(fmt.Sprintf("assertion error: %v is not equal to %v", firstValue, secondValue))
	}
}
