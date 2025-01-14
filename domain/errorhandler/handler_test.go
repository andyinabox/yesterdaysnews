package errorhandler

import (
	"context"
	"errors"
	"testing"
)

func TestMashalJSON(t *testing.T) {
	expected := `{"ErrTypeTest":["test error"]}`

	eh := New(context.Background(), &Config{})
	eh.Add("ErrTypeTest", errors.New("test error"))
	data, err := eh.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != expected {
		t.Errorf("json output does not match:\n%s\n%s\n", expected, string(data))
	}
}

func TestString(t *testing.T) {
	expected := `ErrTypeTest:
  1: "test error"
`

	eh := New(context.Background(), &Config{})
	eh.Add("ErrTypeTest", errors.New("test error"))
	result := eh.String()

	if result != expected {
		t.Errorf("json output does not match:\n%s\n%s\n", expected, result)
	}
}
