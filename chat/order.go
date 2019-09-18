package chat

import (
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/order"
)

type orderOp string

const (
	opAddMeal    orderOp = "add_meal"
	opRemoveMeal orderOp = "remove_meal"
	opWhen       orderOp = "when"
	opWhere      orderOp = "where"

	opWhereHere     = "here"
	opWhereTakeaway = "takeaway"
)

// processOrder handles user order updates
func (s *service) processOrder(u *api.User, opData string) order.Order {
	splited := strings.Split(opData, "=")
	if len(splited) < 2 {
		return order.Order{}
	}

	op := orderOp(splited[0])
	val := splited[1]

	var o *order.Order

	switch op {
	case opAddMeal:
		o = s.order.AddMeal(u, val)

	case opRemoveMeal:
		o = s.order.RemoveMeal(u, val)

	case opWhen:
		t, err := time.Parse("15:04", val)
		if err != nil {
			log.Errorf("parsing time %+v:", val, err)
			return order.Order{}
		}
		o = s.order.SetTime(u, t)

	case opWhere:
		o = s.order.SetTakeaway(u, val == opWhereTakeaway)

	default:
	}

	// no need to pass pointer further
	return *o
}
