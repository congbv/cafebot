package chat

import (
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

func generateTimeSlots(
	startFrom time.Time,
	interval time.Duration,
	totalSlotNum int,
	lowerLimit time.Time,
	upperLimit time.Time,
) []string {
	if totalSlotNum == 0 {
		totalSlotNum = 12
	}

	vals := make([]string, totalSlotNum)

	val := startFrom.Round(interval)
	if val.Before(startFrom) {
		val = val.Add(interval)
	}

	if val.Before(lowerLimit) {
		val = lowerLimit
	}

	for i := 0; i < totalSlotNum; i++ {
		vals[i] = val.Format("15:04")
		if val.Equal(upperLimit) || val.After(upperLimit) {
			break
		}
		val = val.Add(interval)
	}

	return vals
}

func generateTimeSlotsKeyboard(
	slots []string,
	numInRow int,
	nextInteraction string,
) [][]api.InlineKeyboardButton {
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
			buttons[i][j] = api.NewInlineKeyboardButtonData(
				slot,
				nextInteraction+"?when="+slot,
			)
		}
	}

	return buttons
}

func backKeyboardButton(prevInteraction string) []api.InlineKeyboardButton {
	return api.NewInlineKeyboardRow(
		api.NewInlineKeyboardButtonData(
			buttonText["step_back"],
			prevInteraction,
		),
	)
}
