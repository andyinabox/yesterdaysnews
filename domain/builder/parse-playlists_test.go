package builder

import (
	"reflect"
	"testing"

	"gitlab.com/andyinabox/yesterdaysnews/domain"
)

func TestParsePlaylists(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []domain.PlaylistSource
	}{
		{
			name: "bare ids",
			in:   "ABC,DEF",
			want: []domain.PlaylistSource{
				{Name: "ABC", ID: "ABC"},
				{Name: "DEF", ID: "DEF"},
			},
		},
		{
			name: "named pairs",
			in:   "cbs:ABC,nbc:DEF",
			want: []domain.PlaylistSource{
				{Name: "cbs", ID: "ABC"},
				{Name: "nbc", ID: "DEF"},
			},
		},
		{
			name: "mixed bare and named",
			in:   "ABC,nbc:DEF",
			want: []domain.PlaylistSource{
				{Name: "ABC", ID: "ABC"},
				{Name: "nbc", ID: "DEF"},
			},
		},
		{
			name: "whitespace and empty entries",
			in:   " cbs : ABC , , nbc:DEF ,",
			want: []domain.PlaylistSource{
				{Name: "cbs", ID: "ABC"},
				{Name: "nbc", ID: "DEF"},
			},
		},
		{
			name: "missing name falls back to id",
			in:   ":ABC",
			want: []domain.PlaylistSource{
				{Name: "ABC", ID: "ABC"},
			},
		},
		{
			name: "empty id is skipped",
			in:   "cbs:,nbc:DEF",
			want: []domain.PlaylistSource{
				{Name: "nbc", ID: "DEF"},
			},
		},
		{
			name: "empty string yields empty slice",
			in:   "",
			want: []domain.PlaylistSource{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ParsePlaylists(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParsePlaylists(%q) = %#v, want %#v", tc.in, got, tc.want)
			}
		})
	}
}
