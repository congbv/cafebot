package chat

import (
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

func generateTimeSlotsKeyboard(
	nextIntr intrEndpoint,
	slots []string,
	selectedTime *time.Time,
) [][]api.InlineKeyboardButton {
	var selectedSlot string
	if selectedTime != nil {
		selectedSlot = selectedTime.Format("15:04")
	}

	numInRow := 4
	totalSlots := len(slots)
	totalRows := totalSlots / numInRow

	if totalSlots == 0 || numInRow == 0 {
		return [][]api.InlineKeyboardButton{}
	}

	buttons := make([][]api.InlineKeyboardButton, totalRows)
	for i := range buttons {
		buttons[i] = make([]api.InlineKeyboardButton, numInRow)
	}

	for i := 0; i < totalRows; i++ {
		for j := 0; j < numInRow; j++ {
			slot := slots[i*numInRow+j]
			buttons[i][j] = newIntrButton(
				slot,
				newIntrData(nextIntr, "", opWhen, slot),
				selectedSlot == slot,
			)
		}
	}

	return buttons
}

func generateTimeSlots(
	startFrom time.Time,
	interval time.Duration,
	lowerLimit time.Time,
	upperLimit time.Time,
) []string {
	totalSlotNum := 12
	vals := make([]string, 0, totalSlotNum)

	val := startFrom.Round(interval)
	if val.After(upperLimit) || val.Equal(upperLimit) {
		return []string{}
	}

	if val.Before(startFrom) || val.Equal(startFrom) {
		val = val.Add(interval)
	}
	if val.Before(lowerLimit) {
		val = lowerLimit
	}

	for i := 0; i < totalSlotNum; i++ {
		if val.After(upperLimit) {
			break
		}
		vals = append(vals, val.Format("15:04"))
		val = val.Add(interval)
	}

	return vals
}
