package main

import (
	"testing"
)

func Test_getVersion(t *testing.T) {
	tests := []struct {
		testName       string
		prerelase      string
		version        string
		expectedResult string
	}{
		// Test cases
		{
			testName:       "getVersion_WithDevPrelease_ShouldReturnVersionAndVersionPostfix",
			prerelase:      "dev",
			version:        "0.0.0",
			expectedResult: "0.0.0-dev",
		},
		{
			testName:       "getVersion_WithRCPrelease_ShouldReturnVersionAndVersionPostfix",
			prerelase:      "rc",
			version:        "0.0.1",
			expectedResult: "0.0.1-rc",
		},
		{
			testName:       "getVersion_WithoutPrelease_ShouldReturnVersionOnly",
			prerelase:      "",
			version:        "0.0.2",
			expectedResult: "0.0.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			Version = tt.version
			VersionPostfix = tt.prerelase
			printVersion()
			if Version != tt.expectedResult {
				t.Errorf("getVersion() = %v, expected %v", Version, tt.expectedResult)
			}
		})
	}
}
