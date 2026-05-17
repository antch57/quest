package jamz

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/antch57/quest/internal/jamz/jambase"
	"github.com/urfave/cli/v3"
)

type mockShowSearcher struct {
	events []jambase.Event
	err    error
}

func (m mockShowSearcher) SearchShows(_ context.Context, _ SearchOptions) ([]jambase.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func TestSearchCmd(t *testing.T) {
	tests := []struct {
		name        string
		wantName    string
		wantUsage   string
		wantFlags   int
		wantCountry string
		wantRadius  int
		wantAction  bool
	}{
		{
			name:        "returns configured search command",
			wantName:    "search",
			wantUsage:   "search for upcoming shows from Jambase...",
			wantFlags:   7,
			wantCountry: "US",
			wantRadius:  25,
			wantAction:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchCmd()
			if got == nil {
				t.Fatalf("SearchCmd() = nil")
			}

			if got.Name != tt.wantName {
				t.Errorf("SearchCmd().Name = %q, want %q", got.Name, tt.wantName)
			}

			if got.Usage != tt.wantUsage {
				t.Errorf("SearchCmd().Usage = %q, want %q", got.Usage, tt.wantUsage)
			}

			if len(got.Flags) != tt.wantFlags {
				t.Errorf("SearchCmd().Flags length = %d, want %d", len(got.Flags), tt.wantFlags)
			}

			if (got.Action != nil) != tt.wantAction {
				t.Errorf("SearchCmd().Action present = %v, want %v", got.Action != nil, tt.wantAction)
			}

			var countryValue string
			var radiusValue int
			for _, flag := range got.Flags {
				switch f := flag.(type) {
				case *cli.StringFlag:
					if f.Name == "country" {
						countryValue = f.Value
					}
				case *cli.IntFlag:
					if f.Name == "radius" {
						radiusValue = f.Value
					}
				}
			}

			if countryValue != tt.wantCountry {
				t.Errorf("SearchCmd().country default = %q, want %q", countryValue, tt.wantCountry)
			}

			if radiusValue != tt.wantRadius {
				t.Errorf("SearchCmd().radius default = %d, want %d", radiusValue, tt.wantRadius)
			}
		})
	}
}

func Test_runSearchCmd(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
		errIs   error
	}{
		{
			name:    "missing API key returns ErrApiKeyMissing",
			apiKey:  "",
			wantErr: true,
			errIs:   ErrApiKeyMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JAMBASE_API_KEY", tt.apiKey)

			err := runSearchCmd(context.Background(), SearchCmd())
			if (err != nil) != tt.wantErr {
				t.Fatalf("runSearchCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errIs != nil && !errors.Is(err, tt.errIs) {
				t.Errorf("runSearchCmd() error = %v, want errors.Is(..., %v)", err, tt.errIs)
			}
		})
	}
}

func Test_apiKeyFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		want    string
		wantErr bool
		errIs   error
	}{
		{
			name:    "returns api key when present",
			apiKey:  "test-api-key",
			want:    "test-api-key",
			wantErr: false,
		},
		{
			name:    "returns error when missing",
			apiKey:  "",
			wantErr: true,
			errIs:   ErrApiKeyMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JAMBASE_API_KEY", tt.apiKey)

			got, err := apiKeyFromEnv()
			if (err != nil) != tt.wantErr {
				t.Fatalf("apiKeyFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errIs != nil && !errors.Is(err, tt.errIs) {
					t.Fatalf("apiKeyFromEnv() error = %v, want errors.Is(..., %v)", err, tt.errIs)
				}
				return
			}
			if got != tt.want {
				t.Errorf("apiKeyFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchAction(t *testing.T) {
	tests := []struct {
		name         string
		errored      bool
		events       []jambase.Event
		wantW        string
		wantContains []string
		wantErr      bool
	}{
		{
			name:    "returns searcher error",
			errored: true,
			wantErr: true,
		},
		{
			name:    "prints no shows found for empty result",
			events:  []jambase.Event{},
			wantW:   "no shows found\n",
			wantErr: false,
		},
		{
			name: "renders table for results",
			events: []jambase.Event{
				{
					Name:     "Goose at Red Rocks",
					Date:     "2024-07-04T20:00:00",
					DoorTime: "2024-07-04T19:00:00",
					Venue:    "Red Rocks Amphitheatre",
					Timezone: "America/Denver",
				},
			},
			wantContains: []string{"Goose at Red Rocks", "Thu Jul 4, 2024", "8:00 PM", "7:00 PM", "Red Rocks Amphitheatre", "America/Denver"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			searcher := mockShowSearcher{events: tt.events}
			if tt.errored {
				searcher.err = errors.New("boom")
			}

			err := searchAction(context.Background(), w, searcher, SearchOptions{})
			if (err != nil) != tt.wantErr {
				t.Fatalf("searchAction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			gotW := w.String()
			if tt.wantW != "" && gotW != tt.wantW {
				t.Errorf("searchAction() = %q, want %q", gotW, tt.wantW)
			}
			for _, wantFragment := range tt.wantContains {
				if !strings.Contains(gotW, wantFragment) {
					t.Errorf("searchAction() output = %q, want fragment %q", gotW, wantFragment)
				}
			}
		})
	}
}

func Test_formatEventDateAndStartTime(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "empty value",
			args:  args{value: ""},
			want:  "-",
			want1: "-",
		},
		{
			name:  "date only value",
			args:  args{value: "2024-07-04"},
			want:  "Thu Jul 4, 2024",
			want1: "-",
		},
		{
			name:  "datetime value",
			args:  args{value: "2024-07-04T20:00:00"},
			want:  "Thu Jul 4, 2024",
			want1: "8:00 PM",
		},
		{
			name:  "invalid value falls back",
			args:  args{value: "not-a-date"},
			want:  "not-a-date",
			want1: "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := formatEventDateAndStartTime(tt.args.value)
			if got != tt.want {
				t.Errorf("formatEventDateAndStartTime() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("formatEventDateAndStartTime() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_formatDoorTime(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty value",
			args: args{value: ""},
			want: "-",
		},
		{
			name: "hh:mm:ss value",
			args: args{value: "19:00:00"},
			want: "7:00 PM",
		},
		{
			name: "hh:mm value",
			args: args{value: "19:00"},
			want: "7:00 PM",
		},
		{
			name: "datetime value",
			args: args{value: "2024-07-04T19:00:00"},
			want: "7:00 PM",
		},
		{
			name: "invalid value falls back",
			args: args{value: "doors-open-ish"},
			want: "doors-open-ish",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDoorTime(tt.args.value); got != tt.want {
				t.Errorf("formatDoorTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
