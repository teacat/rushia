package rushia

import "time"

// String returns a pointer string for SQL nullable.
func String(v string) *string {
	return &v
}

// Bool returns a pointer boolean for SQL nullable.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer int for SQL nullable.
func Int(v int) *int {
	return &v
}

// SliceInt returns a pointer slice of ints for SQL nullable.
func SliceInt(v []int) *[]int {
	return &v
}

// SliceString returns a pointer string for SQL nullable.
func SliceString(v []string) *[]string {
	return &v
}

// Float64 returns a pointer float64 for SQL nullable.
func Float64(v float64) *float64 {
	return &v
}

// Time returns a pointer time.Time for SQL nullable.
func Time(v time.Time) *time.Time {
	return &v
}
