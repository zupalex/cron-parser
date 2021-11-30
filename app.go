package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type CronField int

const (
	MinutesField   CronField = 0
	HoursField     CronField = 1
	MonthDaysField CronField = 2
	MonthsField    CronField = 3
	WeekDaysField  CronField = 4
)

var (
	debugMode *bool = new(bool)
)

func debug(i ...interface{}) {
	if *debugMode {
		fmt.Print("[D] ")
		fmt.Println(i...)
	}
}

func debugf(format string, i ...interface{}) {
	if *debugMode {
		format = "[D] " + format + "\n"
		fmt.Printf(format, i...)
	}
}

type CronTab struct {
	Command   string
	Minutes   []int
	Hours     []int
	MonthDays []int
	Months    []int
	WeekDays  []int
}

func (ct *CronTab) GetDisplay() string {
	sb := strings.Builder{}

	sb.WriteString("minutes       ")
	sb.WriteString(joinIntSlice(ct.Minutes))
	sb.WriteString("\n")

	sb.WriteString("hour          ")
	sb.WriteString(joinIntSlice(ct.Hours))
	sb.WriteString("\n")

	sb.WriteString("day of month  ")
	sb.WriteString(joinIntSlice(ct.MonthDays))
	sb.WriteString("\n")

	sb.WriteString("month         ")
	sb.WriteString(joinIntSlice(ct.Months))
	sb.WriteString("\n")

	sb.WriteString("day of week   ")
	sb.WriteString(joinIntSlice(ct.WeekDays))
	sb.WriteString("\n")

	sb.WriteString("command       ")
	sb.WriteString(ct.Command)
	sb.WriteString("\n")

	return sb.String()
}

func joinIntSlice(input []int) string {
	sb := strings.Builder{}
	for _, i := range input {
		sb.WriteString(fmt.Sprintf("%v ", i))
	}
	return strings.TrimRight(sb.String(), " ")
}

func sanitizeInput(input string) string {
	tabs := regexp.MustCompile("(\\t+)")
	input = tabs.ReplaceAllString(input, " ")

	multiSpaces := regexp.MustCompile("\\s\\s+")
	input = multiSpaces.ReplaceAllString(input, " ")

	return strings.TrimSpace(input)
}

func getValidRange(fieldType CronField) ([]int, error) {
	switch fieldType {
	case MinutesField:
		return []int{0, 59}, nil
	case HoursField:
		return []int{0, 23}, nil
	case MonthDaysField:
		return []int{1, 31}, nil
	case MonthsField:
		return []int{1, 12}, nil
	case WeekDaysField:
		// Sunday can be 0 or 7
		return []int{0, 7}, nil
	default:
		return nil, errors.New("Invalid Cron field type")
	}
}

func parseTimeField(input string, fieldType CronField) ([]int, error) {
	expanded := map[int]bool{}

	intervals := strings.Split(input, ",")

	rangeRe := regexp.MustCompile("^(\\d\\d?)-(\\d\\d?)$")
	wildcardRe := regexp.MustCompile("^\\*(/\\d\\d?)?$")

	validRange, err := getValidRange(fieldType)
	if err != nil {
		return nil, err
	}

	debug("valid range:", validRange)

	for _, interval := range intervals {
		if m := rangeRe.FindAllStringSubmatch(interval, -1); len(m) == 1 {
			// Technically cannot throw an error here as we matched against digits but...
			start, err := strconv.Atoi(m[0][1])
			if err != nil {
				return nil, err
			}

			end, err := strconv.Atoi(m[0][2])
			if err != nil {
				return nil, err
			}

			if end < start {
				errMsg := fmt.Sprintf("Range error: start cannot be greater than end - %v > %v", start, end)
				return nil, errors.New(errMsg)
			}

			if start < validRange[0] {
				errMsg := fmt.Sprintf("Range error: invalid start of range - %v < %v (minimum value)", start, validRange[0])
				return nil, errors.New(errMsg)
			}

			if end > validRange[1] {
				errMsg := fmt.Sprintf("Range error: invalid end of range - %v > %v (maximum value)", end, validRange[1])
				return nil, errors.New(errMsg)
			}

			for i := start; i <= end; i++ {
				expanded[i] = true
			}
		} else if m := wildcardRe.FindAllStringSubmatch(interval, -1); len(m) == 1 {
			if len(m[0][1]) > 0 {
				frequency, err := strconv.Atoi(m[0][1][1:])
				if err != nil {
					return nil, err
				}

				if frequency > validRange[1] {
					errMsg := fmt.Sprintf("Frequency error: invalid frequency - %v > %v (maximum value)", frequency, validRange[1])
					return nil, errors.New(errMsg)
				}

				debugf("frequency: %v", frequency)

				maxNFreq := (int)(math.Floor((float64)(validRange[1] / frequency)))
				debug("maxNFreq:", maxNFreq)

				expanded[validRange[0]] = true

				for i := validRange[0]; i < maxNFreq+1; i++ {
					expanded[i*frequency+validRange[0]] = true
				}
			} else {
				for i := validRange[0]; i <= validRange[1]; i++ {
					expanded[i] = true
				}
			}
		} else {
			number, err := strconv.Atoi(interval)
			if err != nil {
				return nil, err
			}

			if number < validRange[0] {
				errMsg := fmt.Sprintf("Invalid value - %v < %v (minimum value)", number, validRange[0])
				return nil, errors.New(errMsg)
			}

			if number > validRange[1] {
				errMsg := fmt.Sprintf("Invalid value - %v > %v (maximum value)", number, validRange[1])
				return nil, errors.New(errMsg)
			}

			expanded[number] = true
		}
	}

	r := make([]int, len(expanded))

	idx := 0
	for k, _ := range expanded {
		r[idx] = k
		idx++
	}

	sort.Sort((sort.IntSlice)(r))

	// remove first 0 if we also have 7 (both are Sunday) for week days fields
	if fieldType == WeekDaysField && len(r) >= 2 && r[0] == 0 && r[len(r)-1] == 7 {
		r = r[1:]
	}

	return r, nil
}

func parseInput(input string) (*CronTab, error) {
	input = sanitizeInput(input)

	inputParts := strings.Split(input, " ")
	if len(inputParts) != 6 {
		errMsg := fmt.Sprintf("Malformed input. Expected 6 fields delimited by a space, got %v", len(inputParts))
		return nil, errors.New(errMsg)
	}

	ct := CronTab{Command: inputParts[5]}

	if len(ct.Command) == 0 {
		return nil, errors.New("Command cannot be empty")
	}

	var err error

	if _, err = os.Stat(ct.Command); err != nil {
		return nil, err
	}

	if ct.Minutes, err = parseTimeField(inputParts[0], MinutesField); err != nil {
		return nil, err
	}
	if ct.Hours, err = parseTimeField(inputParts[1], HoursField); err != nil {
		return nil, err
	}
	if ct.MonthDays, err = parseTimeField(inputParts[2], MonthDaysField); err != nil {
		return nil, err
	}
	if ct.Months, err = parseTimeField(inputParts[3], MonthsField); err != nil {
		return nil, err
	}
	if ct.WeekDays, err = parseTimeField(inputParts[4], WeekDaysField); err != nil {
		return nil, err
	}

	return &ct, nil
}

func main() {
	debugMode = flag.Bool("debug", false, "turn on debug mode")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Error: missing cron string")
		os.Exit(1)
	}

	input := flag.Args()[0]
	ct, err := parseInput(input)
	if err != nil {
		fmt.Println("Failed parsing cron string:", err)
		os.Exit(1)
	}

	fmt.Print(ct.GetDisplay())
}
