package builder

import (
	"reflect"
	"testing"
)

func TestGetPrefixesToDelete(t *testing.T) {

	prefixes := []string{"1536696082/", "1736696082/", "1236696082/", "current/", "1136696082/"}
	expected := []string{"1236696082/", "1136696082/", "current/"}

	result := getPrefixesToDelete(prefixes, 2)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected slices to be equal, got:\n%v\n%v\n", expected, result)
	}

}
