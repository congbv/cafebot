package chat

import (
	"net/url"
	"testing"
)

func TestCutOffParams(t *testing.T) {
	t.Parallel()

	cases := []struct {
		test string

		data string

		expectedParams url.Values
		expectedData   string
		expectedHasErr bool
	}{{
		test: "valid case",

		data: "data?param1=value1&param2=value2",

		expectedParams: url.Values{
			"param1": []string{"value1"},
			"param2": []string{"value2"},
		},
		expectedData:   "data",
		expectedHasErr: false,
	}, {
		test: "empty data case",

		data: "",

		expectedParams: nil,
		expectedData:   "",
		expectedHasErr: true,
	}, {
		test: "empty data with params case",

		data: "?hello=world",

		expectedParams: nil,
		expectedData:   "",
		expectedHasErr: true,
	}, {
		test: "no params case",

		data: "data",

		expectedParams: nil,
		expectedData:   "data",
		expectedHasErr: false,
	}, {
		test: "hard with params",

		data: "data?ascii=%3Ckey%3A+0x90%3E",

		expectedParams: url.Values{
			"ascii": []string{"<key: 0x90>"},
		},
		expectedData:   "data",
		expectedHasErr: false,
	}}

	for _, c := range cases {
		c := c
		t.Run(c.test, func(t *testing.T) {
			data, params, err := cutOffParams(c.data)
			if err != nil && !c.expectedHasErr {
				t.Error("invalid err value")
			}

			if data != c.expectedData {
				t.Error("invalid data value")
			}

			for k := range params {
				_, ok := c.expectedParams[k]
				if !ok {
					t.Error("invalid params value")
				}
			}
		})
	}
}
