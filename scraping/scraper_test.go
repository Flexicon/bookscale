package scraping

import "testing"

func TestParsePrice(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":           {input: "24,99 zł", want: "24,99 zł"},
		"with whitespace":  {input: " 32,50  zł  ", want: "32,50 zł"},
		"without currency": {input: " 42,23", want: "42,23 zł"},
		"empik example":    {input: "\n21,99 zł\n34,90 zł\n", want: "21,99 zł"},
		"empty":            {input: "", want: ""},
		"whitespace only":  {input: "       \t ", want: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parsePrice(tc.input)
			if got != tc.want {
				t.Errorf("expected: %s, got: %s", tc.want, got)
			}
		})
	}
}
