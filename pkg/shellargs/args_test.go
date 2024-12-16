package shellargs

import "testing"

func TestArgs(t *testing.T) {
	{
		expected := " -v --test 1 -something a/path/to.txt"

		a := New()
		a.Add("-v")
		a.AddKeyed("--test", "1")
		a.Add("-something")
		a.Add("a/path/to.txt")

		result := a.String()

		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}

	}
}
