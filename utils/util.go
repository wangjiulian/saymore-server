package utils

import "time"

func MaskPhone(phone string) string {
	if len(phone) < 7 { // Ensure phone number is long enough
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func YearsBetween(start, end time.Time) int {
	// Calculate the difference in full years
	years := end.Year() - start.Year()

	// Check if a year needs to be subtracted
	anniversary := start.AddDate(years, 0, 0) // Get the anniversary date
	if anniversary.After(end) {
		years-- // Not a full year yet, subtract one year
	}

	// If less than a year, still return 1
	if years < 1 {
		return 1
	}
	return years
}
