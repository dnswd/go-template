package util

import (
	"math/big"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UUID Tests

func TestUUIDToPgtype(t *testing.T) {
	tests := []struct {
		name  string
		input uuid.UUID
	}{
		{
			name:  "valid uuid",
			input: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		},
		{
			name:  "nil uuid (all zeros)",
			input: uuid.Nil,
		},
		{
			name:  "random uuid",
			input: uuid.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UUIDToPgtype(tt.input)
			assert.True(t, result.Valid)
			// Compare as byte slices to avoid type assertion issues
			assert.Equal(t, tt.input[:], result.Bytes[:])
		})
	}
}

func TestPgtypeToUUID(t *testing.T) {
	tests := []struct {
		name      string
		input     pgtype.UUID
		expected  uuid.UUID
		expectErr bool
	}{
		{
			name: "valid pgtype uuid",
			input: pgtype.UUID{
				Bytes: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Valid: true,
			},
			expected:  uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			expectErr: false,
		},
		{
			name: "null pgtype uuid",
			input: pgtype.UUID{
				Valid: false,
			},
			expected:  uuid.UUID{},
			expectErr: true,
		},
		{
			name: "nil uuid bytes but valid",
			input: pgtype.UUID{
				Bytes: uuid.Nil,
				Valid: true,
			},
			expected:  uuid.Nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PgtypeToUUID(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUUIDRoundTrip(t *testing.T) {
	original := uuid.New()
	pg := UUIDToPgtype(original)
	result, err := PgtypeToUUID(pg)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

// Timestamptz Tests

func TestTimeToPgtype(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
	}{
		{
			name:  "current time",
			input: time.Now(),
		},
		{
			name:  "zero time",
			input: time.Time{},
		},
		{
			name:  "specific date with timezone",
			input: time.Date(2024, 1, 15, 14, 30, 45, 123456789, time.UTC),
		},
		{
			name:  "specific date with non-UTC timezone",
			input: time.Date(2024, 1, 15, 14, 30, 45, 0, time.FixedZone("EST", -5*3600)),
		},
		{
			name:  "unix epoch",
			input: time.Unix(0, 0),
		},
		{
			name:  "far future",
			input: time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:  "far past",
			input: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeToPgtype(tt.input)
			assert.True(t, result.Valid)
			assert.Equal(t, pgtype.Finite, result.InfinityModifier)
			assert.True(t, tt.input.Equal(result.Time))
		})
	}
}

func TestPgtypeToTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		input     pgtype.Timestamptz
		expected  time.Time
		expectErr bool
	}{
		{
			name: "valid timestamp",
			input: pgtype.Timestamptz{
				Time:             now,
				Valid:            true,
				InfinityModifier: pgtype.Finite,
			},
			expected:  now,
			expectErr: false,
		},
		{
			name: "null timestamp",
			input: pgtype.Timestamptz{
				Valid: false,
			},
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name: "positive infinity",
			input: pgtype.Timestamptz{
				Valid:            true,
				InfinityModifier: pgtype.Infinity,
			},
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name: "negative infinity",
			input: pgtype.Timestamptz{
				Valid:            true,
				InfinityModifier: pgtype.NegativeInfinity,
			},
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name: "zero time with finite modifier",
			input: pgtype.Timestamptz{
				Time:             time.Time{},
				Valid:            true,
				InfinityModifier: pgtype.Finite,
			},
			expected:  time.Time{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PgtypeToTime(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.expected.Equal(result))
			}
		})
	}
}

func TestTimestamptzRoundTrip(t *testing.T) {
	original := time.Now()
	pg := TimeToPgtype(original)
	result, err := PgtypeToTime(pg)
	require.NoError(t, err)
	assert.True(t, original.Equal(result))
}

// Numeric/Decimal Tests

func TestDecimalToPgtype(t *testing.T) {
	tests := []struct {
		name    string
		input   *apd.Decimal
		checkFn func(t *testing.T, result pgtype.Numeric)
	}{
		{
			name:  "nil decimal (NULL)",
			input: nil,
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.False(t, result.Valid)
			},
		},
		{
			name:  "zero",
			input: apd.New(0, 0),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int32(0), result.Exp)
				assert.Equal(t, int64(0), result.Int.Int64())
			},
		},
		{
			name: "positive integer",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("12345")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int64(12345), result.Int.Int64())
			},
		},
		{
			name: "negative integer",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("-12345")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int64(-12345), result.Int.Int64())
			},
		},
		{
			name: "decimal with fractional part",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("123.456")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int32(-3), result.Exp)
				assert.Equal(t, int64(123456), result.Int.Int64())
			},
		},
		{
			name: "negative decimal with fractional part",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("-123.456")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int32(-3), result.Exp)
				assert.Equal(t, int64(-123456), result.Int.Int64())
			},
		},
		{
			name: "very small decimal",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("0.000001")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				assert.Equal(t, int32(-6), result.Exp)
				assert.Equal(t, int64(1), result.Int.Int64())
			},
		},
		{
			name: "very large number",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("123456789012345678901234567890")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
				expected := new(big.Int)
				expected.SetString("123456789012345678901234567890", 10)
				assert.Equal(t, expected, result.Int)
			},
		},
		{
			name: "NaN",
			input: func() *apd.Decimal {
				d := &apd.Decimal{}
				d.Form = apd.NaN
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.True(t, result.NaN)
			},
		},
		{
			name: "positive infinity (should be invalid)",
			input: func() *apd.Decimal {
				d := &apd.Decimal{}
				d.Form = apd.Infinite
				d.Negative = false
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.False(t, result.Valid)
			},
		},
		{
			name: "negative infinity (should be invalid)",
			input: func() *apd.Decimal {
				d := &apd.Decimal{}
				d.Form = apd.Infinite
				d.Negative = true
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.False(t, result.Valid)
			},
		},
		{
			name: "scientific notation",
			input: func() *apd.Decimal {
				d, _, _ := apd.NewFromString("1.23E+10")
				return d
			}(),
			checkFn: func(t *testing.T, result pgtype.Numeric) {
				assert.True(t, result.Valid)
				assert.False(t, result.NaN)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecimalToPgtype(tt.input)
			tt.checkFn(t, result)
		})
	}
}

func TestPgtypeToDecimal(t *testing.T) {
	tests := []struct {
		name    string
		input   pgtype.Numeric
		checkFn func(t *testing.T, result *apd.Decimal)
	}{
		{
			name: "null numeric",
			input: pgtype.Numeric{
				Valid: false,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				assert.Nil(t, result)
			},
		},
		{
			name: "zero",
			input: pgtype.Numeric{
				Int:   big.NewInt(0),
				Exp:   0,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, apd.Finite, result.Form)
				assert.Equal(t, "0", result.String())
			},
		},
		{
			name: "positive integer",
			input: pgtype.Numeric{
				Int:   big.NewInt(12345),
				Exp:   0,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, "12345", result.String())
			},
		},
		{
			name: "negative integer",
			input: pgtype.Numeric{
				Int:   big.NewInt(-12345),
				Exp:   0,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.True(t, result.Negative)
				assert.Equal(t, "-12345", result.String())
			},
		},
		{
			name: "decimal with fractional part",
			input: pgtype.Numeric{
				Int:   big.NewInt(123456),
				Exp:   -3,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, int32(-3), result.Exponent)
				assert.Equal(t, "123.456", result.String())
			},
		},
		{
			name: "negative decimal with fractional part",
			input: pgtype.Numeric{
				Int:   big.NewInt(-123456),
				Exp:   -3,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.True(t, result.Negative)
				assert.Equal(t, int32(-3), result.Exponent)
				assert.Equal(t, "-123.456", result.String())
			},
		},
		{
			name: "very small decimal",
			input: pgtype.Numeric{
				Int:   big.NewInt(1),
				Exp:   -6,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, "0.000001", result.String())
			},
		},
		{
			name: "very large number",
			input: pgtype.Numeric{
				Int: func() *big.Int {
					i := new(big.Int)
					i.SetString("123456789012345678901234567890", 10)
					return i
				}(),
				Exp:   0,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, "123456789012345678901234567890", result.String())
			},
		},
		{
			name: "NaN",
			input: pgtype.Numeric{
				NaN:   true,
				Valid: true,
			},
			checkFn: func(t *testing.T, result *apd.Decimal) {
				require.NotNil(t, result)
				assert.Equal(t, apd.NaN, result.Form)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PgtypeToDecimal(tt.input)
			tt.checkFn(t, result)
		})
	}
}

func TestDecimalRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"integer", "12345"},
		{"negative", "-12345"},
		{"decimal", "123.456"},
		{"small decimal", "0.000001"},
		{"large number", "123456789012345678901234567890"},
		{"zero", "0"},
		{"negative decimal", "-123.456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original, _, err := apd.NewFromString(tt.input)
			require.NoError(t, err)

			pg := DecimalToPgtype(original)
			result := PgtypeToDecimal(pg)

			require.NotNil(t, result)
			assert.Equal(t, original.String(), result.String())
		})
	}
}

func TestNilDecimalRoundTrip(t *testing.T) {
	pg := DecimalToPgtype(nil)
	assert.False(t, pg.Valid)

	result := PgtypeToDecimal(pg)
	assert.Nil(t, result)
}
