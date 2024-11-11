package types

import (
	"fmt"
	"strings"
	"time"
)

type TimePrecision int32

const (
	MicroInSecond = 1000000
	MilliInSecond = 1000
)

const (
	TIME_PRECISION_SECOND      TimePrecision = 0
	TIME_PRECISION_MILLISECOND TimePrecision = 1
	TIME_PRECISION_MICROSECOND TimePrecision = 2
)

var TimePrecisions = map[string]TimePrecision{
	"second":      TIME_PRECISION_SECOND,
	"millisecond": TIME_PRECISION_MILLISECOND,
	"microsecond": TIME_PRECISION_MICROSECOND,
}

type Time struct {
	Precision TimePrecision
	Format    string
}

func NewTime() *Time {
	return &Time{}
}

func (t *Time) Parse(data interface{}) error {

	props := data.(map[string]interface{})
	if v, ok := props["precision"]; ok {

		p, ok := TimePrecisions[v.(string)]
		if !ok {
			return fmt.Errorf("Unsupported precision type: %v", v)
		}

		t.Precision = p
	}

	return nil
}

func (t *Time) getValueByPrecision(d int64) time.Time {

	switch t.Precision {
	case TIME_PRECISION_MILLISECOND:
		return time.Unix(d/MilliInSecond, d%MilliInSecond*1000000)
	case TIME_PRECISION_MICROSECOND:
		return time.Unix(d/MicroInSecond, d%MicroInSecond*1000)
	}

	// Auto detect precision
	if d >= 1000000000000000 {
		// with microsecond
		return time.Unix(int64(d)/MicroInSecond, int64(d)%MicroInSecond*1000)
	} else if d >= 1000000000000 {
		// with millisecond
		return time.Unix(int64(d)/MilliInSecond, int64(d)%MilliInSecond*1000000)
	}

	return time.Unix(d, 0)
}

func (t *Time) GetValue(data interface{}) (time.Time, error) {

	switch d := data.(type) {
	case time.Time:
		return d, nil
	case int64:
		return t.getValueByPrecision(d), nil
	case uint64:
		return t.getValueByPrecision(int64(d)), nil
	case string:

		if len(d) == 0 {
			return time.Unix(0, 0), ErrEmptyValue
		}

		t, err := time.Parse(time.RFC3339Nano, d)
		if err != nil {

			str := strings.Replace(d, " ", "T", 1)

			if d[len(d)-1:] != "Z" {
				t, _ = time.Parse(time.RFC3339Nano, str+"Z")
			} else {
				t, _ = time.Parse(time.RFC3339Nano, str)
			}

		}

		return t, nil
	case float64:
		return t.getValueByPrecision(int64(d)), nil
	}

	return time.Unix(0, 0), nil
}
