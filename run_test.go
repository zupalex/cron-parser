package main

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	test := "   */15   0 1,15 *      1-5 /usr/bin/find   "
	sanitizedTest := sanitizeInput(test)
	if sanitizedTest != "*/15 0 1,15 * 1-5 /usr/bin/find" {
		t.Errorf("[sanitizeInput] Unexpected output: %v", sanitizedTest)
	}
}

func compareIntSlices(i1, i2 []int) bool {
	if len(i1) != len(i2) {
		return false
	} else {
		for idx, v := range i1 {
			if v != i2[idx] {
				return false
			}
		}
	}

	return true
}

func TestParseMinutes(t *testing.T) {
	var err error

	checkParseOutput := func(test string, correct []int) {
		res, err := parseTimeField(test, MinutesField)
		if err != nil {
			t.Errorf("[parseTimeField] [MinutesField] Failed to parse '%v' - %v", test, err)
		} else if !compareIntSlices(res, correct) {
			t.Errorf("[parseTimeField] [MinutesField] Unexpected result for '%v' - Got %v, expected %v", test, res, correct)
		}
	}

	// Test we get the correct full range
	test := "*"
	expects := make([]int, 60)
	for i := 0; i < 60; i++ {
		expects[i] = i
	}
	checkParseOutput(test, expects)

	// Test valid interval
	test = "*/10"
	expects = []int{0, 10, 20, 30, 40, 50}
	checkParseOutput(test, expects)

	// Test valid interval at the boundary
	test = "*/31"
	expects = []int{0, 31}
	checkParseOutput(test, expects)

	// Test malformed entry
	test = "*/"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing wildcard '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test malformed entry
	test = "*/invalid"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing wildcard '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test out of range
	test = "*/65"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing wildcard '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test valid single entry
	test = "9"
	expects = []int{9}
	checkParseOutput(test, expects)

	// Test out of range
	test = "70"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test valid range
	test = "4-12"
	expects = []int{4, 5, 6, 7, 8, 9, 10, 11, 12}
	checkParseOutput(test, expects)

	// Test range with start == end
	test = "15-15"
	expects = []int{15}
	checkParseOutput(test, expects)

	// Test invalid range
	test = "10-5"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test out of range boundary
	test = "30-80"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test malformed range (shouldn't happen)
	test = "20 -30"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test valid list
	test = "4,8,15,16,23,42"
	expects = []int{4, 8, 15, 16, 23, 42}
	checkParseOutput(test, expects)

	// Test that same values are filtered out
	test = "4,4,1,10"
	expects = []int{1, 4, 10}
	checkParseOutput(test, expects)

	// Test out of range
	test = "40,80"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test malformed entry (shouldn't even happen)
	test = "4, 8"
	_, err = parseTimeField(test, MinutesField)
	if err == nil {
		t.Errorf("[parseTimeField] [MinutesField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MinutesField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}
}

func TestParseHours(t *testing.T) {
	var err error

	checkParseOutput := func(test string, correct []int) {
		res, err := parseTimeField(test, HoursField)
		if err != nil {
			t.Errorf("[parseTimeField] [HoursField] Failed to parse '%v' - %v", test, err)
		} else if !compareIntSlices(res, correct) {
			t.Errorf("[parseTimeField] [HoursField] Unexpected result for '%v' - Got %v, expected %v", test, res, correct)
		}
	}

	// Test we get the correct full range
	test := "*"
	expects := make([]int, 24)
	for i := 0; i < 24; i++ {
		expects[i] = i
	}
	checkParseOutput(test, expects)

	// Test valid interval
	test = "*/5"
	expects = []int{0, 5, 10, 15, 20}
	checkParseOutput(test, expects)

	// Test an interval with out of range frequency
	test = "*/30"
	_, err = parseTimeField(test, HoursField)
	if err == nil {
		t.Errorf("[parseTimeField] [HoursField] Parsing wildcard '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [HoursField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test a range with out of range bound
	test = "12-24"
	_, err = parseTimeField(test, HoursField)
	if err == nil {
		t.Errorf("[parseTimeField] [HoursField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [HoursField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}

	// Test a valid range
	test = "5-10"
	expects = []int{5, 6, 7, 8, 9, 10}
	checkParseOutput(test, expects)

	// Check a valid combination of range and list
	test = "2-5,20-23"
	expects = []int{2, 3, 4, 5, 20, 21, 22, 23}
	checkParseOutput(test, expects)
}

func TestParseMonthDays(t *testing.T) {
	var err error

	checkParseOutput := func(test string, correct []int) {
		res, err := parseTimeField(test, MonthDaysField)
		if err != nil {
			t.Errorf("[parseTimeField] [MonthDaysField] Failed to parse '%v' - %v", test, err)
		} else if !compareIntSlices(res, correct) {
			t.Errorf("[parseTimeField] [MonthDaysField] Unexpected result for '%v' - Got %v, expected %v", test, res, correct)
		}
	}

	// Test we get the correct full range
	test := "*"
	expects := make([]int, 31)
	for i := 0; i < 31; i++ {
		expects[i] = i + 1
	}
	checkParseOutput(test, expects)

	// Test a valid interval
	test = "*/4"
	expects = []int{1, 5, 9, 13, 17, 21, 25, 29}
	checkParseOutput(test, expects)

	// Test out of range
	test = "32"
	_, err = parseTimeField(test, MonthDaysField)
	if err == nil {
		t.Errorf("[parseTimeField] [MonthDaysField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MonthDaysField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}
}

func TestParseMonths(t *testing.T) {
	var err error

	checkParseOutput := func(test string, correct []int) {
		res, err := parseTimeField(test, MonthsField)
		if err != nil {
			t.Errorf("[parseTimeField] [MonthsField] Failed to parse '%v' - %v", test, err)
		} else if !compareIntSlices(res, correct) {
			t.Errorf("[parseTimeField] [MonthsField] Unexpected result for '%v' - Got %v, expected %v", test, res, correct)
		}
	}

	// Test we get the correct full range
	test := "*"
	expects := make([]int, 12)
	for i := 0; i < 12; i++ {
		expects[i] = i + 1
	}
	checkParseOutput(test, expects)

	// Test out of range
	test = "13"
	_, err = parseTimeField(test, MonthsField)
	if err == nil {
		t.Errorf("[parseTimeField] [MonthsField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [MonthsField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}
}

func TestParseWeekDays(t *testing.T) {
	var err error

	checkParseOutput := func(test string, correct []int) {
		res, err := parseTimeField(test, WeekDaysField)
		if err != nil {
			t.Errorf("[parseTimeField] [WeekDaysField] Failed to parse '%v' - %v", test, err)
		} else if !compareIntSlices(res, correct) {
			t.Errorf("[parseTimeField] [WeekDaysField] Unexpected result for '%v' - Got %v, expected %v", test, res, correct)
		}
	}

	// Test we get the correct full range
	test := "*"
	expects := make([]int, 7)
	for i := 0; i < 7; i++ {
		expects[i] = i + 1
	}
	checkParseOutput(test, expects)

	// Test we don't get Sunday twice (0 and 7)
	test = "0,3,7"
	expects = []int{3, 7}
	checkParseOutput(test, expects)

	// Test out of range
	test = "8"
	_, err = parseTimeField(test, WeekDaysField)
	if err == nil {
		t.Errorf("[parseTimeField] [WeekDaysField] Parsing range '%v' did not produce the expected error", test)
	} else {
		debugf("[parseTimeField] [WeekDaysField] Parsing wildcard '%v' produced expected error: %v", test, err)
	}
}

func TestFullParsing(t *testing.T) {
	var err error

	checkParseOutput := func(test, correct string) {
		ct, err := parseInput(test)
		if err != nil {
			t.Errorf("[parseInput] Failed parsing cron string '%v': %v", test, err)
		} else if ct.GetDisplay() != correct {
			t.Errorf("[parseInput] Unexpected result after parsing cron string '%v'\n---Got---\n%v\n---Expects---\n%v", test, ct.GetDisplay(), correct)
		}
	}

	// Test valid simple input
	test := "*/15 0 1,15 * 1-5 /usr/bin/find"
	expects := `minutes       0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
`

	checkParseOutput(test, expects)

	// Test valid input with range, list and interval
	test = "5,25,55 */3 */7 1-6 * /usr/bin/find"
	expects = `minutes       5 25 55
hour          0 3 6 9 12 15 18 21
day of month  1 8 15 22 29
month         1 2 3 4 5 6
day of week   1 2 3 4 5 6 7
command       /usr/bin/find
`

	checkParseOutput(test, expects)

	// Test valid input with mixed range, list and intervals in same field
	test = "2-5,15-19 12 3,6-8,*/15 * * /usr/bin/find"
	expects = `minutes       2 3 4 5 15 16 17 18 19
hour          12
day of month  1 3 6 7 8 16 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5 6 7
command       /usr/bin/find
`

	checkParseOutput(test, expects)

	// Test input with single digits formatted as '0X', e.g. '05' for '5'
	test = "*/5 00-05 * * 0,6,7 /usr/bin/whoami"
	expects = `minutes       0 5 10 15 20 25 30 35 40 45 50 55
hour          0 1 2 3 4 5
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   6 7
command       /usr/bin/whoami
`

	checkParseOutput(test, expects)

	// Test malformed input (missing field)
	test = "5,25,55 */7 1-6 * /usr/bin/find"
	_, err = parseInput(test)
	if err == nil {
		t.Errorf("[parseInput] Parsing cron string '%v' did not produce the expected error", test)
	} else {
		debugf("[parseInput] Parsing cron string '%v' produced expected error: %v", test, err)
	}

	// Test failure if command reference missing executable
	test = "*/15 0 1,15 * 1-5 /usr/bin/no_such_file"
	_, err = parseInput(test)
	if err == nil {
		t.Errorf("[parseInput] Parsing cron string '%v' did not produce the expected error", test)
	} else {
		debugf("[parseInput] Parsing cron string '%v' produced expected error: %v", test, err)
	}
}
