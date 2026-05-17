// Package tablefmt provides helpers for rendering CLI tables.
package tablefmt

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	type args struct {
		headers []string
		rows    [][]string
	}
	tests := []struct {
		name         string
		args         args
		wantContains []string
		wantW        string
	}{
		{
			name: "renders header and row values",
			args: args{
				headers: []string{"name", "date", "venue"},
				rows: [][]string{
					{"Goose", "2024-07-04", "Red Rocks"},
				},
			},
			wantContains: []string{"name", "date", "venue", "Goose", "2024-07-04", "Red Rocks"},
		},
		{
			name: "renders only headers when rows are empty",
			args: args{
				headers: []string{"id", "title"},
				rows:    [][]string{},
			},
			wantContains: []string{"id", "title"},
		},
		{
			name: "renders empty output when headers and rows are empty",
			args: args{
				headers: []string{},
				rows:    [][]string{},
			},
			wantW: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			Render(w, tt.args.headers, tt.args.rows)

			gotW := w.String()
			if tt.wantW != "" && gotW != tt.wantW {
				t.Errorf("Render() = %v, want %v", gotW, tt.wantW)
			}

			for _, wantFragment := range tt.wantContains {
				if !strings.Contains(gotW, wantFragment) {
					t.Errorf("Render() output = %q, want fragment %q", gotW, wantFragment)
				}
			}

			if len(tt.args.headers) > 0 && gotW == "" {
				t.Errorf("Render() output is empty for non-empty headers")
			}
		})
	}
}
