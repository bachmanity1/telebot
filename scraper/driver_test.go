package scraper

import (
	"reflect"
	"testing"
)

func TestIsValidTimeslot(t *testing.T) {
	for _, test := range []struct {
		prev    string
		next    string
		isValid bool
	}{
		{
			prev:    "",
			next:    "",
			isValid: false,
		},
		{
			prev:    "",
			next:    "2021-04-27 16:10~16:20",
			isValid: true,
		},
		{
			prev:    "2021-04-27 16:10~16:20",
			next:    "",
			isValid: false,
		},
		{
			prev:    "2021-04-27 16:10~16:20",
			next:    "2021-04-27 13:20~13:30",
			isValid: false,
		},
		{
			prev:    "2021-04-27 16:10~16:20",
			next:    "2021-04-28 16:10~16:20",
			isValid: false,
		},
		{
			prev:    "2021-04-27 16:10~16:20",
			next:    "2021-03-26 16:10~16:20",
			isValid: true,
		},
	} {
		if v := isValidTimeslot(test.prev, test.next); v != test.isValid {
			t.Errorf("got: %v, expected: %v -> prev: %s, next: %s", v, test.isValid, test.prev, test.next)
		}
	}
}

func TestGetPhoneNumber(t *testing.T) {
	for _, test := range []struct {
		input  string
		output []string
	}{
		{
			input:  "",
			output: nil,
		},
		{
			input:  "01023456789",
			output: []string{"010", "2345", "6789"},
		},
		{
			input:  "010234567891",
			output: nil,
		},
		{
			input:  "43334",
			output: nil,
		},
		{
			input:  "0-1-0-2-3-4-5-6-7-8-9",
			output: []string{"010", "2345", "6789"},
		},
		{
			input:  "0102456789",
			output: nil,
		},
		{
			input:  "0102345a789",
			output: nil,
		},
	} {
		if output := getPhoneNumber(test.input); !reflect.DeepEqual(output, test.output) {
			t.Errorf("input: %s -> got: %v, expected: %v", test.input, output, test.output)
		}
	}
}
