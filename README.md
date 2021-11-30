<!-- omit in toc -->
# Cron Parser

A simple go program to parse a cron instruction passed as a single string and expands each field to show the times at which it will run.

<!-- omit in toc -->
## Index

- [Build](#build)
- [Usage](#usage)
- [Testing](#testing)

## Build

A [Golang environment](https://go.dev/doc/install) is required to be able to build the executable. There are no external dependencies. It is recommended to use the latest one available (1.17 at the time of writing this documentation).

Assuming a working Go environment is available, the executable can be built by running the following command in the directory containing the source files:

```shell
$ go build -o cron-parser .
```

This will generate a go executable `cron-parser` in the current directory.

## Usage

After following the instructions from the [Build](#build) sections, we can run the executable:

```shell
$ ./cron-parser [OPTIONS] CRON_STRING
```

Where:
- `CRON_STRING` is a cron instruction passed as a single string argument
- `OPTIONS` are optional flags.

Example:
```shell
$ ./cron-parser "*/15 0 1,15 * 1-5 /usr/bin/find"
minutes       0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

<!-- omit in toc -->
### Options list
| flag      | Description |
| ----------- | ----------- |
| --debug      | turn on debug mode (verbose) mode |

## Testing
To run the built-in tests:

```shell
$ go test . -v
```
 
Sample output:
```shell
=== RUN   TestSanitizeInput
--- PASS: TestSanitizeInput (0.00s)
=== RUN   TestParseMinutes
--- PASS: TestParseMinutes (0.00s)
=== RUN   TestParseHours
--- PASS: TestParseHours (0.00s)
=== RUN   TestParseMonthDays
--- PASS: TestParseMonthDays (0.00s)
=== RUN   TestParseMonths
--- PASS: TestParseMonths (0.00s)
=== RUN   TestParseWeekDays
--- PASS: TestParseWeekDays (0.00s)
=== RUN   TestFullParsing
--- PASS: TestFullParsing (0.00s)
PASS
ok      cron-parser     0.004s
```