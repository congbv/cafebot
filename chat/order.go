package chat

import (
	"net/url"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	orderOpAddMeal     = "add_meal"
	orderOpRemoveMeal  = "remove_meal"
	orderOpSetTime     = "set_time"
	orderOpSetTakeaway = "set_takeaway"
)

// processOrder handles user order updates
func (s *service) processOrder(u *api.User, params url.Values) {
	if params == nil {
		return
	}

	for op, vals := range params {
		if len(vals) == 0 {
			continue
		}

		v := vals[0]

		switch op {
		case orderOpAddMeal:
			s.order.AddMeal(u, v)

		case orderOpRemoveMeal:
			s.order.RemoveMeal(u, v)

		case orderOpSetTime:
			t, err := time.Parse("15:04", v)
			if err != nil {
				log.Errorf("processOrder: set time: %s", err)
				continue
			}
			s.order.SetTime(u, t)

		case orderOpSetTakeaway:
			var takeaway bool
			if v == "1" {
				takeaway = true
			}
			s.order.SetTakeaway(u, takeaway)
		default:
		}
	}

}
