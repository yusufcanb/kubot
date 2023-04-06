package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGetAbsolutePath(t *testing.T) {
	type args struct {
		filePath     string
		expectedPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "C:\\test.txt", args: args{filePath: "C:\\test.txt", expectedPath: "\\test.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetAbsolutePath(tt.args.filePath)
			if err != nil {
				assert.Equal(t, path, tt.args.expectedPath)
			}
		})
	}
}

func TestGetAbsolutePath1(t *testing.T) {
	type args struct {
		filePath string
	}

	// Windows test case
	t.Run("Windows", func(t *testing.T) {
		filePath := "C:\\Users\\Documents\\file.txt"
		expectedPath := "\\Users\\Documents\\file.txt"
		got, err := GetAbsolutePath(filePath)
		if err != nil {
			t.Errorf("GetAbsolutePath() error = %v", err)
			return
		}
		if got != expectedPath {
			t.Errorf("GetAbsolutePath() got = %v, want %v", got, expectedPath)
		}
	})

	// Unix test case
	t.Run("Unix", func(t *testing.T) {
		args := args{
			filePath: "/home/user/file.txt",
		}
		expectedPath := "/home/user/file.txt"

		got, err := GetAbsolutePath(args.filePath)
		if err != nil {
			t.Errorf("GetAbsolutePath() error = %v", err)
			return
		}

		if got != expectedPath {
			t.Errorf("GetAbsolutePath() got = %v, want %v", got, expectedPath)
		}
	})
}
