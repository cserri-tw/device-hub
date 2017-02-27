package engine

import (
	"bytes"
	"encoding/binary"
	"testing"

	"bitbucket.org/tsetsova/decode-prototype/hub/expando"
	"github.com/stretchr/testify/assert"
)

func TestRawDecodeValid(t *testing.T) {

	t.Parallel()

	script := expando.Script{
		Main:    "xxx",
		Runtime: expando.Javascript,
		Input:   expando.Raw,
		Contents: `function xxx (input) {
				return {
						'value' : ((input[0] << 8) | input[1]) / 100,
					}
			}`,
	}

	// create the input payload as a byte array
	buf := &bytes.Buffer{}

	// multiple * 100 to ensure non floating point values
	initialValue := 22.33
	binary.Write(buf, binary.BigEndian, int16(initialValue*100))

	input := expando.Input{Payload: buf.Bytes()}

	e := New()
	result, err := e.Execute(script, input)

	assert.Nil(t, err)

	resultAsMap := result.(map[string]interface{})
	value := resultAsMap["value"]

	assert.Equal(t, initialValue, value)
}

func TestCSVDecodeValid(t *testing.T) {

	t.Parallel()

	script := expando.Script{
		Main:    "xxx",
		Runtime: expando.Javascript,
		Input:   expando.CSV,
		Contents: `function xxx (header, lines) {
				return {
						'header' : header,
						'lines' : lines
					}
			}`,
	}

	csv := "column1, column2\none, two\nthree, four\n five,six"
	input := expando.Input{Payload: []byte(csv)}

	e := New()
	result, err := e.Execute(script, input)
	assert.Nil(t, err)
	resultAsMap := result.(map[string]interface{})

	assert.Len(t, resultAsMap["header"], 2)
	assert.Len(t, resultAsMap["lines"], 3)

}

func TestJSONDecodeValid(t *testing.T) {

	t.Parallel()

	script := expando.Script{
		Main:    "xxx",
		Runtime: expando.Javascript,
		Input:   expando.JSON,
		Contents: `function xxx (input) {
				return input
			}`,
	}

	json := "{ \"a\" : 1}"
	input := expando.Input{Payload: []byte(json)}

	e := New()
	result, err := e.Execute(script, input)

	assert.Nil(t, err)
	resultAsMap := result.(map[string]interface{})

	assert.NotNil(t, resultAsMap["a"])
	assert.Equal(t, resultAsMap["a"], float64(1))
}
