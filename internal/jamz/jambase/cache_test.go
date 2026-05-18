package jambase

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func useTempHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
}

func seedCache(t *testing.T, dir, jsonData string) {
	t.Helper()
	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("HOME is not set")
	}

	questDir := filepath.Join(home, ".quest", dir)
	if err := os.MkdirAll(questDir, 0o700); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	cacheFilepath := filepath.Join(questDir, cacheFile)
	if err := os.WriteFile(cacheFilepath, []byte(jsonData), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
}

func Test_saveCache(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "saves cache successfully",
			args: args{
				key:   "testKey",
				value: "testValue",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if err := saveCache(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("saveCache() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			cacheData, err := loadCache()
			if err != nil {
				t.Fatalf("loadCache() error = %v", err)
			}
			if cacheData[tt.args.key] != tt.args.value {
				t.Errorf("saveCache() did not save the correct value, got %v, want %v", cacheData[tt.args.key], tt.args.value)
			}
		})
	}
}

func Test_loadCache(t *testing.T) {
	tests := []struct {
		name    string
		seedDir string
		seed    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "returns error when cache does not exist",
			seed:    "",
			wantErr: true,
		},
		{
			name:    "returns cache data when cache exists",
			seedDir: "jamz",
			seed:    `{"key1": "value1", "key2": "123"}`,
			want: map[string]interface{}{
				"key1": "value1",
				"key2": "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if tt.seed != "" {
				seedCache(t, tt.seedDir, tt.seed)
			}
			got, err := loadCache()
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadCache() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkCache(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		seedDir string
		seed    string
		want    string
		wantErr bool
		errIs   error
	}{
		{
			name: "successfully retrieves value for existing key",
			args: args{
				key: "existingKey",
			},
			seedDir: "jamz",
			seed:    `{"existingKey": "existingValue"}`,
			want:    "existingValue",
			wantErr: false,
		},
		{
			name: "returns error for non-existing key",
			args: args{
				key: "nonExistingKey",
			},
			seedDir: "jamz",
			seed:    `{"existingKey": "existingValue"}`,
			want:    "",
			wantErr: true,
			errIs:   ErrorCacheNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if tt.seed != "" {
				seedCache(t, tt.seedDir, tt.seed)
			}
			got, got1 := checkCache(tt.args.key)
			if got != tt.want {
				t.Errorf("checkCache() got = %v, want %v", got, tt.want)
			}
			if (got1 != nil) != tt.wantErr {
				t.Errorf("checkCache() got1 = %v, wantErr %v", got1, tt.wantErr)
			}
			if tt.errIs != nil && !errors.Is(got1, tt.errIs) {
				t.Errorf("checkCache() error = %v, wantErr %v", got1, tt.errIs)
			}
		})
	}
}
