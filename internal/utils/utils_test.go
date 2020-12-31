package utils

import (
	"testing"
	"time"
)

func Test_PrettyTimeAgo(t *testing.T) {

	cases := map[string]string{
		"1s":         "less than a minute ago",
		"30s":        "less than a minute ago",
		"1m08s":      "about 1 minute ago",
		"15m0s":      "about 15 minutes ago",
		"59m10s":     "about 59 minutes ago",
		"1h10m02s":   "about 1 hour ago",
		"15h0m01s":   "about 15 hours ago",
		"30h10m":     "about 1 day ago",
		"50h":        "about 2 days ago",
		"720h05m":    "about 1 month ago",
		"3000h10m":   "about 4 months ago",
		"8760h59m":   "about 1 year ago",
		"17601h59m":  "about 2 years ago",
		"262800h19m": "about 30 years ago",
	}

	for duration, expected := range cases {
		d, e := time.ParseDuration(duration)
		if e != nil {
			t.Errorf("failed to create a duration: %s", e)
		}

		fuzzy := PrettyTimeAgo(d)
		if fuzzy != expected {
			t.Errorf("unexpected fuzzy duration value: %s for %s", fuzzy, duration)
		}
	}
}

func Test_Pluralize(t *testing.T) {
	testCases := []struct {
		name   string
		word   string
		amount int
		want   string
	}{
		{
			name:   "singular",
			word:   "label",
			amount: 1,
			want:   "1 label",
		},
		{
			name:   "plural",
			word:   "label",
			amount: 3,
			want:   "3 labels",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := Pluralize(tC.amount, tC.word)
			if got != tC.want {
				t.Errorf("Pluralize() got = %s, want = %s", got, tC.want)
			}
		})
	}
}

func Test_PresentInStringSlice(t *testing.T) {
	testCases := []struct {
		name   string
		hay    []string
		needle string
		want   bool
	}{
		{"simple true", []string{"foo", "bar", "baz"}, "bar", true},
		{"simple false", []string{"foo", "bar", "baz"}, "qux", false},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := PresentInStringSlice(tC.hay, tC.needle)
			if got != tC.want {
				t.Errorf("PresentInStringSlice() got = %t, want = %t", got, tC.want)
			}
		})
	}
}
func Test_CommonElementsInStringSlice(t *testing.T) {
	testCases := []struct {
		name   string
		array1 []string
		array2 []string
		want   []string
	}{
		{
			name:   "simple no matching elements",
			array1: []string{"foo", "bar", "baz"},
			array2: []string{"qux", "quux", "quz"},
			want:   []string{},
		},
		{
			name:   "simple matching elements",
			array1: []string{"foo", "quux", "baz"},
			array2: []string{"qux", "quux", "baz"},
			want:   []string{"quux", "baz"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := CommonElementsInStringSlice(tC.array1, tC.array2)
			if len(got) != len(tC.want) {
				t.Errorf("CommonElementsInStringSlice() size of got (%d) and wanted (%d) arrays differ",
					len(got), len(tC.want))
			}
			for i := range got {
				if got[i] != tC.want[i] {
					t.Errorf("CommonElementsInStringSlice() got = %s, want = %s", got[i], tC.want[i])
				}
			}
		})
	}
}
