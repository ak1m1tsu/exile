package client

import (
	"errors"
	"testing"
)

type testCase struct {
	name   string
	input  string
	expErr error
}

func Test_FetchAge(t *testing.T) {
	testCases := []testCase{
		{
			name:   "Good input",
			input:  "Ivan",
			expErr: nil,
		},
		{
			name:   "Bad input",
			input:  "1234567890",
			expErr: ErrFindAge,
		},
		{
			name:   "Empty input",
			input:  "",
			expErr: ErrNameEmpty,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := FetchAge(tc.input)
			if !errors.Is(err, tc.expErr) {
				t.Errorf("Expected error %v, got %v", tc.expErr, err)
			}
		})
	}
}

func Test_FetchGender(t *testing.T) {
	testCases := []testCase{
		{
			name:   "Good input",
			input:  "Ivan",
			expErr: nil,
		},
		{
			name:   "Bad input",
			input:  "1234567890",
			expErr: ErrFindGender,
		},
		{
			name:   "Empty input",
			input:  "",
			expErr: ErrNameEmpty,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := FetchGender(tc.input)
			if !errors.Is(err, tc.expErr) {
				t.Errorf("Expected error %v, got %v", tc.expErr, err)
			}
		})
	}
}

func Test_FetchNationality(t *testing.T) {
	testCases := []testCase{
		{
			name:   "Good input",
			input:  "Ivan",
			expErr: nil,
		},
		{
			name:   "Bad input",
			input:  "1234567890",
			expErr: ErrFindNation,
		},
		{
			name:   "Empty input",
			input:  "",
			expErr: ErrNameEmpty,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := FetchNationality(tc.input)
			if !errors.Is(err, tc.expErr) {
				t.Errorf("Expected error %v, got %v", tc.expErr, err)
			}
		})
	}
}
