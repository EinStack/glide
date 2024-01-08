package health

import (
	"fmt"
	"strconv"
	"strings"
)

const budgetSeparator = "/"

// ErrorBudget parses human-friendly error budget representation and return it as errors & update rate pair
// Error budgets could be set as a string in the following format: "10/s", "5/ms", "100/m" "1500/h"
type ErrorBudget struct {
	budget int
	unit   string
}

func DefaultErrorBudget() ErrorBudget {
	return ErrorBudget{
		budget: 10,
		unit:   "m",
	}
}

// Budget defines max allows number of errors per given time period
func (e *ErrorBudget) Budget() int {
	return e.budget
}

// RecoveryRate defines how much time do we need to wait to get one error token recovered (in microseconds)
func (e *ErrorBudget) RecoveryRate() int {
	return e.budget / e.unitToMicro(e.unit)
}

// MarshalText implements the encoding.TextMarshaler interface.
// This marshals the type and name as one string in the config.
func (b *ErrorBudget) MarshalText() (text []byte, err error) {
	return []byte(b.String()), nil
}

func (b *ErrorBudget) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), budgetSeparator)

	if len(parts) != 2 {
		return fmt.Errorf("invalid format")
	}

	budget, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("error parsing error number: %v", err)
	}

	unit := parts[1]

	if unit != "ms" && unit != "s" && unit != "m" && unit != "h" {
		return fmt.Errorf("invalid unit (supported: ms, s, m, h)")
	}

	b.budget = budget
	b.unit = unit

	return nil
}

func (e *ErrorBudget) unitToMicro(unit string) int {
	switch unit {
	case "ms":
		return 1000 // 1 ms = 1000 microseconds
	case "s":
		return 1000000 // 1 s = 1,000,000 microseconds
	case "m":
		return 60000000 // 1 m = 60,000,000 microseconds
	case "h":
		return 3600000000 // 1 h = 3,600,000,000 microseconds
	default:
		return 0 // or handle error
	}
}

// String returns the ID string representation as "type[/name]" format.
func (b *ErrorBudget) String() string {
	return strconv.Itoa(b.budget) + budgetSeparator + b.unit
}
