package main

import (
	"os/exec"
	"reflect"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

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

func Test_packBinary(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := packBinary(tt.args.name); got != tt.want {
				t.Errorf("packBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name string
		want []*exec.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPackNameFromConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name          string
		args          args
		wantPackNames []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPackNames := getPackNameFromConfig(tt.args.configPath); !reflect.DeepEqual(gotPackNames, tt.wantPackNames) {
				t.Errorf("getPackNameFromConfig() = %v, want %v", gotPackNames, tt.wantPackNames)
			}
		})
	}
}
