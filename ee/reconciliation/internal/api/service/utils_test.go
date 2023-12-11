package service

import "testing"

type testVersion struct {
	version string
}

func (v testVersion) GetVersion() string {
	return v.version
}

func TestIsVersionSupported(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                string
		version             string
		minSupportedVersion string
		expectedResult      bool
	}

	testCases := []testCase{
		{
			name:                "commit hash",
			version:             "1234567890abcdef",
			minSupportedVersion: "v1.0.0",
			expectedResult:      true,
		},
		{
			name:                "higher version",
			version:             "1.0.1",
			minSupportedVersion: "v1.0.0",
			expectedResult:      true,
		},
		{
			name:                "equal version",
			version:             "1.0.0",
			minSupportedVersion: "v1.0.0",
			expectedResult:      true,
		},
		{
			name:                "lower version",
			version:             "0.9.0",
			minSupportedVersion: "v1.0.0",
			expectedResult:      false,
		},
		{
			name:                "versions with beta",
			version:             "2.0.0-beta.2",
			minSupportedVersion: "v2.0.0-beta.1",
			expectedResult:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v := testVersion{version: tc.version}
			result := isVersionSupported(v, tc.minSupportedVersion)
			if result != tc.expectedResult {
				t.Errorf("expected result %v, got %v", tc.expectedResult, result)
			}
		})
	}
}
