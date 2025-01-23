package youtubedownloader

import "testing"

func TestBinVersion(t *testing.T) {

	v1, err := ParseVersion("2024.12.06")
	if err != nil {
		t.Fatal(err)
	}
	v2 := NewVersion(2024, 11, 06)

	expectedV2str := "2024.11.06"
	v2str := v2.String()

	if v2str != expectedV2str {
		t.Errorf("expected %q, got %q", expectedV2str, v2str)
	}

	if v2.IsEqualOrGreaterThan(v1) {
		t.Errorf("expected %s to be less than %s", v1, v2)
	}

}
