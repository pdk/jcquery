package typer

import (
	"strconv"
	"strings"
	"time"
)

// Type is what types can be inferred
type Type string

// Possible values for Type
const (
	Null      Type = "Null"
	Integer   Type = "Integer"
	Real      Type = "Real"
	Text      Type = "Text"
	Boolean   Type = "Boolean"
	Date      Type = "Date"
	Timestamp Type = "Timestamp"
)

func (t Type) String() string {
	return string(t)
}

// TimestampFormats has the list of parsable timestamp formats
var TimestampFormats = []struct {
	Name   string
	Format string
}{
	{Name: "ANSIC", Format: time.ANSIC},
	{Name: "UnixDate", Format: time.UnixDate},
	{Name: "RubyDate", Format: time.RubyDate},
	{Name: "RFC822", Format: time.RFC822},
	{Name: "RFC822Z", Format: time.RFC822Z},
	{Name: "RFC850", Format: time.RFC850},
	{Name: "RFC1123", Format: time.RFC1123},
	{Name: "RFC1123Z", Format: time.RFC1123Z},
	{Name: "RFC3339", Format: time.RFC3339},
	{Name: "RFC3339Nano", Format: time.RFC3339Nano},
	{Name: "Kitchen", Format: time.Kitchen},
	{Name: "Stamp", Format: time.Stamp},
	{Name: "StampMilli", Format: time.StampMilli},
	{Name: "StampMicro", Format: time.StampMicro},
	{Name: "StampNano", Format: time.StampNano},
	// 2008-06-17 10:35:49.618312
	{Name: "DateTimeMicro", Format: "2006-01-02 15:04:05.000000"},
}

// GuessedValue will contain a Type, and a single ...Value value.
type GuessedValue struct {
	GuessedType         Type
	IntValue            int64
	RealValue           float64
	BooleanValue        bool
	TimestampValue      time.Time
	TimestampFormatName string
	TextValue           string
}

// NullValue is a null thing
var NullValue GuessedValue = GuessedValue{GuessedType: Null}

// GuessType infers a type by trying a few conversions on a string
func GuessType(val string) GuessedValue {

	if val == "" || val == "null" {
		return NullValue
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return GuessedValue{
			GuessedType: Integer,
			IntValue:    i,
		}
	}

	r, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return GuessedValue{
			GuessedType: Real,
			RealValue:   r,
		}
	}

	switch strings.ToLower(val) {
	case "true":
	case "yes":
	case "t":
		return GuessedValue{
			GuessedType:  Boolean,
			BooleanValue: true,
		}
	case "false":
	case "no":
	case "f":
		return GuessedValue{
			GuessedType:  Boolean,
			BooleanValue: false,
		}
	}

	for _, format := range TimestampFormats {
		t, err := time.Parse(format.Format, val)
		if err == nil {
			return GuessedValue{
				GuessedType:         Timestamp,
				TimestampValue:      t,
				TimestampFormatName: format.Name,
			}
		}
	}

	return GuessedValue{
		GuessedType: Text,
		TextValue:   val,
	}
}
