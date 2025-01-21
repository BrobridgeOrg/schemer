package schemer_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/BrobridgeOrg/schemer"
	"github.com/stretchr/testify/assert"
)

var testTimeSource = `{
	"time_default": {
	    "type": "time"
	},
	"time_second": {
		"type": "time",
		"precision": "second"
	},
	"time_millisecond": {
		"type": "time",
		"precision": "millisecond"
	},
	"time_microsecond": {
		"type": "time",
		"precision": "microsecond"
	},
	"time_notSupport": {
		"type": "time",
		"precision": "notSupport"
	},
	"time_us": {
		"type": "time",
		"precision": "us"
	},
	"time_MICROSecond": {
		"type": "time",
		"precision": "MICROSecond"
	}
}`

type timeInput struct {
	time_default     string
	time_second      string
	time_millisecond string
	time_microsecond string
	time_notSupport  string
	time_us          string
	time_MICROSecond string
}

type timeExpected struct {
	time_default     time.Time
	time_second      time.Time
	time_millisecond time.Time
	time_microsecond time.Time
	time_notSupport  time.Time
	time_us          time.Time
	time_MICROSecond time.Time
}

var (
	SPECIAL_CHARACTERS = "!@#$%^&*()_+{}:<>?~`-=[]\\\\;',./"
	MAX_LENGTH_STRING  = strings.Repeat("a", 65536)
	NULL_TIME          time.Time
	UTC_PLUS_8         = time.FixedZone("UTC+08", 8*60*60)
)

func TimeTransform(t *testing.T, testSourceSchema *schemer.Schema, transformer *schemer.Transformer, input timeInput) map[string]interface{} {

	jsonInput := fmt.Sprintf(`
	{
		"time_default": "%s",
		"time_second": "%s",
		"time_millisecond": "%s",
		"time_microsecond": "%s",
		"time_notSupport": "%s",
		"time_us": "%s",
		"time_MICROSecond": "%s"
	}`, input.time_default, input.time_second, input.time_millisecond, input.time_microsecond, input.time_notSupport, input.time_us, input.time_MICROSecond)
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &rawData)
	if err != nil {
		t.Error(err)
		return nil
	}

	source := testSourceSchema.Normalize(rawData)

	returnedValue, err := transformer.Transform(nil, source)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	if !assert.Len(t, returnedValue, 1) {
		t.Fatal(err)
	}

	result := returnedValue[0]
	return result
}

func SetupTimeTransformer(t *testing.T, schema string) (*schemer.Transformer, *schemer.Schema) {

	testSourceSchema := schemer.NewSchema()
	err := schemer.UnmarshalJSON([]byte(schema), testSourceSchema)
	if err != nil {
		t.Error(err)
	}

	// Using the same schema for destination
	testDestSchema := schemer.NewSchema()
	err = schemer.UnmarshalJSON([]byte(schema), testDestSchema)
	if err != nil {
		t.Error(err)
	}

	// Create transformer
	transformer := schemer.NewTransformer(testSourceSchema, testDestSchema, schemer.WithRuntime(jsRuntime))
	err = transformer.SetScript(`return source`)
	if err != nil {
		t.Error(err)
	}

	return transformer, testSourceSchema
}

func TestTimePrecision(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)

	// second
	timetest1 := timeInput{time_second: "2024-08-06T15:02:00Z"}
	timetest1Expected := timeExpected{time_second: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)
	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	timetest2 := timeInput{time_second: "2024-08-06T15:02:00.123Z"}
	timetest2Expected := timeExpected{time_second: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest2)
	assert.Equal(t, timetest2Expected.time_second, result["time_second"].(time.Time).UTC())
	timetest3 := timeInput{time_second: "2024-08-06T15:02:00.123456Z"}
	timetest3Expected := timeExpected{time_second: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest3)
	assert.Equal(t, timetest3Expected.time_second, result["time_second"].(time.Time).UTC())
	timetest4 := timeInput{time_second: "2024-08-06T15:02:00.123456789Z"}
	timetest4Expected := timeExpected{time_second: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest4)
	assert.Equal(t, timetest4Expected.time_second, result["time_second"].(time.Time).UTC())
	timetest5 := timeInput{time_second: "2024-08-06T15:02:00.1234567890Z"}
	timetest5Expected := timeExpected{time_second: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest5)
	assert.Equal(t, timetest5Expected.time_second, result["time_second"].(time.Time).UTC())

	// millisecond
	timetest6 := timeInput{time_millisecond: "2024-08-06T15:02:00Z"}
	timetest6Expected := timeExpected{time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest6)
	assert.Equal(t, timetest6Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	timetest7 := timeInput{time_millisecond: "2024-08-06T15:02:00.123Z"}
	timetest7Expected := timeExpected{time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest7)
	assert.Equal(t, timetest7Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	timetest8 := timeInput{time_millisecond: "2024-08-06T15:02:00.123456Z"}
	timetest8Expected := timeExpected{time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest8)
	assert.Equal(t, timetest8Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	timetest9 := timeInput{time_millisecond: "2024-08-06T15:02:00.123456789Z"}
	timetest9Expected := timeExpected{time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest9)
	assert.Equal(t, timetest9Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	timetest10 := timeInput{time_millisecond: "2024-08-06T15:02:00.1234567890Z"}
	timetest10Expected := timeExpected{time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest10)
	assert.Equal(t, timetest10Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())

	// microsecond
	timetest11 := timeInput{time_microsecond: "2024-08-06T15:02:00Z"}
	timetest11Expected := timeExpected{time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest11)
	assert.Equal(t, timetest11Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	timetest12 := timeInput{time_microsecond: "2024-08-06T15:02:00.123Z"}
	timetest12Expected := timeExpected{time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest12)
	assert.Equal(t, timetest12Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	timetest13 := timeInput{time_microsecond: "2024-08-06T15:02:00.123456Z"}
	timetest13Expected := timeExpected{time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 123456000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest13)
	assert.Equal(t, timetest13Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	timetest14 := timeInput{time_microsecond: "2024-08-06T15:02:00.123456789Z"}
	timetest14Expected := timeExpected{time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 123456000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest14)
	assert.Equal(t, timetest14Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	timetest15 := timeInput{time_microsecond: "2024-08-06T15:02:00.1234567890Z"}
	timetest15Expected := timeExpected{time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 123456000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest15)
	assert.Equal(t, timetest15Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())

	// default precision
	timetest21 := timeInput{time_default: "2024-08-06T15:02:00Z"}
	timetest21Expected := timeExpected{time_default: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest21)
	assert.Equal(t, timetest21Expected.time_default, result["time_default"].(time.Time).UTC())
	timetest22 := timeInput{time_default: "2024-08-06T15:02:00.123Z"}
	timetest22Expected := timeExpected{time_default: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest22)
	assert.Equal(t, timetest22Expected.time_default, result["time_default"].(time.Time).UTC())
	timetest23 := timeInput{time_default: "2024-08-06T15:02:00.123456Z"}
	timetest23Expected := timeExpected{time_default: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest23)
	assert.Equal(t, timetest23Expected.time_default, result["time_default"].(time.Time).UTC())
	timetest24 := timeInput{time_default: "2024-08-06T15:02:00.123456789Z"}
	timetest24Expected := timeExpected{time_default: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest24)
	assert.Equal(t, timetest24Expected.time_default, result["time_default"].(time.Time).UTC())
	timetest25 := timeInput{time_default: "2024-08-06T15:02:00.1234567890Z"}
	timetest25Expected := timeExpected{time_default: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest25)
	assert.Equal(t, timetest25Expected.time_default, result["time_default"].(time.Time).UTC())
}

func TestTimeNotSupportPrecision(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)

	// not support
	timetest26 := timeInput{time_notSupport: "2024-08-06T15:02:00Z"}
	timetest26Expected := timeExpected{time_notSupport: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest26)
	assert.Equal(t, timetest26Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
	timetest27 := timeInput{time_notSupport: "2024-08-06T15:02:00.123Z"}
	timetest27Expected := timeExpected{time_notSupport: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest27)
	assert.Equal(t, timetest27Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
	timetest28 := timeInput{time_notSupport: "2024-08-06T15:02:00.123456Z"}
	timetest28Expected := timeExpected{time_notSupport: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest28)
	assert.Equal(t, timetest28Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
	timetest29 := timeInput{time_notSupport: "2024-08-06T15:02:00.123456789Z"}
	timetest29Expected := timeExpected{time_notSupport: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest29)
	assert.Equal(t, timetest29Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
	timetest30 := timeInput{time_notSupport: "2024-08-06T15:02:00.1234567890Z"}
	timetest30Expected := timeExpected{time_notSupport: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest30)
	assert.Equal(t, timetest30Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())

	// MICROSecond (mixed case)
	timetest31 := timeInput{time_MICROSecond: "2024-08-06T15:02:00Z"}
	timetest31Expected := timeExpected{time_MICROSecond: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest31)
	assert.Equal(t, timetest31Expected.time_MICROSecond, result["time_MICROSecond"].(time.Time).UTC())
	timetest32 := timeInput{time_MICROSecond: "2024-08-06T15:02:00.123Z"}
	timetest32Expected := timeExpected{time_MICROSecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest32)
	assert.Equal(t, timetest32Expected.time_MICROSecond, result["time_MICROSecond"].(time.Time).UTC())
	timetest33 := timeInput{time_MICROSecond: "2024-08-06T15:02:00.123456Z"}
	timetest33Expected := timeExpected{time_MICROSecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest33)
	assert.Equal(t, timetest33Expected.time_MICROSecond, result["time_MICROSecond"].(time.Time).UTC())
	timetest34 := timeInput{time_MICROSecond: "2024-08-06T15:02:00.123456789Z"}
	timetest34Expected := timeExpected{time_MICROSecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest34)
	assert.Equal(t, timetest34Expected.time_MICROSecond, result["time_MICROSecond"].(time.Time).UTC())
	timetest35 := timeInput{time_MICROSecond: "2024-08-06T15:02:00.1234567890Z"}
	timetest35Expected := timeExpected{time_MICROSecond: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest35)
	assert.Equal(t, timetest35Expected.time_MICROSecond, result["time_MICROSecond"].(time.Time).UTC())

	// us (abbreviation for microsecond precision)
	timetest36 := timeInput{time_us: "2024-08-06T15:02:00Z"}
	timetest36Expected := timeExpected{time_us: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest36)
	assert.Equal(t, timetest36Expected.time_us, result["time_us"].(time.Time).UTC())
	timetest37 := timeInput{time_us: "2024-08-06T15:02:00.123Z"}
	timetest37Expected := timeExpected{time_us: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest37)
	assert.Equal(t, timetest37Expected.time_us, result["time_us"].(time.Time).UTC())
	timetest38 := timeInput{time_us: "2024-08-06T15:02:00.123456Z"}
	timetest38Expected := timeExpected{time_us: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest38)
	assert.Equal(t, timetest38Expected.time_us, result["time_us"].(time.Time).UTC())
	timetest39 := timeInput{time_us: "2024-08-06T15:02:00.123456789Z"}
	timetest39Expected := timeExpected{time_us: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest39)
	assert.Equal(t, timetest39Expected.time_us, result["time_us"].(time.Time).UTC())
	timetest40 := timeInput{time_us: "2024-08-06T15:02:00.1234567890Z"}
	timetest40Expected := timeExpected{time_us: time.Date(2024, 8, 6, 15, 2, 0, 123000000, time.UTC)}
	result = TimeTransform(t, testTimeSourceSchema, transformer, timetest40)
	assert.Equal(t, timetest40Expected.time_us, result["time_us"].(time.Time).UTC())
}

func TestTimeTransformerWithMySQLFormat(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "2024-08-06 15:02:00", time_millisecond: "2024-08-06 15:02:00", time_microsecond: "2024-08-06 15:02:00", time_default: "2024-08-06 15:02:00", time_notSupport: "2024-08-06 15:02:00"}
	timetest1Expected := timeExpected{
		time_second:      time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC),
		time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC),
		time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC),
		time_default:     time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC),
		time_notSupport:  time.Date(2024, 8, 6, 15, 2, 0, 0, time.UTC),
	}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithTimeZone(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "2024-08-06T15:02:00+08:00", time_millisecond: "2024-08-06T15:02:00+08:00", time_microsecond: "2024-08-06T15:02:00+08:00", time_default: "2024-08-06T15:02:00+08:00", time_notSupport: "2024-08-06T15:02:00+08:00"}
	timetest1Expected := timeExpected{
		time_second:      time.Date(2024, 8, 6, 15, 2, 0, 0, UTC_PLUS_8).UTC(),
		time_millisecond: time.Date(2024, 8, 6, 15, 2, 0, 0, UTC_PLUS_8).UTC(),
		time_microsecond: time.Date(2024, 8, 6, 15, 2, 0, 0, UTC_PLUS_8).UTC(),
		time_default:     time.Date(2024, 8, 6, 15, 2, 0, 0, UTC_PLUS_8).UTC(),
		time_notSupport:  time.Date(2024, 8, 6, 15, 2, 0, 0, UTC_PLUS_8).UTC(),
	}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithNull(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "", time_millisecond: "", time_microsecond: "", time_default: "", time_notSupport: ""}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, nil, result["time_second"])
	assert.Equal(t, nil, result["time_millisecond"])
	assert.Equal(t, nil, result["time_microsecond"])
	assert.Equal(t, nil, result["time_default"])
	assert.Equal(t, nil, result["time_notSupport"])
}

func TestTimeTransformerWithSpace(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: " ", time_millisecond: " ", time_microsecond: " ", time_default: " ", time_notSupport: " "}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithRandomAlphabet(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "abc", time_millisecond: "abc", time_microsecond: "abc", time_default: "abc", time_notSupport: "abc"}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithChinese(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "中文", time_millisecond: "中文", time_microsecond: "中文", time_default: "中文", time_notSupport: "中文"}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithSpecialCharacters(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: SPECIAL_CHARACTERS, time_millisecond: SPECIAL_CHARACTERS, time_microsecond: SPECIAL_CHARACTERS, time_default: SPECIAL_CHARACTERS, time_notSupport: SPECIAL_CHARACTERS}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithMaxLenString(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: MAX_LENGTH_STRING, time_millisecond: MAX_LENGTH_STRING, time_microsecond: MAX_LENGTH_STRING, time_default: MAX_LENGTH_STRING, time_notSupport: MAX_LENGTH_STRING}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}

func TestTimeTransformerWithOnlyNumber(t *testing.T) {

	transformer, testTimeSourceSchema := SetupTimeTransformer(t, testTimeSource)
	timetest1 := timeInput{time_second: "5", time_millisecond: "5", time_microsecond: "5", time_default: "5", time_notSupport: "5"}
	timetest1Expected := timeExpected{time_second: NULL_TIME, time_millisecond: NULL_TIME, time_microsecond: NULL_TIME, time_default: NULL_TIME, time_notSupport: NULL_TIME}
	result := TimeTransform(t, testTimeSourceSchema, transformer, timetest1)

	assert.Equal(t, timetest1Expected.time_second, result["time_second"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_millisecond, result["time_millisecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_microsecond, result["time_microsecond"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_default, result["time_default"].(time.Time).UTC())
	assert.Equal(t, timetest1Expected.time_notSupport, result["time_notSupport"].(time.Time).UTC())
}
