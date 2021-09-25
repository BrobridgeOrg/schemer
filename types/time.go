package types

import (
	"fmt"
	"time"
)

type TimePercision int32

const (
	MicroInSecond = 1000000
	MilliInSecond = 1000
)

const (
	TIME_PERCISION_SECOND      TimePercision = 0
	TIME_PERCISION_MILLISECOND TimePercision = 1
	TIME_PERCISION_MICROSECOND TimePercision = 2
)

var TimePercisions = map[string]TimePercision{
	"second":      TIME_PERCISION_SECOND,
	"millisecond": TIME_PERCISION_MILLISECOND,
}

type Time struct {
	Percision TimePercision
}

func NewTime() *Time {
	return &Time{}
}

func (t *Time) Parse(data interface{}) error {

	props := data.(map[string]interface{})
	if v, ok := props["percision"]; ok {

		p, ok := TimePercisions[v.(string)]
		if !ok {
			return fmt.Errorf("Unsupported percision type: %v", v)
		}

		t.Percision = p
	}

	return nil
}

func (t *Time) getValueByPercision(d int64) time.Time {

	switch t.Percision {
	case TIME_PERCISION_MILLISECOND:
		return time.Unix(d/MilliInSecond, d%MilliInSecond*1000000)
	case TIME_PERCISION_MICROSECOND:
		return time.Unix(d/MicroInSecond, d%MicroInSecond*1000)
	}

	if d >= 1000000000000000 {
		// with microsecond
		return time.Unix(int64(d)/MicroInSecond, int64(d)%MicroInSecond*1000)
	} else if d >= 1000000000000 {
		// with millisecond
		return time.Unix(int64(d)/MilliInSecond, int64(d)%MilliInSecond*1000000)
	}

	return time.Unix(d, 0)
}

func (t *Time) GetValue(data interface{}) time.Time {

	switch d := data.(type) {
	case time.Time:
		return d
	case int64:
		return t.getValueByPercision(d)
	case uint64:
		return t.getValueByPercision(int64(d))
	case string:
		t, _ := time.Parse(time.RFC3339Nano, d)
		return t
	case float64:
		return t.getValueByPercision(int64(d))
	}

	return time.Unix(0, 0)
}
