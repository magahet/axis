package types

import (
	"time"

	"github.com/pkg/errors"
)

var timezone *time.Location

func init() {
	name, offset := time.Now().Zone()
	timezone = time.FixedZone(name, offset)
}

func ParseTime(timeString string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "could not parse time")
	}
	return t.In(timezone), nil
}

func ToLocalTime(timeString string) (string, error) {
	t, err := ParseTime(timeString)
	if err != nil {
		return timeString, errors.Wrap(err, "could not parse time")
	}
	return t.Format(time.RFC3339), nil
}
