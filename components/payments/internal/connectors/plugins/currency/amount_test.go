package currency

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAmountWithPrecisionFromString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		amountString string
		precision    int
		expected     *big.Int
		expecterErr  error
	}

	testCases := []testCase{
		{
			name:         "nominal float string with precision 2",
			amountString: "123.45",
			precision:    2,
			expected:     big.NewInt(12345),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 2 (2)",
			amountString: "0.00",
			precision:    2,
			expected:     big.NewInt(0),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 0",
			amountString: "123.",
			precision:    0,
			expected:     big.NewInt(123),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 0 and decimals = precision",
			amountString: "123.",
			precision:    0,
			expected:     big.NewInt(123),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 0 and decimals = precision",
			amountString: "123.4567",
			precision:    4,
			expected:     big.NewInt(1234567),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 2 and decimals < precision",
			amountString: "123.",
			precision:    2,
			expected:     big.NewInt(12300),
			expecterErr:  nil,
		},
		{
			name:         "nominal float string with precision 2 and decimals < precision (2)",
			amountString: "123.1",
			precision:    2,
			expected:     big.NewInt(12310),
			expecterErr:  nil,
		},
		{
			name:         "nominal integer string with precision 2",
			amountString: "123",
			precision:    2,
			expected:     big.NewInt(12300),
			expecterErr:  nil,
		},
		{
			name:         "nominal integer string with precision 0",
			amountString: "123",
			precision:    0,
			expected:     big.NewInt(123),
			expecterErr:  nil,
		},

		// Error cases
		{
			name:        "negative precision",
			precision:   -1,
			expected:    nil,
			expecterErr: ErrInvalidPrecision,
		},
		{
			name:         "invalid amount multiple dots",
			amountString: "123.45.67",
			precision:    2,
			expected:     nil,
			expecterErr:  ErrInvalidAmount,
		},
		{
			name:         "invalid float amount",
			amountString: "123.4a",
			precision:    2,
			expected:     nil,
			expecterErr:  ErrInvalidAmount,
		},
		{
			name:         "invalid float amount (2)",
			amountString: "12a3.4",
			precision:    2,
			expected:     nil,
			expecterErr:  ErrInvalidAmount,
		},
		{
			name:         "invalid integer amount",
			amountString: "123a",
			precision:    2,
			expected:     nil,
			expecterErr:  ErrInvalidAmount,
		},
		{
			name:         "float number with decimal part > precision",
			amountString: "123.456",
			precision:    2,
			expected:     nil,
			expecterErr:  ErrInvalidPrecision,
		},
		{
			name:         "float number with decimal part > precision (2)",
			amountString: "123.456",
			precision:    0,
			expected:     nil,
			expecterErr:  ErrInvalidPrecision,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			amount, err := GetAmountWithPrecisionFromString(tc.amountString, tc.precision)
			if tc.expecterErr != nil {
				require.True(t, errors.Is(err, tc.expecterErr), "expected error %v, got %v", tc.expecterErr, err)
				return
			}

			if amount.Cmp(tc.expected) != 0 {
				t.Errorf("expected %v, got %v", tc.expected, amount)
			}
		})
	}
}

func TestGetStringAmountFromBigIntWithPrecision(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		amount         *big.Int
		precision      int
		expectedAmount string
		expectedError  error
	}

	testCases := []testCase{
		{
			name:           "precision 0",
			amount:         big.NewInt(12345),
			precision:      0,
			expectedAmount: "12345",
			expectedError:  nil,
		},
		{
			name:           "precision 0 (2)",
			amount:         big.NewInt(0),
			precision:      0,
			expectedAmount: "0",
			expectedError:  nil,
		},
		{
			name:           "precision 0 (3)",
			amount:         big.NewInt(123),
			precision:      0,
			expectedAmount: "123",
			expectedError:  nil,
		},
		{
			name:           "precision 2",
			amount:         big.NewInt(12345),
			precision:      2,
			expectedAmount: "123.45",
			expectedError:  nil,
		},
		{
			name:           "precision 2 (2)",
			amount:         big.NewInt(123),
			precision:      2,
			expectedAmount: "1.23",
			expectedError:  nil,
		},
		{
			name:           "precision 2 (3)",
			amount:         big.NewInt(12),
			precision:      2,
			expectedAmount: "0.12",
			expectedError:  nil,
		},
		{
			name:           "precision 2 (4)",
			amount:         big.NewInt(0),
			precision:      2,
			expectedAmount: "0.00",
			expectedError:  nil,
		},
		{
			name:           "precision > length of amount",
			amount:         big.NewInt(123),
			precision:      6,
			expectedAmount: "0.000123",
			expectedError:  nil,
		},

		// Error cases
		{
			name:           "negative precision",
			amount:         big.NewInt(123),
			precision:      -1,
			expectedAmount: "",
			expectedError:  ErrInvalidPrecision,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			amount, err := GetStringAmountFromBigIntWithPrecision(tc.amount, tc.precision)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError), "expected error %v, got %v", tc.expectedError, err)
				return
			}

			if amount != tc.expectedAmount {
				t.Errorf("expected %v, got %v", tc.expectedAmount, amount)
			}
		})
	}
}
