package core

import (
	"reflect"
	"testing"
)

func Test_userHomeDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ensure this function returns a string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := userHomeDir()
			if reflect.TypeOf(value) != reflect.TypeOf("a string") {
				t.Errorf("userHomeDir() should always return a string")
			}
		})
	}
}
