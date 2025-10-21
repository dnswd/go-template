package util

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToPgtype converts uuid.UUID to pgtype.UUID
func UUIDToPgtype(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// PgtypeToUUID converts pgtype.UUID to uuid.UUID
// Returns error if pgtype.UUID is not valid
func PgtypeToUUID(pu pgtype.UUID) (uuid.UUID, error) {
	if !pu.Valid {
		return uuid.UUID{}, fmt.Errorf("pgtype.UUID is not valid")
	}
	return pu.Bytes, nil
}

// TimeToPgtype converts time.Time to pgtype.Timestamptz
func TimeToPgtype(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// PgtypeToTime converts pgtype.Timestamptz to time.Time
// Returns error if pgtype.Timestamptz is not valid or is infinite
func PgtypeToTime(pt pgtype.Timestamptz) (time.Time, error) {
	if !pt.Valid {
		return time.Time{}, fmt.Errorf("pgtype.Timestamptz is not valid")
	}
	if pt.InfinityModifier != pgtype.Finite {
		return time.Time{}, fmt.Errorf("pgtype.Timestamptz is infinite")
	}
	return pt.Time, nil
}

// DecimalToPgtype converts *apd.Decimal to pgtype.Numeric
// Returns NULL pgtype.Numeric if input is nil
func DecimalToPgtype(d *apd.Decimal) pgtype.Numeric {
	if d == nil {
		return pgtype.Numeric{
			Valid: false,
		}
	}

	// Handle special forms
	if d.Form == apd.Infinite {
		// PostgreSQL numeric doesn't support infinity
		return pgtype.Numeric{
			Valid: false,
		}
	}

	var num pgtype.Numeric
	num.Exp = d.Exponent
	num.NaN = d.Form == apd.NaN
	num.Valid = true

	// Convert coefficient to big.Int
	// apd.Decimal stores the absolute value in Coeff and sign separately
	num.Int = d.Coeff.MathBigInt()

	// Apply the negative sign to the big.Int
	if d.Negative {
		num.Int.Neg(num.Int)
	}

	return num
}

// PgtypeToDecimal converts pgtype.Numeric to *apd.Decimal
// Returns nil if pgtype.Numeric is not valid (NULL)
func PgtypeToDecimal(pn pgtype.Numeric) *apd.Decimal {
	if !pn.Valid {
		return nil
	}

	d := &apd.Decimal{
		Exponent: pn.Exp,
	}

	// Handle special values
	if pn.NaN {
		d.Form = apd.NaN
		return d
	}

	d.Form = apd.Finite

	// Set coefficient from big.Int
	if pn.Int != nil {
		// Check if negative
		d.Negative = pn.Int.Sign() < 0

		// apd.Decimal stores absolute value in Coeff
		if d.Negative {
			// Make a copy and get absolute value
			absInt := new(big.Int)
			absInt.Abs(pn.Int)
			d.Coeff.SetMathBigInt(absInt)
		} else {
			d.Coeff.SetMathBigInt(pn.Int)
		}
	}

	return d
}
