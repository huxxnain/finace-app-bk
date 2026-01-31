package utils

import "time"

// GetCurrentMonthYear returns the current month and year
func GetCurrentMonthYear() (year int, month int) {
	now := time.Now()
	return now.Year(), int(now.Month())
}

// GetMonthName returns the name of a month
func GetMonthName(month int) string {
	monthNames := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if month < 1 || month > 12 {
		return ""
	}
	return monthNames[month-1]
}
