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
