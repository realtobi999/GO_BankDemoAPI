package utils

import "time"

func CalculateAge(birthday time.Time) int {
	today := time.Now()
	age := today.Year() - birthday.Year()
	if today.Month() < birthday.Month() || (today.Month() == birthday.Month() && today.Day() < birthday.Day()) {
		age--
	}
	return age
}