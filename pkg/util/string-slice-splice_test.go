package util

import (
	"reflect"
	"testing"
)

func TestStringSliceSplice(t *testing.T) {
	slice := []string{"one", "two", "three", "four", "five"}
	sliceCopy := make([]string, len(slice))
	copy(sliceCopy, slice)

	expectedSlice := []string{"one", "three", "four", "five"}
	expectedString := "two"

	str, result := StringSliceSplice(slice, 1)

	if str != expectedString {
		t.Errorf("expected %q, got %q", expectedString, str)
	}

	if !reflect.DeepEqual(result, expectedSlice) {
		t.Errorf("expected %v, got %v", expectedSlice, result)
	}

	if !reflect.DeepEqual(slice, sliceCopy) {
		t.Errorf("original slice was modified: %v", slice)
	}

}
