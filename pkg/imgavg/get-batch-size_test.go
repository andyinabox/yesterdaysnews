package imgavg

import (
	"testing"
)

func TestGetBatchSize(t *testing.T) {

	{
		len := 510
		expected := 255
		result := getBatchSize(len)

		if result != expected {
			t.Errorf("expected %d got %d", expected, result)
		}
	}

	{
		len := 500
		expected := 250
		result := getBatchSize(len)

		if result != expected {
			t.Errorf("expected %d got %d", expected, result)
		}
	}

	{
		len := 700
		expected := 234
		result := getBatchSize(len)

		if result != expected {
			t.Errorf("expected %d got %d", expected, result)
		}
	}

}
