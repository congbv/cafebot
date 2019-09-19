package chat

import (
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

func whereKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhen

		hereButton := newIntrButton(
			buttonText["here"],
			newIntrData(nextIntr, "", opWhere, opWhereHere),
			o.Takeaway != nil && *o.Takeaway == false,
		)
		takeawayButton := newIntrButton(
			buttonText["takeaway"],
			newIntrData(nextIntr, "", opWhere, opWhereTakeaway),
			o.Takeaway != nil && *o.Takeaway == true,
		)

		keyboard := api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				hereButton,
				takeawayButton,
			),
		)

		return &keyboard
	}
}

func whenKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		nextIntr := intrWhat
		prevIntr := intrWhere

		// we need everything except hour and minute to be 0
		now, _ := time.Parse("15:04", time.Now().Format("15:04"))

		interval := time.Duration(conf.TimeSlotIntervalMin) * time.Minute

		buttonRows := append(
			generateTimeSlotsKeyboard(
				nextIntr,
				generateTimeSlots(
					now,
					interval,
					time.Time(conf.FirstOrderTime),
					time.Time(conf.LastOrderTime),
				),
				o.Time,
			),
			backKeyboardButton(prevIntr),
		)
		keyboard := api.NewInlineKeyboardMarkup(buttonRows...)

		return &keyboard
	}
}

func whatKeyboardFactory(conf config.CafeConfig) keyboardFunc {
	return func(intrData string, o order.Order) *api.InlineKeyboardMarkup {
		prevIntr := intrWhen

		if intrData == "" {
			buttonRows := append(
				generateMenuCategoryButtonRows(
					conf.Menu.Categories,
					intrWhat,
				),
				backKeyboardButton(prevIntr),
			)
			keyboard := api.NewInlineKeyboardMarkup(buttonRows...)
			return &keyboard
		}

		meals, ok := conf.Menu.Map[intrData]
		if !ok {
			log.Errorf("invalid menu category requested: %s", intrData)
			return nil
		}

		buttonRows := append(
			generateMenuMealButtonRows(intrData, meals, intrWhat, o),
			backKeyboardButton(intrWhat),
		)
		keyboard := api.NewInlineKeyboardMarkup(buttonRows...)
		return &keyboard
	}
}

func generateMenuCategoryButtonRows(
	categories []string,
	nextIntr intrEndpoint,
) [][]api.InlineKeyboardButton {
	buttonRows := make([][]api.InlineKeyboardButton, 0, len(categories))
	for _, category := range categories {
		intrData := newIntrData(nextIntr, category, "", "")
		buttonRows = append(
			buttonRows,
			[]api.InlineKeyboardButton{
				api.NewInlineKeyboardButtonData(
					category,
					intrData,
				),
			},
		)
	}
	return buttonRows
}

func generateMenuMealButtonRows(
	category string,
	meals []config.Meal,
	nextIntr intrEndpoint,
	o order.Order,
) [][]api.InlineKeyboardButton {
	buttonRows := make([][]api.InlineKeyboardButton, 0, len(meals))
	for _, meal := range meals {
		mealOp := opAddMeal
		selected := contains(o.Meal, meal.Val)
		if selected {
			mealOp = opRemoveMeal
		}

		buttonRows = append(
			buttonRows,
			[]api.InlineKeyboardButton{
				newIntrButton(
					meal.Val,
					newIntrData(
						nextIntr,
						category,
						mealOp,
						meal.Hash,
					),
					selected,
				),
			},
		)
	}
	return buttonRows
}

func backKeyboardButton(prevIntr intrEndpoint) []api.InlineKeyboardButton {
	return api.NewInlineKeyboardRow(
		api.NewInlineKeyboardButtonData(
			buttonText["back"],
			string(prevIntr),
		),
	)
}

func newIntrButton(text, data string, selected bool) api.InlineKeyboardButton {
	if text == "" || data == "" {
		return api.InlineKeyboardButton{}
	}
	if selected {
		text = buttonText["selected"] + text
	}
	return api.NewInlineKeyboardButtonData(text, data)
}

func contains(mm []string, m string) bool {
	for _, v := range mm {
		if v == m {
			return true
		}
	}
	return false
}
