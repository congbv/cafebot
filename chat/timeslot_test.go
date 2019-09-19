package chat

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateTimeSlots(t *testing.T) {
	t.Parallel()

	lowerLimit, _ := time.Parse("15:04", "09:00")
	upperLimit, _ := time.Parse("15:04", "21:00")

	interval := 15 * time.Minute
	cases := []struct {
		test          string
		nowFunc       func() time.Time
		expectedSlots []string
	}{{
		test: "simple case",
		nowFunc: func() time.Time {
			now, _ := time.Parse("15:04", "15:04")
			return now
		},
		expectedSlots: []string{
			"15:15", "15:30", "15:45", "16:00",
			"16:15", "16:30", "16:45", "17:00",
			"17:15", "17:30", "17:45", "18:00",
		},
	}, {
		test: "round time case",
		nowFunc: func() time.Time {
			now, _ := time.Parse("15:04", "15:00")
			return now
		},
		expectedSlots: []string{
			"15:15", "15:30", "15:45", "16:00",
			"16:15", "16:30", "16:45", "17:00",
			"17:15", "17:30", "17:45", "18:00",
		},
	}, {
		test: "late night case",
		nowFunc: func() time.Time {
			now, _ := time.Parse("15:04", "22:00")
			return now
		},
		expectedSlots: []string{},
	}, {
		test: "early morning case",
		nowFunc: func() time.Time {
			now, _ := time.Parse("15:04", "05:00")
			return now
		},
		expectedSlots: []string{
			"09:00", "09:15", "09:30", "09:45",
			"10:00", "10:15", "10:30", "10:45",
			"11:00", "11:15", "11:30", "11:45",
		},
	}, {
		test: "partly covered case",
		nowFunc: func() time.Time {
			now, _ := time.Parse("15:04", "19:30")
			return now
		},
		expectedSlots: []string{
			"19:45", "20:00", "20:15", "20:30",
			"20:45", "21:00",
		},
	}}

	for _, c := range cases {
		c := c
		t.Run(c.test, func(t *testing.T) {
			slots := generateTimeSlots(
				c.nowFunc(),
				interval,
				lowerLimit,
				upperLimit,
			)
			expectedGot := expectedGot(
				c.expectedSlots,
				slots,
			)
			if len(c.expectedSlots) != len(slots) {
				t.Fatal(expectedGot)
			}
			for i, s := range slots {
				if s != c.expectedSlots[i] {
					t.Fatal(expectedGot)
				}
			}
		})
	}

}

func expectedGot(expected, got []string) string {
	return fmt.Sprintf("expected: %+v, got: %+v", expected, got)
}
