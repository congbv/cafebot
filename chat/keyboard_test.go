package chat

import (
	"testing"
	"time"
)

func TestGenerateTimeSlots(t *testing.T) {
	now, err := time.Parse("15:04", "15:04")
	if err != nil {
		t.Error(err)
	}

	interval := 15 * time.Minute
	count := 12
	expectedSlots := []string{
		"15:15", "15:30", "15:45", "16:00",
		"16:15", "16:30", "16:45", "17:00",
		"17:15", "17:30", "17:45", "18:00",
	}
	lowerLimit, _ := time.Parse("15:04", "09:00")
	upperLimit, _ := time.Parse("15:04", "21:00")

	slots := generateTimeSlots(now, interval, count, lowerLimit, upperLimit)
	for i, s := range slots {
		if s != expectedSlots[i] {
			t.Fatal("time slots were generated wrong")
		}
	}
}
